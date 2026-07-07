package main

import (
	"fmt"
	"os"
	"sync"
)

type KVPair struct {
	key   string
	value string
}

func main() {
	cities := []string{"Tokyo", "Delhi", "Shanghai", "Sao_Paulo", "Mexico_City", "Cairo", "Mumbai", "Beijing", "Dhaka", "Osaka", "New_York", "Karachi", "Buenos_Aires", "Istanbul", "Kolkata", "Lagos", "Moscow", "London", "Paris", "Los_Angeles"}

	//Reads in input file for each city and calls Map function
	ch := make(chan KVPair)
	var wg sync.WaitGroup
	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			path := "data/" + city + ".txt"
			input, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			map_out := Map(path, string(input))
			for _, item := range map_out {
				ch <- item
			}
		}(city)
	}

	// Goroutine waits to close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	// range over channel repeatedly reads from channel until it is closed
	kv_pairs := make(map[string][]string)
	for item := range ch {
		// Group every temperature (item.value) under its year (item.key).
		kv_pairs[item.key] = append(kv_pairs[item.key], item.value)
	}

	// Call Reduce on each year (in parallel) and print the max temperature.
	var rwg sync.WaitGroup
	for year, temps := range kv_pairs {
		rwg.Add(1)
		go func(year string, temps []string) {
			defer rwg.Done()
			max := Reduce(year, temps)
			fmt.Printf("Highest Temp in %s was %v\n", year, max)
		}(year, temps)
	}
	rwg.Wait()
}
