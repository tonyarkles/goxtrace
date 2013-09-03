package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	cutset := " "
	for {
		for scanner.Scan() {
			text := scanner.Text()
			chunks := strings.SplitN(text, ":", 2)
			Log("Chunks:", len(chunks))
			if len(chunks) == 2 {
				key := strings.Trim(chunks[0], cutset)
				value := strings.Trim(chunks[1], cutset)
				Log("Key:", key, "Value:", value)
			} else if text == "" {
				Log("Completed record")
			} else {
				Log("Unparseable input:", text)
			}
		}
		if err := scanner.Err(); err != nil {
			Log("Error reading from socket:", err)
		}
	}
	conn.Write([]byte("Hello world!\n"))
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func runServer(binding string) {
	ln, err := net.Listen("tcp", binding)
	if err != nil {

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			Log("Error from Accept():", err)
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
