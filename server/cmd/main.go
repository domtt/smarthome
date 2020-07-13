package main

import (
	"bufio"
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

func notify(c net.Conn, event string, message string) {
	msg := event + " " + message
	c.Write([]byte(msg + "\n"))
}

func notifyAll(id, event string, message string) {
	connectionMap.WithEvent(event, func(key string, value models.Connection) {
		// make sure receipient != sender
		if key != id {
			notify(value.Socket, event, message)
		}
	})
}

func updateActive() {
	notifyAll("", ":active", strings.Join(connectionMap.Keys(), ","))
}

func handleTcp(c net.Conn) {
	fmt.Printf("Connected to %s\n", c.RemoteAddr().String())
	id := ""

	for {
		cmd, e := bufio.NewReader(c).ReadString('\n')
		if e != nil {
			connectionMap.Unregister(id)
			updateActive()
			break
		}
		words := strings.Split(strings.TrimSpace(strings.ToLower(cmd)), " ")
		switch words[0] {
		case "register":
			if len(words) == 3 {
				id = words[1]
				events := strings.Split(words[2], ",")
				e := connectionMap.Register(id, c, events)
				if e != nil {
					c.Write([]byte(e.Error() + "\n"))
					break
				}
				updateActive()
			}
		case "trigger":
			message := ""
			if len(words) > 2 {
				message = words[2]
			}
			notifyAll(id, words[1], message)
			c.Write([]byte("triggered\n"))
		default:
			c.Write([]byte("command not found\n"))
		}
	}
}

func main() {
	tcpServer()
}
