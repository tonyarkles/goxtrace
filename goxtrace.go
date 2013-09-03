package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

type GoxEngine struct {
	Quit chan bool
	Db   *GoxDb
}

func handleXTraceConnection(conn net.Conn, engine *GoxEngine) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	cutset := " "
	currentRecord := make(map[string]interface{})
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
			record := NewGoxtraceRecord(currentRecord)
			if record != nil {
				engine.Db.Write(record)
				Log("Completed record:", record)
			} else {
				Log("Missing X-Trace in record:", currentRecord)
			}
			currentRecord = make(map[string]interface{})
		} else {
			Log("Unparseable input:", text)
		}
	}
	if err := scanner.Err(); err != nil {
		Log("Error reading from socket:", err)
	}
}

func handleJsonConnection(conn net.Conn, engine *GoxEngine) {
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
		record := NewGoxtraceRecord(currentRecord)
		if record != nil {
			engine.Db.Write(record)
			Log("Completed record:", record)
		} else {
			Log("Missing X-Trace in record:", currentRecord)
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

type server func(net.Conn, *GoxEngine)

func runServer(binding string, listener server, engine *GoxEngine) {
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
			go listener(conn, engine)
		}
		engine.Quit <- true
	}()
}

func main() {
	engine := &GoxEngine{Quit: make(chan bool), Db: NewGoxDb("data.db3")}
	runServer(":4444", handleXTraceConnection, engine)
	runServer(":4445", handleJsonConnection, engine)
	<-engine.Quit
}
