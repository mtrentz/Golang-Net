package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	// Create a goroutine that clears the terminal
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// fmt.Print("\033[H\033[2J")
			io.Copy(os.Stdout, strings.NewReader("\033[H\033[2J"))
		}
	}()

	// Printing random stuff to terminal every once in a while
	for {
		// Pick from random assortment of symbols !@#$%&*
		symbols := []rune("!@#$%&*")
		// Pick a random index from the symbols array
		index := rand.Intn(len(symbols))
		// Pick a random symbol from the symbols array
		symbol := symbols[index]

		fmt.Println(strings.Repeat(string(symbol), 50))
		time.Sleep(time.Second * 1)
	}
}
