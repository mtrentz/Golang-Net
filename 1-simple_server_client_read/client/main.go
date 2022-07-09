package main

import (
	"io"
	"log"
	"net"
	"os"
)

func copy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	copy(os.Stdout, conn)

}
