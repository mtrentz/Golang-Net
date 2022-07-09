package main

import (
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
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// io.WriteString(conn, "Salve salve do cliente"+"\n")

	// Keep reading from conn
	for {
		copy(os.Stdout, conn)
	}

	// copy(os.Stdout, conn)
}
