// Every executable Go Program should contain a package called main.
// This tells the Go compiler to compile the package into an executable
// program rather than a shared library.
package main

import (
	"fmt"
	"os"
	"time"
)

// searchForWord reads the file at filepath and prints the index of every
// occurrence of target within its contents.
func searchForWord(filepath string, target string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("error reading %v: %v\n", filepath, err)
		return
	}

	input := string(data)

	for i := 0; i+len(target) <= len(input); i++ {
		if input[i:i+len(target)] == target {
			fmt.Printf("found %v @ %v\n", target, i)
		}
	}
}

// The entry point of a Go program should be the main function of main package.
// When the executable is run, main() is automatically called.
func main() {
	fmt.Println("Hello World")

	// --- Step 1: searching a hard-coded string for "cat" ---
	// input := "There once was a cat named Barry. He was a very good cat. This cat lived in Boston. He loved doing Boston-related activities (that were good for cats). He walked the esplanade. He shopped on Newbury. He ate at Tatte. He sometimes even went to TD Garden. Did you know that cats are not allowed in TD Garden?"
	// target := "cat"
	// for i := 0; i+len(target) <= len(input); i++ {
	// 	if input[i:i+len(target)] == target {
	// 		fmt.Printf("found cat @ %v\n", i)
	// 	}
	// }

	// --- Final step: search the dictionary for two words at the same time ---
	go searchForWord("dictionary.txt", "fish")
	go searchForWord("dictionary.txt", "dog")

	// Give the goroutines time to finish before main exits. Removing this
	// line causes the program to exit before the goroutines print anything.
	time.Sleep(2 * time.Second)
}
