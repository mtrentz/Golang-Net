package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func copy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

func handleConn(c net.Conn) {
	fmt.Println("Accepted connection from: ", c.RemoteAddr())
	defer c.Close()
	time.Sleep(1 * time.Second)

	// Read from connection c, store into variable.
	// It needs to be scanner by line if want to reply back
	// or else the connection with client would need to close
	// to get and EOF
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Received: ", line)

		// Return the line duplicated
		io.WriteString(c, line+line)
	}
}

func main() {
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
