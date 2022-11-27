package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	gRPC "github.com/BeckieBecksen/Distri05/Auction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "5000", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.CommClient      //the server
var ServerConn *grpc.ClientConn //the server connection
var LTime = int32(0)

func main() {
	//parse flag/arguments
	flag.Parse()

	//log to file instead of console
	//f := setLog()
	//defer f.Close()

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")
	ConnectToServer()
	defer ServerConn.Close()

	readInput()
}

// connect to server
func ConnectToServer() {

	//dial options
	//the server is not using TLS, so we use insecure credentials
	//(should be fine for local testing but not in the real world)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	//dial the server, with the flag "server", to get a connection to it
	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%s", *serverPort), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
	server = gRPC.NewCommClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
}

func readInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Auction!, Today we have a mystery item for sale, how much $$$ will thee bid?")
	fmt.Println("--------------------------------------------------------------------------------------------")

	//Infinite loop to listen for clients input.
	for {
		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		if !conReady(server) {
			log.Printf("Client %s: something was wrong with the connection to the server :(", *clientsName)
			continue
		}

		if len(input) >= 4 {
			if strings.ToLower(input[0:4]) == "bid_" {
				//use re for bid method
				re := regexp.MustCompile("[0-9]+")
				st := re.FindAllString(input, -1)
				if st != nil {
					s1 := string(st[0])
					num, err := strconv.ParseInt(s1, 10, 32)
					if err == nil {
						placeBid(int32(num))
					} else {
						fmt.Println("please input a valid $$ bid")
					}
				} else {
					fmt.Println("please input a valid $$ bid")
				}
			}

			if len(input) >= 7 {
				if strings.ToLower(input[0:7]) == "status_" {
					getStatus()
				}
			}
		}

	}

}

func placeBid(bidA int32) {

	var stId, _ = strconv.ParseInt(string(*clientsName), 10, 32)
	var myId = int32(stId)
	fmt.Println(myId)

	bid := gRPC.BidAmount{
		Id:        myId,
		bidAmount: bidA,
		Lamptime:  LTime,
	}

	//Make gRPC call to server with amount, and recieve acknowlegdement back.
	ack, err := server.Bid(context.Background(), &bid)
	//ack, err := server.updateWinningBid(context.Background(), &bid)

	if err != nil {
		log.Printf("Client %s: no response from the server, attempting to reconnect", *clientsName)
		log.Println(err)
	}
	fmt.Println("the server says " + ack.Response)
}

func getStatus() {

}

func conReady(s gRPC.CommClient) bool {
	return ServerConn.GetState().String() == "READY"
}
