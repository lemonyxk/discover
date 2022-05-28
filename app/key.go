/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-03-04 10:50
**/

package app

import (
	"sync"

	"github.com/lemonyxk/kitty/v2/socket/websocket/server"
)

type key struct {
	mux  sync.Mutex
	conn map[string][]server.Conn
}

func (l *key) All() map[string][]server.Conn {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.conn
}

func (l *key) Get(key string) []server.Conn {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.conn[key]
	if !ok {
		return nil
	}
	return list
}

func (l *key) Add(key string, conn server.Conn) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.conn[key]
	if !ok {
		l.conn[key] = append(l.conn[key], conn)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].FD() == conn.FD() {
			return false
		}
	}

	l.conn[key] = append(l.conn[key], conn)

	return true
}

func (l *key) Delete(key string, conn server.Conn) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.conn[key]
	if !ok {
		return false
	}

	var index = -1
	for i := 0; i < len(list); i++ {
		if list[i].FD() == conn.FD() {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	list = append(list[0:index], list[index+1:]...)

	if len(list) == 0 {
		delete(l.conn, key)
		return true
	}

	// put back
	// notice
	l.conn[key] = list

	return true
}

func (l *key) Destroy() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.conn = make(map[string][]server.Conn)
}
