package main

import (
	"io"
	"log"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	_, err := io.WriteString(c, "Hello World\n")
	if err != nil {
		log.Println(err)
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
