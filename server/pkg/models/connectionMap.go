package models

import (
	"net"
	"sync"
)

type ConnectionMap struct {
	inner sync.Map
}

func (cs *ConnectionMap) Get(id string) *Connection {
	res, ok := cs.inner.Load(id)
	if !ok {
		return nil
	}
	return res.(*Connection)
}

func (cs *ConnectionMap) Keys() []string {
	keys := []string{}
	cs.inner.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (cs *ConnectionMap) forEach(fn func(key string, value Connection) bool) {
	cs.inner.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(Connection))
	})
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func (cs *ConnectionMap) WithEvent(event string, fn func(key string, value Connection)) {
	cs.forEach(func(key string, value Connection) bool {
		if contains(value.Events, event) {
			fn(key, value)
		}
		return true
	})
}

func (cs *ConnectionMap) Unregister(id string) {
	cs.inner.Delete(id)
}

func (cs *ConnectionMap) Register(id string, c net.Conn, events []string) {
	cs.inner.Store(id, Connection{
		Socket: c,
		Events: events,
	})
}
