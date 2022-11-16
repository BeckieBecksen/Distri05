package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	//FAULTY PACKAGE
	//"grpcChatServer/chatserver"
)

func main() {

	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" //default Port set to 5000 if PORT is not set in env
	}

	//init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	//REWRITE
	//cs := chatserver.ChatServer{}
	//chatserver.RegisterServicesServer(grpcserver,&cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}

}
