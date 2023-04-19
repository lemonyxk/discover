/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-27 15:05
**/

package app

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/discover/structs"
	"github.com/lemonyxk/exception"
	client2 "github.com/lemonyxk/kitty/socket/udp/client"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

var Node = &node{
	StartTime: time.Now(),
}

type node struct {
	Store    *store.Store
	Addr     *message.Server
	Config   *Config
	Client   *client2.Client
	Server   *server.Server
	Register *register
	Alive    *alive
	Key      *key

	// start time
	StartTime time.Time

	lock sync.Mutex
}

func (n *node) GetMaster() *message.Server {
	addr, _ := n.Store.Raft().LeaderWithID()
	if addr == "" {
		return nil
	}
	return RaftAddr2Addr(string(addr))
}

func (n *node) IsReady() bool {

	var cfg = n.Store.Raft().GetConfiguration()
	if cfg.Error() != nil {
		return false
	}

	if len(cfg.Configuration().Servers) == 0 {
		return false
	}

	return n.Store.IsReady
}

func (n *node) IsMaster() bool {
	return n.Store.Raft().State() == raft.Leader
}

func (n *node) InitRegister() {
	n.Register = &register{data: make(map[int64]*structs.Register)}
}

func (n *node) InitAlive() {
	n.Alive = &alive{data: make(map[string][]*message.ServerInfo), conn: make(map[string][]server.Conn)}
}

func (n *node) InitListen() {
	n.Key = &key{conn: make(map[string][]server.Conn)}
}

func (n *node) InitStore() {
	n.Store = store.New(Node.Config.Dir, Node.Addr.Raft)
	n.Store.OnKeyChange = n.OnKeyChange
	n.Store.OnLeaderChange = n.OnLeaderChange
	n.Store.OnPeerChange = n.OnPeerChange
	exception.Assert.LastNil(n.Store.Open())
}

func (n *node) InitAddr() {

	addr, err := net.ResolveTCPAddr("tcp", n.Config.Addr)
	exception.Assert.LastNil(err)

	host, _, err := net.SplitHostPort(n.Config.Addr)
	exception.Assert.LastNil(err)

	var http = n.Config.Http
	if http == "" {
		http = fmt.Sprintf("%s:%d", host, addr.Port)
	}

	var rf = n.Config.Raft
	if rf == "" {
		rf = fmt.Sprintf("%s:%d", host, addr.Port+1000)
	}

	var tcp = n.Config.Tcp
	if tcp == "" {
		tcp = fmt.Sprintf("%s:%d", host, addr.Port+2000)
	}

	n.Addr = &message.Server{
		Addr: http,
		Http: http,
		Raft: rf,
		Tcp:  tcp,
	}
}

func (n *node) InitConfig(config *Config) {

	if config == nil {
		panic("config is empty")
	}

	if config.Addr == "" {
		panic("addr is empty")
	}

	if config.Dir == "" {
		panic("dir is empty")
	}

	n.Config = config
}

func (n *node) GetServerList() []*message.Address {
	var servers = n.Store.Raft().GetConfiguration().Configuration().Servers
	var list []*message.Address
	for _, s := range servers {
		var leader, _ = n.Store.Raft().LeaderWithID()
		list = append(list, &message.Address{
			Server: RaftAddr2Addr(string(s.Address)),
			Master: s.Address == leader,
		})
	}
	return list
}

func (n *node) GetMasterAddr() *message.Address {
	var servers = n.Store.Raft().GetConfiguration().Configuration().Servers
	for _, s := range servers {
		var leader, _ = n.Store.Raft().LeaderWithID()
		if s.Address == leader {
			return &message.Address{
				Server: RaftAddr2Addr(string(s.Address)),
				Master: true,
			}
		}
	}
	return nil
}

func (n *node) Lock() {
	n.lock.Lock()
}

func (n *node) Unlock() {
	n.lock.Unlock()
}
