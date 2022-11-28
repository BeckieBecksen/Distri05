# A Distributed Auction System

This is a distributed auction system that uses replication, and is resilient to single server failure.

## How to start the program
To get an auction started, use the terminal to navigate to the /Auctioneer/ folder, then use "go run Auctioneer.go 0" to launch the first auctioneer.
Repeat this with 2 other terminals and increase the number by 1 for each of the terminals.

Now you have your auctioneers ready for action, but no buyers. Navigate to the /Aristocrat/ folder, then use "go run Aristocrat.go 0" to launch the Aristocrat.
you can repeat this while increasing the interger to create more Aristocrats, but it's not necessary.

## Commands
The Auctioneers has no implemented commands, they simply stand around and wait for the Aristocrat(s) to make up their mind.
but they can be forcefully shut down using "ctrl + c" up to 2 of the 3 Auctioneers can be shut down without affecting the auction.

The Aristocrats have the following commands:
- use "bid_" followed by a positive integer to place a bid at the auction (the Auction has a timer of 1 minute that will begin upon the first bid being placed)
- use "status_" to see the current winning bid, or if the auction is over (when the auction is over no more bids can be placed)

to start a new Auction repeat the steps in "How to start the program"

## Authors

* Christopher Ryder Cortsen crco@itu.dk
* Emil Boesgaard Nørbjerg emno@itu.dk
* Ida Barkou Vilstrup ivil@itu.dk
* Rebecca Due Mylenberg remy@itu.dk

