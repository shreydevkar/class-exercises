package main

import (
    "log"
	"net"
    "net/rpc"
	"net/http"
)

type Move struct {
	Color int
    Col int
}

type Board struct {
	BoardString string
}

type ConnectGame int

func (t *ConnectGame) Move(args *Move, reply *int) error {
	return nil
}

func (t *ConnectGame) Get(args *int, reply *Board) error {
    reply.BoardString = "Hello World"
	return nil
}

func main() {
    cg := new(ConnectGame)
    rpc.Register(cg)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("listen error:", err)
    }
	log.Println("Serving on PORT 1234")
    http.Serve(l, nil)
}