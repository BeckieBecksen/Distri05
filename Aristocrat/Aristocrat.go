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
	"time"

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

	fmt.Println(a.id)

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
							fmt.Println("please input a valid $$ bid [" + time.Now().Local().Format(time.Stamp) + "]")
						}
					} else {
						fmt.Println("please input a valid $$ bid [" + time.Now().Local().Format(time.Stamp) + "]")
					}
				} else {
					fmt.Println("please input a valid $$ bid [" + time.Now().Local().Format(time.Stamp) + "]")
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

	bid := gRPC.BidAmount{
		Id:       a.id,
		Amount:   bidA,
		Lamptime: a.LampTime,
	}

	//calls all serves and gets replies/errors
	for idNeer, neer := range a.Auctioneers {
		a.LampTime++
		ack, err := neer.Bid(a.ctx, &bid)
		if err != nil {
			fmt.Println("Auctioneer " + fmt.Sprint(idNeer) + " says ERROR [" + time.Now().Local().Format(time.Stamp) + "]")
			a.LampTime++
		} else {
			fmt.Println("Auctioneer " + strconv.Itoa(int(ack.Id)) + " says " + ack.Response + " [" + time.Now().Local().Format(time.Stamp) + "]")
			a.LampTime++
		}

	}

}

func (a *Aristocrat) getStatus() {

	status := gRPC.Request{
		Id:       a.id,
		Lamptime: a.LampTime,
	}

	for idNeer, neer := range a.Auctioneers {
		a.LampTime++
		ack, err := neer.Message(a.ctx, &status)
		if err != nil {
			fmt.Println("Auctioneer " + fmt.Sprint(idNeer) + " says ERROR [" + time.Now().Local().Format(time.Stamp) + "]")
			a.LampTime++
		} else {
			fmt.Println("Auctioneer " + strconv.Itoa(int(ack.Id)) + " says " + ack.Comment + " [" + time.Now().Local().Format(time.Stamp) + "]")
			a.LampTime++
		}

	}

}
