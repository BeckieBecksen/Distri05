package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
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
		aristocrats:   make(map[string]*gRPC.PingClient),
		LampTime:      0,
		WinningBidder: make(map[int32]int32, 1),
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
	aristocrats                  map[string]*gRPC.PingClient // map of streams
	WinningBidder                map[int32]int32             //map of the winningbidder
}

func join() {

}

func startAuction(minutes time.Duration) {
	//after some ammount of time, end the auction
	time.AfterFunc(minutes, AuctionEnd)
}

func AuctionEnd() {
	fmt.Println("The Auction is over")
}

func (s *Server) Ping(ctx context.Context, req *gRPC.Request) (*gRPC.Reply, error) {
	//if the
	for el := range s.WinningBidder {
		if _, ok := s.WinningBidder[el]; ok || s.WinningBidder[el] > req.Amount {
			delete(s.WinningBidder, el)
			s.WinningBidder[req.Id] = req.Amount
			return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!"}, nil
		} else {
			fmt.Println("Client " + string(req.Id) + "'s bid has been denied")
			return &gRPC.Reply{Response: "Your bid was rejected, another Aristocrat currently has a higher bid"}, nil
		}
	}
	return &gRPC.Reply{Response: "Your bid was accepted, you are the leading bidder!" + string(req.Id)}, nil
}

func auctionStatus() {
	//returns a message of the status of the auction
}
