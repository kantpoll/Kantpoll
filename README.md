# About us 

[Kantcoin](https://kantcoin.org) is a free and open-source election management system guided by the following principles:

## Our 10 principles

1. **Transparency** - people should not be fooled or manipulated.
2. **Privacy** - people should not fear being observed, measured and analysed.
3. **Anonymity** - which ensures that people will not be subjected to sanctions or reprimands.
4. **Safety** - protection of voters' data and communications.
5. **Decentralization** - data about campaigns and candidates should not be subject to censorship.
6. **Internationalization** - language should not be an obstacle to voter participation.
7. **Representativeness** - technology should help to bring closer voters to candidates.
8. **Simplicity** - everyone should be able to vote and every organization should be able to create a campaign.
9. **Scalability** - which ensures that campaigns will be processed efficiently.
10. **Flexibility** - there is no ideal format for a campaign, the code must be open to innovation.

## Table of Contents
1. [About us](#about-us)
2. [Our 10 principles](#our-10-principles)
3. [Motivation](#motivation)
4. [How it works](#how-it-works)
    1. [Structure of a campaign](#structure-of-a-campaign)
    2. [Campaign Creation](#campaign-creation)
    3. [Campaign data stored on IPFS](#campaign-data-stored-on-ipfs)
    4. [Campaign data stored in the blockchain](#campaign-data-stored-in-the-blockchain)
    5. [Candidate registration](#candidate-registration)
    6. [Candidate profile](#candidate-profile)
    7. [Group administrators](#group-administrators)
    8. [Group administrators registration](#group-administrators-registration)
    9. [Voter registration](#voter-registration)
    10. [Giving the right to vote](#giving-the-right-to-vote)
    11. [Choosing a group](#choosing-a-group)
    12. [Sending votes](#sending-votes)
    13. [Vote verification](#vote-verification)
    14. [Vote confirmation](#vote-confirmation)
    15. [Vote allocation](#vote-allocation)

## Motivation

On May 12, 2018, the Iraqis will choose their next representatives in what will be the fourth multiparty election after 2003. Fortunately, they are realizing that resolving conflicts through politics is the best way to avoid violence and chaos.

Meanwhile, in the West, people are **losing the faith in democracy**. As stated in this [article](https://www.economist.com/news/essays/21596796-democracy-was-most-successful-political-idea-20th-century-why-has-it-run-trouble-and-what-can-be-do):
> Faith in democracy flares up in moments of triumph, such as the overthrow of unpopular regimes in Cairo or Kiev, only to sputter out once again. Outside the West, democracy often advances only to collapse. And within the West, democracy has too often become associated with debt and dysfunction at home and overreach abroad. Democracy has always had its critics, but now old doubts are being treated with renewed respect as the weaknesses of democracy in its Western strongholds, and the fragility of its influence elsewhere, have become increasingly apparent. Why has democracy lost its forward momentum?

We think that one of the reasons for this is the **lack of innovation**.

In the old days, in Athens, democracy was identified with sortition, a kind of lottery in which people were selected randomly to occupy public positions. The American Revolution brought us the modern representative democracy in the end of the 18th century. Little has changed since then in relation to the _algorithm_ behind the selection of government officials.

There are differences between electoral systems. Some of them are comprehensible, for example, the two-round system in multiparty democracies. Other features of the electoral systems seem a bit arbitrary, such as the methods for allocating seats in parliament. Few of them, however, reflect the technological advances of the last decades.

As you know, technology has changed the way we _communicate_, the way we _learn_ and _work_, but the way we _vote_ remains virtually unchanged. Although, in some countries, people can vote from home, this is far from taking full advantage of widespread technologies. For most people, voting is a tedious and _pro forma_ task.

Before Wikipedia, people would hardly believe that accurate and complete information could emerge from the cooperation of thousands of anonymous individuals on the Internet [(link)](http://www.businessinsider.com/henry-blodget-sorry-britannica-wikipedias-not-only-bigger-but-better-2009-4). Nowadays, we can find disbelief in the capacity of democracy to well organize society. This disbelief is common in countries like China, where its political and economic elite claims for more centralization and uniformity. [(link)](https://www.ft.com/content/c4df31cc-4d26-11e8-97e4-13afc22d86d4)  How can we show them the merits of our system when, even among us, there are serious concerns about the dysfunctionality of it?

Democracy has an **advantage** compared to centralized regimes. People like to know they have influence on the future of their communities and nations. Therefore, it is really sad when citizens, in democratic countries, feel unrepresented and ignored by their politicians.

We need **open-source platforms** to spillover innovation across democracies and to face the criticisms democracy rightly receives. Centralized regimes are investing in artificial intelligence, big data, and all sort of technologies available. **_Will democracies be stuck in the past?_**

Let's use technology to _increase interest_ in politics among citizens, to _encourage participation_ in elections, and to _favor cooperation_ over conflict. We believe that all this is possible. The Kantcoin Project was created to champion **transparency** and **representativeness** in politics. We hope our contribution will help democracy thrive.

## How it works

### Structure of a campaign 

![Kantcoin network infochart](https://raw.githubusercontent.com/kantcoin/Kantcoin/master/infochart1.jpg)

### Campaign Creation

When a campaign is created, a new, exclusive blockchain is created.  This means that all users who connect to a campaign have a copy of all transactions (votes, confirmations, etc) that occur within that campaign.

In addition, each campaign have an IPNS address where relevant information is be stored, such as the enodes needed to connect to the blockchain.

### Campaign data stored on IPFS

IPFS is being used to store data that is not part of the voting process. For example, the description of the campaign, the enodes of group chairpersons, and lists of candidates and parties.

### Campaign data stored in the blockchain

Data stored in the blockchain are basically composed of votes, confirmations, candidates' vote balances, groups' composition, and voters' data (e.g. public key).

### Candidate registration

Candidates are registered by the creator of the campaign, who inserts candidates' data into the candidates IPFS page. For each round, new candidates must be inserted.

For example, in a given campaign, ten candidates compete in the first round. In the second round, only two of them compete.

A campaign may have up to five voting rounds.

### Candidate profile

Candidates may use a WYSIWYG editor to create their profiles. They can insert tables, links, images, and they can also change the font and style of the text with this editor. 

### Group administrators

Group administrators are Kantcoin network nodes who check and register votes in blockchain. They exist to make campaigns scalable. Since we need to wait a few seconds to check each vote, performing all the checks of a ten thousand voters campaign, with only one computer, could take a whole day. Therefore, in large campaigns, these votes should be processed diffusely. 

### Group administrators registration

Group administrators have to be listening to messages sent via Tor in order to receive votes. In addition, the creator of the campaign must attribute groups to group administrators so that voters will know to whom send their votes.

Another role that these group administrators perform is to allow voters to be able to access blockchain data via Tor. For this reason, their onion addresses must be informed to voters by the creator of the campaign beforehand.

### Voter registration

Most of e-voting systems today rely on authentication by verifying user's face and id photos. However, currently we authenticate voters with their cellphone numbers or email accounts. 
In order to make this work, we have created the role of the "login provider", a service that is external to the open-source project (kantcoin.org). In our case (kantcoin.com), a service that uses Amazon's Lambda and DynamoDB. 
The idea is that the creator of the campaign can freely choose the login provider. For example, he or she could use as login provider a service offered by their country's electoral authority, allowing voters to log in (which in our project means to create a "vault" whose public key (address) is stored in the login provider's database) with their electoral document.

### Giving the right to vote

Authorization to participate in an election campaign occurs after verifying the voter's username.

To make this work, it is necessary that the validation logic be embedded in the username, so that we can check, with a regular expression, if a user is able to participate in a campaign. 

To make this clearer, let's give some examples (regular expression - user groups that can participate in the campaign):

Telephones with the prefix 55-61 - residents of Brasília.
Emails ending with @unb.br - staff, students and professors of this university.

Besides validating the username with a regular expression, the creator of the campaign must check if the voter has the private key that is associated with that username in a login provider. This is done by checking the signature of a standardized message. If everything is OK, the campaign creator will grant the user the right to vote.

### Choosing a group

Before the voter can vote, he or she must choose a group. Currently, groups can consist of up to 15 members. Raising this number would lead to a more time-consuming process for validating ring signatures. Within each group, its members have the same weight. However, in donation rounds, the donated amount can vary between groups.

Groups may represent distinct social groups or may be anonymous. Voters can choose whichever group they want. However, voters can only participate in groups whose donation value is the same as theirs.

Examples of groups:
Women (donation value: 100 finney)
Men (donation value: 120 finney)

The value of each group, as well as its description, is defined by the creator of the campaign. Values are expressed in finney (one thousandth of ether)

### Sending votes

Votes are sent, via the Tor network, to group administrators. Votes are then stored locally, reordered and sent all at once to the blockchain, so that it is not possible to know when each was sent to the group administrator.

### Vote verification

Each vote is associated with an unique ring signature (URS). This signature allows group administrators to verify if a vote comes from a group, but it does not allow the identification of the author of that signature.

At each round, there is a standard message to be sent by all voters. If the same voter sends a vote more than once, it will be easily detected, so that we can guarantee the uniqueness of the voter.

### Vote confirmation

In the previous step, we can not guarantee if the group administrator is storing the votes for the correct candidates. For this reason, there is a confirmation step.

If a group adminstrator attempts to defraud a vote, the voter will be able to send an error message to the blockchain. This error message will block the allocation of votes and donations to candidates within that group.

The confirmation step occurs after the round (ballot) is stopped by the creator of the campaign.

### Vote allocation

In order for the votes to be assigned to candidates, the creator of the campaign should first close the round. Then, group administrators should call the "send votes", "pre-commit", and "commit" methods.

### Vote cancellation

An important concern about e-voting schemes is the possibility of coercion and vote-selling. According to this [article](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/11/rvc-jets0101.pdf):
> It is more than a little surprising to think that the best defense against coercion today may not
be found in an in-person voting system. Instead, the approach that offers the most resistance to
coercion may be a remote voting system using the paradigm of in which voters
who have successfully registered without coercion can convincingly pretend to vote according to
a subsequent coercers wishes while secretly voting their own true preferences. It may be possible
to leverage this approach by allowing voters who have been coerced during registration to quietly
invalidate their coerced credentials at some later opportunity and to receive new valid credentials.
For this reason, we have implemented a way to secretly, outside the blockchain, cancel votes. While the blockchain and URS are intended to protect voters against the government, this step is intended to protect voters against society.

