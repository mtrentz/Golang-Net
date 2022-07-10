package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// I'll try to do it sync.Cond again

type Room struct {
	name     string
	cond     *sync.Cond
	messages []string
	n        int // amount of msgs
}

func (room *Room) insertMessage(msg string) {
	room.cond.L.Lock()
	room.messages = append(room.messages, msg)
	room.n += 1
	room.cond.Broadcast()
	room.cond.L.Unlock()
}

func getRoom(name string) *Room {
	return &Room{
		name: name,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func simulateUser(username string, room *Room) {
	n := room.n
	for {
		room.cond.L.Lock()
		for room.n == n {
			fmt.Println("User", username, "with room", room.name, "is waiting for messages")
			room.cond.Wait()
		}

		// Get all msgs from room separated by string
		msgs := strings.Join(room.messages, "|||")

		// Print them with username and room name
		fmt.Println("User", username, "with room", room.name, "has messages:", msgs)

		// Update the amount of msgs the Room had
		n = room.n

		room.cond.L.Unlock()

	}
}

func main() {
	// Start a channel of open interface to stop my program from ending
	c := make(chan interface{})

	r1 := getRoom("room1")
	r2 := getRoom("room2")

	// Simulate one user for each room
	go simulateUser("user1", r1)
	go simulateUser("user2", r2)

	// Simulate msg for room1 every 2 seconds
	go func() {
		for {
			r1.insertMessage("Msg to room1")
			time.Sleep(2 * time.Second)
		}
	}()

	// Simulate msg for room2 every 3 seconds
	go func() {
		for {
			r2.insertMessage("Msg to room2")
			time.Sleep(3 * time.Second)
		}
	}()

	// Prevent the program from ending
	<-c
}
