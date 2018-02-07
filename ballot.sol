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
    }

    struct Confirmation {
        address voter;
        bool ok;
    }

    struct Whisper {
        bytes32 part1;
        bytes32 part2;
        bytes32 part3;
    }

    //The creator of the campaign
    address public chairperson;

    //Groups that are composed of voters
    Group[] public groups;

    //Ballots (for donations, votation, questions...)
    Ballot[] public ballots;

    //Voters, who must also be registered in groups
    mapping (address => Voter) voters;

    //How many donation rounds there will be
    uint public donationRounds;

    //How many votation rounds there will be
    uint public votationRounds;

    //How many donation rounds are left
    uint public remainingDonationRounds;

    //How many votation rounds are left
    uint public remainingVotationRounds;

    //It represents the ballot that voters are voting in
    uint public currentBallot;

    //It represents the current standard vote message that all voters must submit
    bytes32 public currentVoteMessage;

    //Maximum group size
    uint constant public mgz = 12;

    //Group chairpersons' whispers
    mapping (address => Whisper) whispers;

    //Group mappings
    mapping (uint => mapping(uint => address)) gVoters;

    //Ballot mapping
    mapping (uint => mapping(uint => Candidate)) bCandidates;
    uint[255] bCandidatesCounter;

    //For statistics
    mapping (uint => mapping(uint => mapping (bytes32 => uint))) votesPerBallotCandidateBCategory;

    //Ballot + group mappings
    mapping (uint => mapping(uint => mapping (uint => Vote))) gbVotes;
    mapping (uint => mapping(uint => mapping (uint => Confirmation))) gbConfirmations;
    mapping (uint => mapping(uint => bool)) gbCommitted;
    mapping (uint => mapping(uint => bool)) gbMayCommit;
    mapping (uint => mapping(uint => bool)) gbCommittedStatistics;

    //Functions

    //Create a new campaign which can have several ballots within
    function Campaign(uint vRounds, uint dRounds) public payable{
        //The maximum number of rounds is 10 (donations: 9)
        require (vRounds + dRounds >= 1 && vRounds + dRounds <= 10);

        chairperson = msg.sender;
        votationRounds = vRounds;
        remainingVotationRounds = vRounds;
        donationRounds = dRounds;
        remainingDonationRounds = dRounds;
    }

    //The insertion should be done after the creation, since there will be many candidates lists
    //Different ballots may have different lists of candidates
    function addCandidateIntoBallot(uint ballot, uint position, bytes32 ipfs) public {
        require (msg.sender == chairperson);
        require (bCandidates[ballot][position].ipfs == bytes32(0));
        bCandidates[ballot][position].ipfs = ipfs;
    }

    //In order to know how many candidates there are in a ballot
    function iterateCandidatesCounter(uint ballot) public{
        bCandidatesCounter[ballot] += 1;
    }

    //Get the candidate's ipfs
    function getCandidate(uint ballot, uint candidate) public view returns (bytes32 ipfs, uint count){
        ipfs = bCandidates[ballot][candidate].ipfs;
        count = bCandidates[ballot][candidate].votesCount;
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

    //Sets the group chairperson's whispers
    function defineWhisper(address person, bytes32 part1, bytes32 part2, bytes32 part3) public {
        require (msg.sender == person);
        whispers[person].part1 = part1;
        whispers[person].part2 = part2;
        whispers[person].part3 = part3;
    }

    //Returns the group chairperson's whisper address
    function getWhisper(address person) public view returns (bytes32, bytes32, bytes32){
        return (whispers[person].part1, whispers[person].part2, whispers[person].part3);
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

    //Adgd the voter to a group to he/she can vote
    //May only be called by the voter.
    function addVoterToGroup(uint grp) public {
        require (!groups[grp].closed);
        require (!voters[msg.sender].hasGroup1);
        require (groups[grp].size < mgz);
        require (groups[grp].value == voters[msg.sender].cash);
        //The chairperson should give right to vote to this voter first
        require (voters[msg.sender].prefix > 0);

        //Making the voter part of a group
        voters[msg.sender].group = grp;
        voters[msg.sender].hasGroup1 = true;
        groups[grp].size += 1;
    }

    //Adding voter to gVoters array
    function addVoterToGVoters(uint grp, uint position) public {
        require (position < mgz);
        require (gVoters[grp][position] == address(0));
        require (voters[msg.sender].group == grp);
        require (voters[msg.sender].hasGroup1);
        require (!voters[msg.sender].hasGroup2);

        gVoters[grp][position] = msg.sender;
        voters[msg.sender].hasGroup2 = true;
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

    //Returns the addresses of the members of a group
    function getGroupVoters(uint group) public view returns (address[mgz]){
        address[mgz] memory addresses;
        for (uint i = 0; i < mgz; i++){
            addresses[i] = gVoters[group][i];
        }
        return addresses;
    }

    //Returns the pubkeys of the members of a group
    function getGroupPubkeys(uint group) public view returns (uint[mgz], bytes32[mgz]){
        bytes32[mgz] memory pubkeys;
        uint[mgz] memory prefixes;

        for (uint i = 0; i < mgz; i++){
            pubkeys[i] = voters[gVoters[group][i]].pubkey;
            prefixes[i] = voters[gVoters[group][i]].prefix;
        }
        return (prefixes, pubkeys);
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
        require (gbVotes[ballot][grp][position].fNumber == bytes32(0));

        //Verify if this "first number" has already been entered in the array
        for (uint i = 0; i < mgz; i++) {
            if (gbVotes[ballot][grp][i].fNumber == first_number){
                return;
            }
        }

        gbVotes[ballot][grp][position].fNumber = first_number;
        gbVotes[ballot][grp][position].candidate = the_candidate;
    }

    //For the statistics
    function getVotesPerBallotCandidateCategory(uint ballot, uint candidate, bytes32 category) public view returns (uint){
        return votesPerBallotCandidateBCategory[ballot][candidate][category];
    }

    //The voters have to verify if their messages was correctly assigned and then confirm the vote list
    function confirm(uint ballot, uint position, bool ok) public {
        //Get the sender's group
        uint grp = voters[msg.sender].group;

        require (position < mgz);
        require (gbConfirmations[ballot][grp][position].voter == address(0));
        require (!ballots[ballot].closed);
        require (ballots[ballot].stopped);
        //The voter should be part of a group in order to confirm the votation
        require (voters[msg.sender].hasGroup1);
        require (voters[msg.sender].hasGroup2);
        //The voter should have the right to vote
        require (voters[msg.sender].prefix > 0);

        for (uint i = 0; i < mgz; i++) {
            if (gbConfirmations[ballot][grp][i].voter == msg.sender){
                return;
            }
        }

        gbConfirmations[ballot][grp][position].voter = msg.sender;
        gbConfirmations[ballot][grp][position].ok = ok;
    }

    //Returns all sent confirmations regarding a ballot and a group
    function getConfirmations(uint ballot, uint grp) public view returns (address[mgz], bool[mgz]){
        address[mgz] memory addresses;
        bool[mgz] memory oks;
        for (uint i = 0; i < mgz; i++){
            addresses[i] = gbConfirmations[ballot][grp][i].voter;
            oks[i] = gbConfirmations[ballot][grp][i].ok;
        }
        return (addresses, oks);
    }

    //Returns all sent votes regarding a ballot and a group
    function getVotes(uint ballot, uint grp) public view returns (bytes32[mgz], uint[mgz]){
        bytes32[mgz] memory numbers;
        uint[mgz] memory candidates;
        for (uint i = 0; i < mgz; i++){
            numbers[i] = gbVotes[ballot][grp][i].fNumber;
            candidates[i] = gbVotes[ballot][grp][i].candidate;
        }
        return (numbers, candidates);
    }

    //Checks the requirements before committing
    //Solidity/Geth apparently does not work with long functions
    function preCommit(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the votation
        require (groups[grp].cPerson == msg.sender);
        require (!gbMayCommit[ballot][grp]);

        uint count1 = 0;
        uint count2 = 0;
        for (uint j = 0; j < mgz; j++) {
            //If a voter send an error message (not ok), something went wrong
            if ((gbConfirmations[ballot][grp][j].voter != address(0)) && !gbConfirmations[ballot][grp][j].ok){
                return;
            }
            if (gbVotes[ballot][grp][j].fNumber != bytes32(0)){
                count1 += 1;
            }
            if (gbConfirmations[ballot][grp][j].voter != address(0)){
                count2 += 1;
            }
        }
        //The same number of people who voted should confirm the vote
        require (count1 == count2);
        gbMayCommit[ballot][grp] = true;
    }

    //Check if the group chairperson may commit
    function mayCommit(uint ballot, uint grp) public view returns (bool){
        return gbMayCommit[ballot][grp];
    }

    //Check whether the votation/donation was committed or not
    function committed(uint ballot, uint grp) public view returns (bool){
        return gbCommitted[ballot][grp];
    }

    //Check whether the votation/donation statistics was committed or not
    function committedStatistics(uint ballot, uint grp) public view returns (bool){
        return gbCommittedStatistics[ballot][grp];
    }

    //Committing the results and casting the votes
    function commitVotation(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        require (groups[grp].cPerson == msg.sender);
        //The ballot must not be a "donations ballot"
        require (!ballots[ballot].donations);
        //The votes must not have been committed before
        require (!gbCommitted[ballot][grp]);

        for (uint i = 0; i < mgz; i++) {
            if (gbVotes[ballot][grp][i].fNumber != bytes32(0)){
                //Get the chosen candidate
                uint candidate = gbVotes[ballot][grp][i].candidate;
                //Add this vote
                bCandidates[ballot][candidate].votesCount += 1;
            }
        }
        gbCommitted[ballot][grp] = true;
    }

    //Committing the statistics regarding the votation
    function commitVotationStatistics(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        require (groups[grp].cPerson == msg.sender);
        //The ballot must not be a "donations ballot"
        require (!ballots[ballot].donations);
        //The votes must not have been committed before
        require (!gbCommittedStatistics[ballot][grp]);

        for (uint i = 0; i < mgz; i++) {
            if (gbVotes[ballot][grp][i].fNumber != bytes32(0)){
                //Get the chosen candidate
                uint candidate = gbVotes[ballot][grp][i].candidate;
                //Statistics
                bytes32 category = groups[grp].category;
                votesPerBallotCandidateBCategory[ballot][candidate][category] += 1;
            }
        }
        gbCommittedStatistics[ballot][grp] = true;
    }

    //Committing the results and casting the donations
    function commitDonations(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the donations
        require (groups[grp].cPerson == msg.sender);
        //The ballot must be a "donations ballot"
        require (ballots[ballot].donations);
        //There must be cash in the group
        require (groups[grp].value > 0);
        //The donations must not have been committed before
        require (!gbCommitted[ballot][grp]);

        //Calculating the donation amount
        uint donationAmount = groups[grp].value / donationRounds;

        for (uint i = 0; i < mgz; i++) {
            if (gbVotes[ballot][grp][i].fNumber != bytes32(0)){
                //Get the chosen candidate
                uint candidate = gbVotes[ballot][grp][i].candidate;
                //Add this donation to the candidate/candidate's budget
                bCandidates[ballot][candidate].votesCount += donationAmount;
            }
        }
        gbCommitted[ballot][grp]= true;
    }

    //Committing statistics regarding donations
    function commitDonationsStatistics(uint ballot, uint grp) public {
        require (!groups[grp].closed);
        require (ballots[ballot].closed); //The ballot must be closed
        //Only the group chairperson can commit the donations
        require (groups[grp].cPerson == msg.sender);
        //The ballot must be a "donations ballot"
        require (ballots[ballot].donations);
        //There must be cash in the group
        require (groups[grp].value > 0);
        //The donations must not have been committed before
        require (!gbCommittedStatistics[ballot][grp]);

        //Calculating the donation amount
        uint donationAmount = groups[grp].value / donationRounds;

        for (uint i = 0; i < mgz; i++) {
            if (gbVotes[ballot][grp][i].fNumber != bytes32(0)){
                //Get the chosen candidate
                uint candidate = gbVotes[ballot][grp][i].candidate;
                //Statistics
                bytes32 category = groups[grp].category;
                votesPerBallotCandidateBCategory[ballot][candidate][category] += donationAmount;
            }
        }
        gbCommittedStatistics[ballot][grp]= true;
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