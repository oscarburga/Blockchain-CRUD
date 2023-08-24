# Simulated Blockchain - Simple Elections CRUD

This was a project for a Concurrent & Distributed Programming university course during my 2022-01 
semester. The goal of the project was to implement a simple blockchain from scratch to be used for any 
purpose (our decision). 

My group decided to implement a very basic CRUD website where you could register candidates and 
voters to prepare for elections. The "database" (a simulated one, its just arrays) would be stored 
in the blockchain, and every block would represent a "transaction" (CRUD operation) on the database.

The blockchain has a simple consensus-by-vote mechanic performed after every transaction, where they 
each vote to count matching transaction hashes and the blockchain with the most matching hashes is the 
one that gets finalized as the authentic one and replicated to the nodes that had a corrupted/obsolete 
blockchain. We also implement simple deletion of nodes from the network if they time out when trying 
to communicate with them.

## Project structure

I no longer have the original project files since I changed computers. I could only find a zipped copy of 
the final version, but I thought it was a cool little project and wanted to have it public on my GitHub,
so that is why this repository exists and has no commits.

* The folder "go" you'll find all the source code for the blockchain network. 
* The folder "TF-back-final" contains the source code for the back-end server of the website.
* The source code for the website's front end was made mostly by my teammates and can be found on one of 
their profiles. [Link](https://github.com/ChristianEspirituCueva/React-Front).