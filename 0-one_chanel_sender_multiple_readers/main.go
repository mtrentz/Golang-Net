package main

import (
	"fmt"
	"time"
)

func reader(c <-chan int, id int) {
	for num := range c {
		fmt.Println("Reader", id, "received", num)
	}
}

func main() {
	c := make(chan int)
	// Objective here is to create three goroutines that
	// will listen on a channel and print
	for i := 0; i < 3; i++ {
		go func(i int) {
			for {
				reader(c, i)
			}
		}(i)
	}

	// Wait 1s for the goroutines to start
	time.Sleep(time.Second)

	// Then create one that will send a few messages to the channel
	go func() {
		for i := 0; i <= 10; i++ {
			c <- i
		}
	}()

	// Wait a few seconds
	time.Sleep(time.Second * 3)
}
