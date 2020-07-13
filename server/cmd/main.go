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

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func notify(c net.Conn, event string, message string) {
	msg := event + " " + message
	c.Write([]byte(msg + "\n"))
}

func notifyAll(id, event string, message string) {
	connectionMap.Range(func(key string, value models.Connection) bool {
		// make sure receipient != sender
		if key != id {
			if contains(value.Events, event) {
				notify(value.Socket, event, message)
				return false
			}
		}
		return true
	})
}

func registerConnection(id string, c net.Conn, events []string) error {
	connectionMap.Store(id, models.Connection{
		Socket: c,
		Events: events,
	})
	updateActive()
	return nil
}

func unregister(id string) {
	connectionMap.Delete(id)
	updateActive()
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
			unregister(id)
			break
		}
		words := strings.Split(strings.TrimSpace(strings.ToLower(cmd)), " ")
		switch words[0] {
		case "register":
			if len(words) == 3 {
				id = words[1]
				events := strings.Split(words[2], ",")
				e := registerConnection(id, c, events)
				if e != nil {
					c.Write([]byte(e.Error() + "\n"))
					break
				}
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
