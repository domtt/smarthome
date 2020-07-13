package models

import "sync"

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

func (cs *ConnectionMap) Range(fn func(key string, value Connection) bool) {
	cs.inner.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(Connection))
	})
}

func (cs *ConnectionMap) Store(key string, value Connection) {
	cs.inner.Store(key, value)
}

func (cs *ConnectionMap) Delete(key string) {
	cs.inner.Delete(key)
}
