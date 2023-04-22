/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-03-04 10:50
**/

package app

import (
	"sync"

	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

type key struct {
	mux     sync.Mutex
	senders map[string][]socket.Emitter[server.Conn]
}

func (l *key) All() map[string][]socket.Emitter[server.Conn] {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.senders
}

func (l *key) Get(key string) []socket.Emitter[server.Conn] {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.senders[key]
	if !ok {
		return nil
	}
	return list
}

func (l *key) Add(key string, sender socket.Emitter[server.Conn]) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.senders[key]
	if !ok {
		l.senders[key] = append(l.senders[key], sender)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].Conn().FD() == sender.Conn().FD() {
			return false
		}
	}

	l.senders[key] = append(l.senders[key], sender)

	return true
}

func (l *key) Delete(key string, fd int64) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	var list, ok = l.senders[key]
	if !ok {
		return false
	}

	var index = -1
	for i := 0; i < len(list); i++ {
		if list[i].Conn().FD() == fd {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	list = append(list[0:index], list[index+1:]...)

	if len(list) == 0 {
		delete(l.senders, key)
		return true
	}

	// put back
	// notice
	l.senders[key] = list

	return true
}

func (l *key) Destroy() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.senders = make(map[string][]socket.Emitter[server.Conn])
}
