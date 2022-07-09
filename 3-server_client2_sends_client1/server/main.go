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

	// Keep reading lines from client and add to msgs
	go func() {
		for {
			scanner := bufio.NewScanner(c)
			for scanner.Scan() {
				line := scanner.Text()
				msgs <- line
			}
		}
	}()

	// Write msgs to client
	for m := range msgs {
		io.WriteString(c, m)
	}
	// go func() {

	// }()
}

var msgs chan string

func main() {

	msgs = make(chan string)

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
