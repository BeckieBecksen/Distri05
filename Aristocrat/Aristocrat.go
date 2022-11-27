package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	Auction "github.com/BeckieBecksen/Distri05/Auction"
	"google.golang.org/grpc"
)

func main() {
	//Selects the port for each user starting at 5000 with the argument 0
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &bid{
		id: ownPort,
		//map of Auctioneers
		auctioneers: make(map[int32]Auction.PingServer),
		ctx:         ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	grpcServer := grpc.NewServer()
	Auction.RegisterPingServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	//connects to all clients except self (maximum 3 clients)
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)
		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := Auction.NewPingClient(conn)
		p.auctioneers[port] = c
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		Sendbid(1)
	}
}

func Sendbid(bid int32) {

}

type bid struct {
	Auction.UnimplementedPingServer
	id          int32
	auctioneers map[int32]Auction.PingServer
	ctx         context.Context
	lampTime    int32
	bidAmount   int32
}
