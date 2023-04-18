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
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

type alive struct {
	mux  sync.Mutex
	conn map[string][]server.Conn
	data map[string][]*message.ServerInfo
}

func (s *alive) AllConn() map[string][]server.Conn {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.conn
}

func (s *alive) GetConn(serverName string) []server.Conn {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.conn[serverName]
	if !ok {
		return nil
	}
	return list
}

func (s *alive) AddConn(serverName string, conn server.Conn) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	var list, ok = s.conn[serverName]
	if !ok {
		s.conn[serverName] = append(s.conn[serverName], conn)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].FD() == conn.FD() {
			return false
		}
	}

	s.conn[serverName] = append(s.conn[serverName], conn)

	return true
}

func (s *alive) DeleteConn(serverName string, conn server.Conn) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	var list, ok = s.conn[serverName]
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
		delete(s.conn, serverName)
		return true
	}

	// put back
	// notice
	s.conn[serverName] = list

	return true
}

func (s *alive) DestroyConn() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.conn = make(map[string][]server.Conn)
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

func (s *alive) AddData(serverName, addr string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	var info = message.ServerInfo{
		ServerName: serverName,
		Addr:       addr,
	}

	var list, ok = s.data[serverName]
	if !ok {
		s.data[serverName] = append(s.data[serverName], &info)
		return true
	}

	// already in here
	for i := 0; i < len(list); i++ {
		if list[i].Addr == addr {
			return false
		}
	}

	s.data[serverName] = append(s.data[serverName], &info)

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
