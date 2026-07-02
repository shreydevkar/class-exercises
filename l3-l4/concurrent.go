package main

import (
	"log"
	"net/rpc"
)

type Move struct {
	Color int
	Col   int
}

type Board struct {
	BoardString string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	// A channel is a pipe between goroutines. Each white goroutine will
	// send one number into it: 1 if its move succeeded, 0 if it failed.
	ch := make(chan int)

	for i := range 10 {
		go func() {
			var reply int
			moveWhite := Move{0, i % 5}
			errWhite := client.Call("ConnectGame.Move", &moveWhite, &reply)
			if errWhite != nil {
				log.Println("RPC error:", errWhite)
				ch <- 0 // send: this move failed
			} else {
				ch <- 1 // send: this move succeeded
			}
		}()
	}

	// Receive 10 times. This BLOCKS until all 10 goroutines have sent,
	// so main cannot continue until every white move is done.
	sum := 0
	for range 10 {
		sum = sum + <-ch // receive one number and add it to the total
	}
	log.Println("Successful Moves:", sum)

	// Only now — after all whites finished — do we send the black move.
	var replyB int
	moveBlack := Move{1, 0}
	errBlack := client.Call("ConnectGame.Move", &moveBlack, &replyB)
	if errBlack != nil {
		log.Println("RPC error:", errBlack)
	}
}
