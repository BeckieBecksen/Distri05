package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	gRPC "github.com/BeckieBecksen/Distri05/Auction"
	"google.golang.org/grpc"
)

var serverName = flag.String("name", "default", "Senders name")
var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal

func main() {

	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server := &Server{
		port:          *port,
		aristocrats:   make(map[string]*gRPC.CommClient),
		LampTime:      0,
		WinningBidder: make(map[int32]int32, 1),
	}

	gRPC.RegisterCommServer(grpcServer, server) //Registers the server to the gRPC server.
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

}

type Server struct {
	gRPC.UnimplementedCommServer
	port          string                      // Not required but useful if your server needs to know what port it's listening to
	LampTime      int64                       // the Lamport time of the server
	aristocrats   map[string]*gRPC.CommClient // map of streams
	WinningBidder map[int32]int32             //map of the winningbidder
}

func join() {
	//add the users to a map
}

var AuctionStatus = true

func AuctionTime(minutes time.Duration) {
	//after some ammount of time, end the auction
	time.AfterFunc(minutes, AuctionEnd)
}

func AuctionEnd() {
	AuctionStatus = false
}

func (s *Server) Bid(ctx context.Context, req *gRPC.BidAmount) (*gRPC.Reply, error) {
	if len(s.WinningBidder) == 0 {
		s.WinningBidder[req.Id] = req.Amount
		fmt.Println("Client %v\n's bid has been accepted", req.Id)
		return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!"}, nil
	}
	for el := range s.WinningBidder {
		if s.WinningBidder[el] < req.Amount {
			delete(s.WinningBidder, el)
			s.WinningBidder[req.Id] = req.Amount
			fmt.Println("Client %v\n's bid has been accepted", req.Id)
			return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!"}, nil
		} else {

			if el == req.Id {
				return &gRPC.Reply{Response: "You have already have the highest bid! at at whopping $" + strconv.Itoa(int(s.WinningBidder[el]))}, nil
			}
			fmt.Println("Client %v\n's bid has been denied", req.Id)
			return &gRPC.Reply{Response: "Your bid was rejected, another Aristocrat currently has a higher bid"}, nil
		}
	}
	return &gRPC.Reply{Response: "Something went wrong"}, nil
}

func (s *Server) Message(ctx context.Context, reqStat *gRPC.Request) (*gRPC.CurrentStatus, error) {
	//returns a message of the status of the auction
	if AuctionStatus {
		if len(s.WinningBidder) == 0 {
			return &gRPC.CurrentStatus{Comment: "The Auction is still running! Noone has bid yet!"}, nil
		}
		for el := range s.WinningBidder {
			if el == reqStat.Id {
				return &gRPC.CurrentStatus{Comment: "The Auction is still running! You have the highest bid with $" + strconv.Itoa(int(s.WinningBidder[el]))}, nil

			}
			return &gRPC.CurrentStatus{Comment: "The Auction is still running! The Highest Bid is $" + strconv.Itoa(int(s.WinningBidder[el]))}, nil
		}
	}
	for el := range s.WinningBidder {
		return &gRPC.CurrentStatus{Comment: "The Auction is over! The new owner is " + strconv.Itoa(int(el))}, nil
	}
	return &gRPC.CurrentStatus{Comment: "The Auction is over! noone vlaimed the item!"}, nil
}
