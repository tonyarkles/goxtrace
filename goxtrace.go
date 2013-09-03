package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

func handleXTraceConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	cutset := " "
	currentRecord := make(map[string]string)
	for scanner.Scan() {
		text := scanner.Text()
		chunks := strings.SplitN(text, ":", 2)
		Log("Chunks:", len(chunks))
		if len(chunks) == 2 {
			key := strings.Trim(chunks[0], cutset)
			value := strings.Trim(chunks[1], cutset)
			currentRecord[key] = value
			Log("Key:", key, "Value:", value)
		} else if text == "" {
			Log("Completed record:", currentRecord)
			currentRecord = make(map[string]string)
		} else {
			Log("Unparseable input:", text)
		}
	}
	if err := scanner.Err(); err != nil {
		Log("Error reading from socket:", err)
	}
}

func handleJsonConnection(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		var currentRecord map[string]interface{}
		if err := dec.Decode(&currentRecord); err == io.EOF {
			break
		} else if err != nil {
			Log("Error reading from connection:", err)
			break
		}
		Log("Completed record:", currentRecord)
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

type server func(net.Conn)

func runServer(binding string, listener server, quit chan bool) {
	ln, err := net.Listen("tcp", binding)
	if err != nil {
		Log("Listen error")
		return
	}
	fmt.Printf("Listening on %s\n", binding)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				Log("Error from Accept():", err)
				continue
			}
			go listener(conn)
		}
		quit <- true
	}()
}

func main() {
	quit := make(chan bool)
	runServer(":4444", handleXTraceConnection, quit)
	runServer(":4445", handleJsonConnection, quit)
	<-quit
}
