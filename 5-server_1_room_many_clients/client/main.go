package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func prepareMsg(username string, message string) string {
	// Separated as a function to maybe add datetime later
	return username + ": " + message + "\n"
}

func main() {

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
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				line := scanner.Text()
				// Clears the terminal
				fmt.Print("\033[H\033[2J")
				fmt.Println(line)
				// I receive all msgs as a single string separated by "|||",
				// So I'll replace by a newline to print to terminal
				line = strings.Replace(line, "|||", "\n", -1)
				fmt.Println(line)
				// Prints the prompt, to stay at the bottom of terminal
				fmt.Print("Enter your message: ")
			}
		}
	}()

	// Loop to read from user input
	for {
		fmt.Print("Enter your message: ")
		inputScanner := bufio.NewScanner(os.Stdin)
		if inputScanner.Scan() {
			msg := inputScanner.Text()
			io.WriteString(conn, prepareMsg(username, msg))
		}
	}
}
