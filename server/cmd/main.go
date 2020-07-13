package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/td0m/smarthome/server/pkg/models"
)

var connectionMap models.ConnectionMap

func tcpServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic("failed to start the server")
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panicln("fuck")
		}
		go handleTcp(conn)
	}
}

type Message struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

func notify(c net.Conn, event string, payload interface{}) {
	msg := Message{
		event,
		payload,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("fuck")
		return
	}
	c.Write([]byte(string(msgBytes) + "\n"))
}

func notifyAll(id, event string, message interface{}) {
	connectionMap.WithEvent(event, func(key string, value models.Connection) {
		// make sure receipient != sender
		if key != id {
			notify(value.Socket, event, message)
		}
	})
}

func updateActive() {
	notifyAll("", ":active", connectionMap.Keys())
}

func handleTcp(c net.Conn) {
	fmt.Printf("Connected to %s\n", c.RemoteAddr().String())
	id := ""

	for {
		cmd, e := bufio.NewReader(c).ReadString('\n')
		if e != nil {
			connectionMap.Unregister(id)
			fmt.Println("disconnect " + c.RemoteAddr().String())
			updateActive()
			break
		}
		words := strings.Split(strings.TrimSpace(strings.ToLower(cmd)), " ")
		switch words[0] {
		case "register":
			if len(words) == 3 {
				id = words[1]
				events := strings.Split(words[2], ",")
				connectionMap.Register(id, c, events)
				updateActive()
			}
		case "trigger":
			var payload interface{}
			if len(words) > 2 {
				payloadStr := strings.Join(words[2:], " ")
				e := json.Unmarshal([]byte(payloadStr), &payload)
				if e != nil {
					fmt.Println(e)
				}
			}
			notifyAll(id, words[1], payload)
			//c.Write([]byte("triggered\n"))
		default:
			//c.Write([]byte("command not found\n"))
		}
	}
}

func main() {
	tcpServer()
}
