package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	gRPC "github.com/BeckieBecksen/Distri05/Auction"
	"google.golang.org/grpc"
)

var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal

func main() {

	fmt.Printf("Attempts to create listener on port %s\n", *port)

	// Create listener tcp on given port or default port 5400
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		fmt.Printf("Failed to listen on port %s: %v", *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server := &Server{
		port:     *port,
		streams:  make(map[string]*gRPC.PingClient),
		LampTime: 0,
	}

	gRPC.RegisterPingServer(grpcServer, server) //Registers the server to the gRPC server.

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

}

type Server struct {
	gRPC.UnimplementedPingServer                             // You need this line if you have a server struct
	port                         string                      // Not required but useful if your server needs to know what port it's listening to
	LampTime                     int64                       // the Lamport time of the server
	streams                      map[string]*gRPC.PingClient // map of streams
	WinningBidder                map[int32]int32             //map of the winningbidder
}

func join() {

}

func startAuction() {
	//after some ammount of time, end the auction
}

func updateWinningBid() {
	//updates the current winning bidder
}

func auctionStatus() {
	//returns a message of the status of the auction
}
