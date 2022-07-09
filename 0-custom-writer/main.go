package main

import (
	"io"
	"os"
	"strings"
)

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// I want to create a class that will implement the Writer interface.
// The point of it will be to just print a image in reverse to the stdout.
type ReverseWriter struct {
	_ struct{} // empty struct to make it non-empty
}

func (rw ReverseWriter) Write(p []byte) (n int, err error) {
	io.Copy(os.Stdout, strings.NewReader(reverseString(string(p))))
	return len(p), nil
}

func main() {
	rw := ReverseWriter{}
	io.Copy(rw, strings.NewReader("\nAbcdefg\n"))
}
