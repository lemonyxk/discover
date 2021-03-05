/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 14:25
**/

package app

import (
	"math"
	"sync"

	"discover/structs"
)

// NOTICE
// this module just for the first time
// when the system not initialize
// so it's not accurate

type serverMap struct {
	servers map[string]structs.WhoIsMaster
	mux     sync.Mutex
}

func (s *serverMap) Set(addr string, wim structs.WhoIsMaster) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.servers[addr] = wim
}

func (s *serverMap) Delete(addr string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.servers, addr)
}

func (s *serverMap) Get(addr string) structs.WhoIsMaster {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.servers[addr]
}

func (s *serverMap) All() map[string]structs.WhoIsMaster {
	return s.servers
}

func (s *serverMap) GetMaster() structs.WhoIsMaster {
	var master structs.WhoIsMaster
	var tm = int64(math.MaxInt64)
	for _, t := range s.servers {

		if t.IsMaster {
			return t
		}

		if t.Timestamp < tm {
			master = t
			tm = t.Timestamp
		}
	}

	return master
}
