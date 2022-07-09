package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func copy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareMsg(username string, message string) string {
	return username + ": " + message + "\n"
}

func main() {

	var msg string
	var username string

	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get user name
	fmt.Print("Enter your username: ")
	fmt.Scanln(&username)

	// Send to background something to keep reading the channel back
	go func() {
		for {
			// When found something, clear terminal, print the chat and get user msg
			// fmt.Print("\033[H\033[2J")
			// copy(os.Stdout, strings.NewReader("-----------"))
			copy(os.Stdout, conn)
		}
	}()

	// First message goes separately
	fmt.Print("Enter your message: ")
	fmt.Scanln(&msg)
	io.WriteString(conn, prepareMsg(username, msg))

	for {
		fmt.Scanln(&msg)
		io.WriteString(conn, prepareMsg(username, msg))
		// Clear terminal
		fmt.Print("\033[H\033[2J")
	}

	// io.WriteString(conn, "Salve salve do cliente"+"\n")

	// copy(os.Stdout, conn)
}
