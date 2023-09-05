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

	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

type alive struct {
	mux     sync.Mutex
	senders map[string][]socket.Emitter[server.Conn]
	data    map[string][]*message.ServerInfo
}

func (s *alive) AllConn() map[string][]socket.Emitter[server.Conn] {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.senders
}

func (s *alive) GetConn(serverName string) []socket.Emitter[server.Conn] {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.senders[serverName]
	if !ok {
		return nil
	}
	return list
}

func (s *alive) AddConn(serverName string, sender socket.Emitter[server.Conn]) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.senders[serverName]
	if !ok {
		s.senders[serverName] = append(s.senders[serverName], sender)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].Conn().FD() == sender.Conn().FD() {
			return false
		}
	}

	s.senders[serverName] = append(s.senders[serverName], sender)

	return true
}

func (s *alive) DeleteConn(serverName string, fd int64) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	var list, ok = s.senders[serverName]
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
		delete(s.senders, serverName)
		return true
	}

	// put back
	// notice
	s.senders[serverName] = list

	return true
}

func (s *alive) DestroyConn() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.senders = make(map[string][]socket.Emitter[server.Conn])
}

func (s *alive) AllData() map[string][]*message.ServerInfo {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.data
}

func (s *alive) GetData(serverName string) []*message.ServerInfo {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.data[serverName]
	if !ok {
		return nil
	}
	return list
}

func (s *alive) AddData(info message.ServerInfo) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	var list, ok = s.data[info.Name]
	if !ok {
		s.data[info.Name] = append(s.data[info.Name], &info)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].Addr == info.Addr {
			return false
		}
	}

	s.data[info.Name] = append(s.data[info.Name], &info)

	return true
}

func (s *alive) DeleteData(serverName, addr string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.data[serverName]
	if !ok {
		return false
	}

	var index = -1
	for i := 0; i < len(list); i++ {
		if list[i].Addr == addr {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	list = append(list[0:index], list[index+1:]...)

	if len(list) == 0 {
		delete(s.data, serverName)
		return true
	}

	// put back
	// notice
	s.data[serverName] = list

	return true
}

func (s *alive) DestroyData() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data = map[string][]*message.ServerInfo{}
}
