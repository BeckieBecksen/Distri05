package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	gRPC "github.com/BeckieBecksen/Distri05/Auction"
	"google.golang.org/grpc"
)

var serverName = flag.String("name", "default", "Senders name")
var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal
var LTime = int32(0)

func main() {
	//print statement for forced crash (ctrl + c)
	cs := make(chan os.Signal)
	signal.Notify(cs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cs
		fmt.Println("***CRASH*** [" + time.Now().Local().Format(time.Stamp) + "]")
		os.Exit(1)
	}()

	//establishing port
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5400
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", ownPort))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	server := &Server{
		port:          fmt.Sprint(ownPort),
		aristocrats:   make(map[string]*gRPC.CommClient),
		LampTime:      0,
		WinningBidder: make(map[int32]int32, 1),
	}

	gRPC.RegisterCommServer(grpcServer, server) //Registers the server to the gRPC server.

	fmt.Println("The Auction House has opened! I mr." + server.port + " Will be of your service [" + time.Now().Local().Format(time.Stamp) + "]")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

type Server struct {
	gRPC.UnimplementedCommServer
	port          string                      // port of the Auctioneer
	LampTime      int64                       // the Lamport time of the server
	aristocrats   map[string]*gRPC.CommClient // map of aristocrats ports + connections
	WinningBidder map[int32]int32             //map of the winningbidder
}

var AuctionStatus = true

func AuctionStartTime() {
	//after 1 minute, end the auction
	time.AfterFunc(time.Duration(time.Minute)*1, AuctionEnd)

}

func AuctionEnd() {
	AuctionStatus = false
	fmt.Println("the Auction is over [" + time.Now().Local().Format(time.Stamp) + "]")
}

func (s *Server) Bid(ctx context.Context, req *gRPC.BidAmount) (*gRPC.Reply, error) {
	AuctioneerId, _ := strconv.ParseInt(s.port, 10, 32)
	LTime += req.Lamptime + 1
	if AuctionStatus {
		if len(s.WinningBidder) == 0 {
			//first bidder starts the timer of 1 minute
			AuctionStartTime()
			s.WinningBidder[req.Id] = req.Amount
			fmt.Println("The Auction has now begun! [" + time.Now().Local().Format(time.Stamp) + "]")
			fmt.Println(fmt.Sprint(req.Id) + " bids $" + fmt.Sprint(req.Amount) + " Accepted [" + time.Now().Local().Format(time.Stamp) + "]")
			return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!", LampTime: LTime, Id: int32(AuctioneerId)}, nil
		}
		for el := range s.WinningBidder {
			if s.WinningBidder[el] < req.Amount {
				//if the new bid is greater
				delete(s.WinningBidder, el)
				s.WinningBidder[req.Id] = req.Amount
				fmt.Println(fmt.Sprint(req.Id) + " bids $" + fmt.Sprint(req.Amount) + " Accepted [" + time.Now().Local().Format(time.Stamp) + "]")
				return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!", LampTime: LTime, Id: int32(AuctioneerId)}, nil
			} else {
				if el == req.Id {
					//if bidder already has the highest bid
					fmt.Println(fmt.Sprint(req.Id) + " bids $" + fmt.Sprint(req.Amount) + " Denied: already has highest bid $" + fmt.Sprint(s.WinningBidder[el]) + " [" + time.Now().Local().Format(time.Stamp) + "]")
					return &gRPC.Reply{Response: "You have already have the highest bid! at at whopping $" + strconv.Itoa(int(s.WinningBidder[el])), LampTime: int32(s.LampTime), Id: int32(AuctioneerId)}, nil
				}
				//if someone else has the highest bid
				fmt.Println(fmt.Sprint(req.Id) + "'s bid of $" + fmt.Sprint(req.Amount) + " Denied: winning someone else has bid $" + fmt.Sprint(s.WinningBidder[el]) + " [" + time.Now().Local().Format(time.Stamp) + "]")
				return &gRPC.Reply{Response: "Your bid was rejected, another Aristocrat currently has a higher bid", LampTime: int32(s.LampTime), Id: int32(AuctioneerId)}, nil
			}
		}
		//shouldn't happen
		return &gRPC.Reply{Response: "Something went wrong"}, nil
	}
	//Auction is over
	return &gRPC.Reply{Response: "The auction is over!", Id: int32(AuctioneerId)}, nil
}

func (s *Server) Message(ctx context.Context, reqStat *gRPC.Request) (*gRPC.CurrentStatus, error) {
	AuctioneerId, _ := strconv.ParseInt(s.port, 10, 32)
	LTime += reqStat.Lamptime + 1
	if AuctionStatus {
		if len(s.WinningBidder) == 0 {
			fmt.Println(fmt.Sprint(reqStat.Id) + " requested Auction status, The Auction has yet to start [" + time.Now().Local().Format(time.Stamp) + "]")
			return &gRPC.CurrentStatus{Comment: "The Auction begins when the first bid is placed, Noone has bid yet!", LampTime: LTime, Id: int32(AuctioneerId)}, nil
		}
		for el := range s.WinningBidder {
			if el == reqStat.Id {
				fmt.Println(fmt.Sprint(reqStat.Id) + " requested the Auction status, " + fmt.Sprint(reqStat.Id) + " is the highest bidder with $" + strconv.Itoa(int(s.WinningBidder[el])) + " [" + time.Now().Local().Format(time.Stamp) + "]")
				return &gRPC.CurrentStatus{Comment: "The Auction is still running! You have the highest bid with $" + strconv.Itoa(int(s.WinningBidder[el])), LampTime: LTime, Id: int32(AuctioneerId)}, nil
			}
			fmt.Println(fmt.Sprint(reqStat.Id) + " requested the Auction status, " + fmt.Sprint(el) + " is the highest bidder with $" + strconv.Itoa(int(s.WinningBidder[el])) + " [" + time.Now().Local().Format(time.Stamp) + "]")
			return &gRPC.CurrentStatus{Comment: "The Auction is still running! The Highest Bid is $" + strconv.Itoa(int(s.WinningBidder[el])), LampTime: LTime, Id: int32(AuctioneerId)}, nil
		}
	}
	for el := range s.WinningBidder {
		fmt.Println(fmt.Sprint(reqStat.Id) + " requested the Auction status, the auction has ended, the winner is " + strconv.Itoa(int(s.WinningBidder[el])) + " [" + time.Now().Local().Format(time.Stamp) + "]")
		return &gRPC.CurrentStatus{Comment: "The Auction is over! The new owner is " + strconv.Itoa(int(el)), LampTime: LTime, Id: int32(AuctioneerId)}, nil
	}

	return &gRPC.CurrentStatus{Comment: "The Auction is over! noone claimed the item!", LampTime: int32(s.LampTime), Id: int32(AuctioneerId)}, nil
}
