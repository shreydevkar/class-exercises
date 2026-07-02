
package main

import (
    "log"
    "net/rpc"
)

type Move struct {
	Color int
    Col int
}

type Board struct {
	BoardString string
}

func main() {
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    // Synchronous call
    var reply Board
	var args int
    err = client.Call("ConnectGame.Get", &args, &reply)
    if err != nil {
        log.Fatal("game error:", err)
    }
   	log.Printf("Game: %v", reply)
}