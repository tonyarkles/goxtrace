package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	conn.Write([]byte("Hello world!\n"))
	conn.Close()
}

func runServer(binding string) {
	ln, err := net.Listen("tcp", binding)
	if err != nil {

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	binding := ":4444"
	fmt.Printf("Starting server on %s\n", binding)
	runServer(binding)
}
