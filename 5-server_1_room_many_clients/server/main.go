package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type Chat struct {
	n        int
	messages []string
}

func (chat *Chat) insertMessage(msg string) {
	c.L.Lock()
	chat.messages = append(chat.messages, msg)
	chat.n += 1
	c.Broadcast()
	c.L.Unlock()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// When connection is open, I'll store the amount of msgs
	// the Chat had when the user first connected
	n := chat.n

	go func() {
		// Keeps reading from the connection
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			// Case got a new line, add it to chat
			chat.insertMessage(line)
		}
	}()

	for {
		// Now I want to send all messages
		// to client every time that the Chat changes
		c.L.Lock()
		for chat.n == n {
			// Wait for the Chat to change
			c.Wait()
		}

		// Send all messages to client as a single string.
		// Separate them by 3 pipes |
		msgs := strings.Join(chat.messages, "|||")

		io.WriteString(conn, msgs+"\n")

		// Update the amount of msgs the Chat had
		n = chat.n

		c.L.Unlock()
	}
}

var chat Chat
var mu sync.Mutex
var c *sync.Cond

func main() {

	chat = Chat{}
	mu = sync.Mutex{}
	c = sync.NewCond(&mu)

	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
