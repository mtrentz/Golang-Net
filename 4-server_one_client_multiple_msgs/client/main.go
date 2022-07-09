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

func main() {

	var msg string

	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send to background something to keep reading the channel back
	go func() {
		for {
			copy(os.Stdout, conn)
		}
	}()

	for {
		fmt.Print("Enter a message: ")
		fmt.Scanln(&msg)
		io.WriteString(conn, msg+"\n")
	}

	// io.WriteString(conn, "Salve salve do cliente"+"\n")

	// copy(os.Stdout, conn)
}
