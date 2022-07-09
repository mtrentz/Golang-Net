package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// func copy(dst io.Writer, src io.Reader) {
// 	_, err := io.Copy(dst, src)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func handleConn(c net.Conn) {
	fmt.Println("Accepted connection from: ", c.RemoteAddr())
	defer c.Close()
	// time.Sleep(1 * time.Second)

	// Read from channel, prints, sends back to client
	// Send the reader to background
	// go func() {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Received: ", line)

		// Return the msg
		io.WriteString(c, "From server: "+line+"\n")
	}
	// }()

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
