// Kantcoin Project
// https://kantcoin.org
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.

pragma solidity ^0.4.16;
///@title Voting with groups.
contract Campaign {
    //This represents a single voter.
    struct Voter {
        bytes32 pubkey; //The public keys that the ring signatures are composed of
        uint prefix; //02 or 03
        uint group;
        uint cash;
        bool hasGroup1; //If true, that person has already entered a group
        bool hasGroup2; //gVoter
    }

    //A group with N voters
    struct Group {
        address cPerson; //Who can vote on behalf of the other voters
        bytes32 category; //For statistics and/or donations (label from GroupCategory)
        uint value; //The amount of cash each voter has in this group
        uint size; //Number of voters
        bool closed; //A group can be closed in order to create another one
    }

    //Ballot info
    struct Ballot {
        bytes32 id;
        bool closed; //If it is closed, no more interaction is possible
        bool stopped; //Only confirmations are allowed
        bool donations; //Is this ballot intended to compute donations?
    }

    //This struct represents a vote containing a ring signature
    struct Vote {
        bytes32 fNumber; //First number of the URS
        uint candidate; //The chosen candidate
    }

    //This is a type for a single candidate.
    struct Candidate {
        bytes32 ipfs; //IPFS page with the candidate info
        uint votesCount; //Number of accumulated votes/donations
        address donee; //The address to receive donations on the Mainnet
    }

    struct Confirmation {
        address voter;
        bool ok;
    }

    //Can voters cancel their votes after being sent to candidates?
    bool public canCancel;

    //The creator of the campaign
    address public chairperson;

    //Who validates not fully confirmed votations
    address public validator;

    //Groups that are composed of voters
    Group[] public groups;

    //Ballots (for donations, votation, questions...)
    Ballot[] public ballots;

    //Voters, who must also be registered in groups
    mapping (address => Voter) voters;

    //Voters, who must also be registered in groups
    mapping (address => uint) votersNTRUHash;

    //User's hashcodes - should be unique
    uint[] public hashcodes;

    //How many donation ballots there will be
    uint public donationRounds;

    //How many votation ballots there will be
    uint public votationRounds;

    //How many donation ballots are left
    uint public remainingDonationRounds;

    //How many votation ballots are left
    uint public remainingVotationRounds;

    //It represents the ballot that voters are voting in
    uint public currentBallot;

    //It represents the current standard vote message that all voters must submit
    bytes32 public currentVoteMessage;

    //Maximum group size
    uint constant public mgz = 15;

    //Group chairpersons' tor addresses
    mapping (address => mapping (uint => bytes32)) tors;

    //Group mappings
    mapping (uint => mapping (uint => address)) gVoters;

    //Ballot mapping
    mapping (uint => mapping (uint => Candidate)) bCandidates;
    uint[255] bCandidatesCounter;

    //Ballot + candidate mapping
    mapping (uint => mapping (uint => uint)) cancellations;

    //For statistics
    mapping (uint => mapping (uint => mapping (bytes32 => uint))) votesPerBallotCandidateGCategory;

    //Ballot + group mappings
    mapping (uint => mapping (uint => mapping (uint => Vote))) bgVotes;
    mapping (uint => mapping (uint => mapping (uint => Confirmation))) bgConfirmations;
    mapping (uint => mapping (uint => mapping (uint => bool))) bgpCommitted;
    mapping (uint => mapping (uint => mapping (uint => bool))) bgpCommittedStatistics;
    mapping (uint => mapping (uint => bool)) bgMayCommit;
    mapping (uint => mapping (uint => bool)) validations;

    //Functions

    //Create a new campaign which can have several ballots within
    function Campaign(uint vRounds, uint dRounds) public payable{
        //The maximum number of ballots is 6 (donations: 1)
        require (vRounds > 0 && vRounds <= 5 && dRounds <= 1 && (vRounds + dRounds) <= 5);

        chairperson = msg.sender;
        votationRounds = vRounds;
        remainingVotationRounds = vRounds;
        donationRounds = dRounds;
        remainingDonationRounds = dRounds;
    }

    //It defines whether voters can cancel their votes after being sent to candidates
    function defineCanCancel(bool b){
        require (msg.sender == chairperson);
        canCancel = b;
    }

    //The insertion should be done after the creation, since there will be many candidates lists
    //Different ballots may have different lists of candidates
    function addCandidateIntoBallot(uint ballot, uint position, bytes32 ipfs, address donee) public {
        require (msg.sender == chairperson);
        require (bCandidates[ballot][position].ipfs == bytes32(0));
        bCandidates[ballot][position].ipfs = ipfs;
        bCandidates[ballot][position].donee = donee;
    }

    //In order to know how many candidates there are in a ballot
    function iterateCandidatesCounter(uint ballot) public{
        bCandidatesCounter[ballot] += 1;
    }

    //Get the candidate's ipfs
    function getCandidate(uint ballot, uint candidate) public view returns (bytes32 ipfs, uint count, address donee){
        ipfs = bCandidates[ballot][candidate].ipfs;
        count = bCandidates[ballot][candidate].votesCount;
        donee = bCandidates[ballot][candidate].donee;
    }

    //Insert new ballot in ballots array
    function addBallot(bytes32 id, bool isdonationsballot) public {
        require (msg.sender == chairperson);
        if (isdonationsballot){
            //There may only be a determined number of donation rounds
            require (remainingDonationRounds > 0);
            remainingDonationRounds -= 1;
        } else {
            //The same for votation rounds
            require (remainingVotationRounds > 0);
            remainingVotationRounds -= 1;
        }

        ballots.push(Ballot({
            id: id,
            closed: false,
            stopped: false,
            donations: isdonationsballot
            }));
    }

    //Voters interaction ends
    function closeBallot(uint ballot) public {
        require (msg.sender == chairperson);
        require (ballot < (donationRounds + votationRounds));
        require (!ballots[ballot].closed);
        ballots[ballot].closed = true;
    }

    //This will prevent voters from voting in this ballot
    function stopBallot(uint ballot) public {
        require (msg.sender == chairperson);
        require (ballot < (donationRounds + votationRounds));
        require (!ballots[ballot].closed);
        ballots[ballot].stopped = true;
    }

    //Allow voters to vote again
    function unstopBallot(uint ballot) public {
        require (msg.sender == chairperson);
        require (ballot < (donationRounds + votationRounds));
        require (!ballots[ballot].closed);
        ballots[ballot].stopped = false;
    }

    //It defines the ballot in which voters will vote and confirm the vote
    function defineCurrentBallot(uint ballot) public {
        require (msg.sender == chairperson);
        require (ballot < (donationRounds + votationRounds));
        require (!ballots[ballot].closed);
        currentBallot = ballot;
    }

    //Define the current standard vote message that all voters must submit
    function defineCurrentVoteMessage(bytes32 message) public {
        require (msg.sender == chairperson);
        currentVoteMessage = message;
    }

    //It sets the group chairperson's tor addresses and pubkeys
    function defineTor(address person, uint pos, bytes32 value) public {
        require (msg.sender == person);
        tors[person][pos] = value;
    }

    //It returns the group chairperson's tor address
    function getTor(address person, uint pos) public view returns (bytes32){
        return tors[person][pos];
    }

    //It increases by one unit the number of cancellations of some candidate
    function incrementCancellations(uint ballot, uint candidate) public {
        require (msg.sender == chairperson);
        require (ballots[ballot].closed);

        cancellations[ballot][candidate] += 1;
    }

    //It returns the number of cancellations of some candidate
    function getCancellations(uint ballot, uint candidate) public view returns (uint){
        return cancellations[ballot][candidate];
    }

    //Adding a group with its chairperson
    function addGroup(address cPerson) public {
        require (msg.sender == chairperson);
        groups.push(Group({
            cPerson: cPerson,
            category: bytes32(0),
            value: 0,
            size: 0,
            closed: false
            }));
    }

    //Defining the amount of cash of each voter inside this group
    function defineGroupValue(uint grp, uint value) public {
        require (msg.sender == chairperson);
        require ((value % donationRounds) == 0);
        require (groups[grp].size == 0);
        require (groups[grp].value == 0);

        groups[grp].value = value;
    }

    //Defining the category used to do the statistics
    function defineGroupCategory(uint grp, bytes32 category) public {
        require (msg.sender == chairperson);
        require (groups[grp].size == 0);
        require (groups[grp].category == bytes32(0));

        groups[grp].category = category;
    }

    //Close this group in order to replace its voters
    function closeGroup(uint grp) public {
        require (msg.sender == chairperson);
        groups[grp].closed = true;

        //Its voters are now free to enter into another group
        for (uint i = 0; i < mgz; i++) {
            voters[gVoters[grp][i]].hasGroup1 = false;
            voters[gVoters[grp][i]].hasGroup2 = false;
        }
    }

    //Give the voter the right to vote on this ballot.
    function giveRightToVote(address toVoter, uint prefix, bytes32 pubkey, uint cash) public {
        require (msg.sender == chairperson);
        voters[toVoter].pubkey = pubkey;
        voters[toVoter].prefix = prefix;
        voters[toVoter].cash = cash;
    }

    //If this voter is a troll
    function removeRightToVote(address toVoter) public {
        require (msg.sender == chairperson);
        voters[toVoter].prefix = 0;
    }

    //Add the voter to a group to he/she can vote
    function addVoterToGroup(address voter, uint grp) public {
        require (msg.sender == chairperson);
        require (!groups[grp].closed);
        require (!voters[voter].hasGroup1);
        require (groups[grp].size < mgz);
        require (groups[grp].value == voters[voter].cash);
        //The chairperson should give right to vote to this voter first
        require (voters[voter].prefix > 0);

        //Making the voter part of a group
        voters[voter].group = grp;
        voters[voter].hasGroup1 = true;
        groups[grp].size += 1;
    }

    //Adding voter to gVoters array
    function addVoterToGVoters(address voter, uint grp, uint position) public {
        require (msg.sender == chairperson);
        require (position < mgz);
        require (gVoters[grp][position] == address(0));
        require (voters[voter].group == grp);
        require (voters[voter].hasGroup1);
        require (!voters[voter].hasGroup2);

        gVoters[grp][position] = voter;
        voters[voter].hasGroup2 = true;
    }

    //Add voter's NTRU hashcode
    function defineVoterNTRUHash(address voter, uint hashcode) public {
        require (msg.sender == chairperson);
        votersNTRUHash[voter] = hashcode;
    }

    //Get voter's NTRU hascode to confirm his or her public key
    function getVoterNTRUHash(address voter) public view returns (uint){
        return votersNTRUHash[voter];
    }

    //Add the user hashcode to a list (it can be generated from the user name)
    function addVoterHashcode(uint hashcode) public {
        require (msg.sender == chairperson);
        hashcodes.push(hashcode);
    }

    //Search a hashcode (before inserting a new voter)
    function findVoterHashcode(uint hashcode) public view returns (bool){
        for (uint i = 0; i < hashcodes.length; i++){
            if (hashcode == hashcodes[i]){
                return true;
            }
        }
        return false;
    }

    //Get voter's info
    function getVoter(address voter) public view returns (bytes32 pubkey, uint prefix, uint group, uint cash, bool hasGroup1, bool hasGroup2){
        pubkey = voters[voter].pubkey;
        prefix = voters[voter].prefix;
        group = voters[voter].group;
        cash = voters[voter].cash;
        hasGroup1 = voters[voter].hasGroup1;
        hasGroup2 = voters[voter].hasGroup2;
    }

    //It returns the addresses of the members of a group
    function getGroupVoters(uint group) public view returns (address[mgz]){
        address[mgz] memory addresses;
        for (uint i = 0; i < mgz; i++){
            addresses[i] = gVoters[group][i];
        }
        return addresses;
    }

    //It returns the pubkeys of the members of a group
    function getGroupPubkeys(uint group) public view returns (uint[mgz], bytes32[mgz]){
        bytes32[mgz] memory pubkeys;
        uint[mgz] memory prefixes;

        for (uint i = 0; i < mgz; i++){
            pubkeys[i] = voters[gVoters[group][i]].pubkey;
            prefixes[i] = voters[gVoters[group][i]].prefix;
        }
        return (prefixes, pubkeys);
    }

    //It validates a ballot, allowing it to be committed even if there are not enough confirmations
    function validate(uint ballot, uint grp) public {
        require (msg.sender == validator);
        validations[ballot][grp] = true;
    }

    //Who validates not fully confirmed votations
    function defineValidator(address person) public {
        require (msg.sender == chairperson);
        require (validator == address(0));

        validator = person;
    }
    //Someone that will be in charge of checking the signatures and voting
    function defineGroupChairperson(address person, uint grp) public {
        require (msg.sender == chairperson);
        require (!groups[grp].closed);
        require ((grp < groups.length));

        groups[grp].cPerson = person;
    }

    //The group chairperson sends the votes, which are then confirmed by the voters
    function vote(uint ballot, uint grp, uint position, bytes32 first_number, uint the_candidate) public {
        require (msg.sender == groups[grp].cPerson);
        require (!groups[grp].closed);
        require (!ballots[ballot].closed);
        require (!ballots[ballot].stopped);
        require (ballot < (donationRounds + votationRounds));
        require (position < mgz);
        require (bgVotes[ballot][grp][position].fNumber == bytes32(0));

        //Verify if this "first number" has already been entered in the array
        for (uint i = 0; i < mgz; i++) {
            if (bgVotes[ballot][grp][i].fNumber == first_number){
                return;
            }
        }

        bgVotes[ballot][grp][position].fNumber = first_number;
        bgVotes[ballot][grp][position].candidate = the_candidate;
    }

    //For the statistics
    function getVotesPerBallotCandidateCategory(uint ballot, uint candidate, bytes32 category) public view returns (uint){
        return votesPerBallotCandidateGCategory[ballot][candidate][category];
    }

    //The voters have to verify if their messages was correctly assigned and then confirm the vote list
    function confirm(uint ballot, uint position, bool ok) public {
        //Get the sender's group
        uint grp = voters[msg.sender].group;

        require (position < mgz);
        require (bgConfirmations[ballot][grp][position].voter == address(0));
        require (!ballots[ballot].closed);
        require (ballots[ballot].stopped);
        //The voter should be part of a group in order to confirm the votation
        require (voters[msg.sender].hasGroup1);
        require (voters[msg.sender].hasGroup2);
        //The voter should have the right to vote
        require (voters[msg.sender].prefix > 0);

        for (uint i = 0; i < mgz; i++) {
            if (bgConfirmations[ballot][grp][i].voter == msg.sender){
                return;
            }
        }

        bgConfirmations[ballot][grp][position].voter = msg.sender;
        bgConfirmations[ballot][grp][position].ok = ok;
    }

    //It returns all sent confirmations regarding a ballot and a group
    function getConfirmations(uint ballot, uint grp) public view returns (address[mgz], bool[mgz]){
        address[mgz] memory addresses;
        bool[mgz] memory oks;
        for (uint i = 0; i < mgz; i++){
            addresses[i] = bgConfirmations[ballot][grp][i].voter;
            oks[i] = bgConfirmations[ballot][grp][i].ok;
        }
        return (addresses, oks);
    }

    //It returns all sent votes regarding a ballot and a group
    function getVotes(uint ballot, uint grp) public view returns (bytes32[mgz], uint[mgz]){
        bytes32[mgz] memory numbers;
        uint[mgz] memory candidates;
        for (uint i = 0; i < mgz; i++){
            numbers[i] = bgVotes[ballot][grp][i].fNumber;
            candidates[i] = bgVotes[ballot][grp][i].candidate;
        }
        return (numbers, candidates);
    }

    //It checks the requirements before committing
    //Solidity/Geth apparently does not work with long functions
    function preCommit(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the votation
        require (groups[grp].cPerson == msg.sender);
        require (!bgMayCommit[ballot][grp]);

        uint count1 = 0;
        uint count2 = 0;
        for (uint j = 0; j < mgz; j++) {
            //If a voter send an error message (not ok), something went wrong
            if ((bgConfirmations[ballot][grp][j].voter != address(0)) && !bgConfirmations[ballot][grp][j].ok && !validations[ballot][grp]){
                return;
            }
            if (bgVotes[ballot][grp][j].fNumber != bytes32(0)){
                count1 += 1;
            }
            if (bgConfirmations[ballot][grp][j].voter != address(0)){
                count2 += 1;
            }
        }
        //The same number of people who voted should confirm the vote
        //require (count1 == count2 || validations[ballot][grp]);
        bgMayCommit[ballot][grp] = true;
    }

    //Check if the group chairperson may commit
    function mayCommit(uint ballot, uint grp) public view returns (bool){
        return bgMayCommit[ballot][grp];
    }

    //Check whether the votation/donation was committed (for that ballot, group and position), or not
    function committed(uint ballot, uint grp, uint position) public view returns (bool){
        return bgpCommitted[ballot][grp][position];
    }

    //Check whether the votation/donation statistics was committed (for that ballot, group and position), or not
    function committedStatistics(uint ballot, uint grp, uint position) public view returns (bool){
        return bgpCommittedStatistics[ballot][grp][position];
    }

    //Committing the results and casting the votes
    function commitVotationPerPosition(uint ballot, uint grp, uint position) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        require (groups[grp].cPerson == msg.sender);
        //The ballot must not be a "donations ballot"
        require (!ballots[ballot].donations);
        //The votes must not have been committed before
        require (!bgpCommitted[ballot][grp][position]);

        if (bgVotes[ballot][grp][position].fNumber != bytes32(0)){
            //Get the chosen candidate
            uint candidate = bgVotes[ballot][grp][position].candidate;
            //Add this vote
            bCandidates[ballot][candidate].votesCount += 1;
            bgpCommitted[ballot][grp][position] = true;
        }
    }

    //Committing the statistics regarding the votation
    function commitVotationStatisticsPerPosition(uint ballot, uint grp, uint position) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        require (groups[grp].cPerson == msg.sender);
        //The ballot must not be a "donations ballot"
        require (!ballots[ballot].donations);
        //The votes must not have been committed before
        require (!bgpCommittedStatistics[ballot][grp][position]);

        if (bgVotes[ballot][grp][position].fNumber != bytes32(0)){
            //Get the chosen candidate
            uint candidate = bgVotes[ballot][grp][position].candidate;
            //Statistics
            bytes32 category = groups[grp].category;
            votesPerBallotCandidateGCategory[ballot][candidate][category] += 1;
            bgpCommittedStatistics[ballot][grp][position] = true;
        }
    }

    //Committing the results and casting the donations
    function commitDonations(uint ballot, uint grp, uint position) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the donations
        require (groups[grp].cPerson == msg.sender);
        //The ballot must be a "donations ballot"
        require (ballots[ballot].donations);
        //There must be cash in the group
        require (groups[grp].value > 0);
        //The donations must not have been committed before
        require (!bgpCommitted[ballot][grp][position]);

        //Calculating the donation amount
        uint donationAmount = groups[grp].value / donationRounds;

        if (bgVotes[ballot][grp][position].fNumber != bytes32(0)){
            //Get the chosen candidate
            uint candidate = bgVotes[ballot][grp][position].candidate;
            //Add this donation to the candidate/candidate's budget
            bCandidates[ballot][candidate].votesCount += donationAmount;
            bgpCommitted[ballot][grp][position]= true;
        }
    }

    //Committing statistics regarding donations
    function commitDonationsStatistics(uint ballot, uint grp, uint position) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the donations
        require (groups[grp].cPerson == msg.sender);
        //The ballot must be a "donations ballot"
        require (ballots[ballot].donations);
        //There must be cash in the group
        require (groups[grp].value > 0);
        //The donations must not have been committed before
        require (!bgpCommittedStatistics[ballot][grp][position]);

        //Calculating the donation amount
        uint donationAmount = groups[grp].value / donationRounds;

        if (bgVotes[ballot][grp][position].fNumber != bytes32(0)){
            //Get the chosen candidate
            uint candidate = bgVotes[ballot][grp][position].candidate;
            //Statistics
            bytes32 category = groups[grp].category;
            votesPerBallotCandidateGCategory[ballot][candidate][category] += donationAmount;
            bgpCommittedStatistics[ballot][grp][position] = true;
        }
    }

    //groups.length
    function howManyGroups() public view returns (uint){
        return groups.length;
    }

    //ballots.length
    function howManyBallots() public view returns (uint){
        return ballots.length;
    }

    //Candidates length
    function howManyCandidatesInBallot(uint ballot) public view returns (uint){
        return bCandidatesCounter[ballot];
    }
}