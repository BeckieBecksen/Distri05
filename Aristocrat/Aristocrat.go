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
)

var myPort int64

var ServerConn *grpc.ClientConn //the server connection

func main() {
	//parse flag/arguments
	flag.Parse()
	Port, _ := strconv.ParseInt(os.Args[1], 10, 32)
	myPort = Port

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")
	ConnectToServer()
	defer ServerConn.Close()

}

// connect to server
func ConnectToServer() {

	a := &Aristocrat{
		id:          int32(myPort) + 5000,
		Auctioneers: make(map[int32]gRPC.CommClient),
		ctx:         context.Background(),
		LampTime:    0,
	}
	fmt.Print(a.id)
	for i := 0; i < 3; i++ {
		port := int32(5400) + int32(i)

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := gRPC.NewCommClient(conn)
		a.Auctioneers[port] = c
	}

	a.readInput()
}

type Aristocrat struct {
	gRPC.UnimplementedCommServer
	id          int32
	Auctioneers map[int32]gRPC.CommClient
	ctx         context.Context
	LampTime    int32
}

func (a *Aristocrat) readInput() {
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

		if len(input) >= 4 {
			if strings.ToLower(input[0:4]) == "bid_" {
				//use re for bid method
				re := regexp.MustCompile("[0-9]+")
				st := re.FindAllString(input, -1)
				if st != nil {
					s1 := string(st[0])
					num, err := strconv.ParseInt(s1, 10, 32)
					if err == nil {
						if num > 0 {
							a.placeBid(int32(num))
						} else {
							fmt.Println("please input a valid $$ bid")
						}
					} else {
						fmt.Println("please input a valid $$ bid")
					}
				} else {
					fmt.Println("please input a valid $$ bid")
				}
			}

			if len(input) >= 7 {
				if strings.ToLower(input[0:]) == "status_" {

					a.getStatus()
				}
			}
		}

	}

}

func (a *Aristocrat) placeBid(bidA int32) {
	a.LampTime++
	bid := gRPC.BidAmount{
		Id:       a.id,
		Amount:   bidA,
		Lamptime: a.LampTime,
	}

	c := make(chan *gRPC.Reply)

	//calls all serves and gets first reply
	for _, neer := range a.Auctioneers {
		go checkBidResponse(a.ctx, &bid, c, neer)
	}

	firstresponse := <-c

	a.LampTime += firstresponse.LampTime + 1

	fmt.Println("Auctioneer " + strconv.Itoa(int(firstresponse.Id)) + " says " + firstresponse.Response)
}

func checkBidResponse(cx context.Context, b *gRPC.BidAmount, channel chan *gRPC.Reply, AuctioneerConn gRPC.CommClient) {
	ack, _ := AuctioneerConn.Bid(cx, b)
	if ack != nil {
		channel <- ack
	}
}

func checkStatusResponse(cx context.Context, r *gRPC.Request, channel2 chan *gRPC.CurrentStatus, AuctioneerConn gRPC.CommClient) {
	ack, _ := AuctioneerConn.Message(cx, r)
	if ack != nil {
		channel2 <- ack
	}
}

func (a *Aristocrat) getStatus() {
	var stId, _ = strconv.ParseInt(string(a.id), 10, 32)
	var myId = int32(stId)
	fmt.Println(myId)

	cha := make(chan *gRPC.CurrentStatus)

	status := gRPC.Request{
		Id:       myId,
		Lamptime: a.LampTime,
	}

	for _, neer := range a.Auctioneers {
		go checkStatusResponse(a.ctx, &status, cha, neer)
	}

	firstresp := <-cha
	fmt.Println("Auctioneer " + strconv.Itoa(int(firstresp.Id)) + " says " + firstresp.Comment)
	a.LampTime++
}
