package models

import "net"

type Connection struct {
	Socket net.Conn
	Events []string
}
