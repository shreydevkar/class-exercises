package main

import (
	"fmt"
	"log"
	"maps"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type PeerServer struct {
	Address string
	Client  *rpc.Client
}

type Args struct {
	GossipLive map[string]int
	Round      int
	Sender     string
}

type Server struct {
	live    map[string]int
	lock    sync.Mutex
	Round   int
	Address string
	peers   []PeerServer
}

func (t *Server) Heartbeat(args *Args, reply *int) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// If the incoming round is ahead of ours, catch up
	if args.Round > t.Round {
		t.Round = args.Round
	}

	// Mark the sender as live as of this round
	t.live[args.Sender] = t.Round

	// Merge gossip: keep the most recent round we've heard for each server
	for addr, round := range args.GossipLive {
		if addr == t.Address {
			continue // no need to track ourselves
		}
		if round > t.live[addr] {
			t.live[addr] = round
		}
	}

	*reply = t.Round
	return nil
}

func (t *Server) sendHeartbeat(to PeerServer) {
	t.lock.Lock()
	t.Round++
	args := &Args{
		GossipLive: maps.Clone(t.live),
		Round:      t.Round,
		Sender:     t.Address,
	}
	t.lock.Unlock()

	var reply int
	err := to.Client.Call("Server.Heartbeat", args, &reply)
	if err != nil {
		log.Println("RPC error:", err)
	}
}

func (t *Server) GenerateReport() {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Println("REPORT!")
	log.Println("ROUND", t.Round)
	log.Println(t.live)
	fmt.Println("---- LIVE SERVERS ----")
	for addr, round := range t.live {
		if t.Round-round <= 10 {
			fmt.Println(addr)
		}
	}
}

func main() {

	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	my_address := "10.193.25.197:1234"
	server.Address = my_address
	server.Round = 0
	server.peers = make([]PeerServer, 0)
	server.live = make(map[string]int)
	peer_addresses := []string{
		"10.239.244.33:1234",
		"10.239.38.177:1234",
		"10.239.246.218:1234",
	}

	time.Sleep(10 * time.Second) // WAIT to start other servers

	for _, addr := range peer_addresses {
		if addr == my_address {
			continue
		}
		client, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		server.peers = append(server.peers, PeerServer{addr, client})
	}

	// Send a heartbeat to one random peer every second (non-blocking)
	go func() {
		for range time.Tick(time.Second) {
			if len(server.peers) == 0 {
				continue
			}
			peer := server.peers[rand.Intn(len(server.peers))]
			go server.sendHeartbeat(peer)
		}
	}()

	// Generate a report every 5 seconds
	for range time.Tick(5 * time.Second) {
		server.GenerateReport()
	}
}
