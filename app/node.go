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
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/discover/structs"
	"github.com/lemonyxk/discover/utils"
	"github.com/lemonyxk/exception"
	"github.com/lemonyxk/kitty/socket"
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
	Client   *client2.Client[any]
	Server   *server.Server[any]
	Register *register
	Alive    *alive
	Key      *key

	// start time
	StartTime time.Time

	lock sync.Mutex
}

func (n *node) GetMaster() *message.Server {
	_, leader := n.Store.Raft().LeaderWithID()
	if leader == "" {
		return nil
	}
	return RaftAddr2Addr(string(leader))
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
	n.Alive = &alive{data: make(map[string][]*message.ServerInfo), senders: make(map[string][]socket.Emitter[server.Conn])}
}

func (n *node) InitListen() {
	n.Key = &key{senders: make(map[string][]socket.Emitter[server.Conn])}
}

func (n *node) InitStore() {
	n.Store = store.New(Node.Config.Dir, Node.Addr.Raft)
	n.Store.OnKeyChange = n.OnKeyChange
	n.Store.OnLeaderChange = n.OnLeaderChange
	n.Store.OnPeerChange = n.OnPeerChange
	exception.Assert.LastNil(n.Store.Open())
}

func (n *node) InitAddr() {

	host, port, err := utils.SplitHostPort(n.Config.Addr)
	exception.Assert.LastNil(err)

	var http = n.Config.Http
	if http == "" {
		http = fmt.Sprintf("%s:%d", host, port)
	}

	var rf = n.Config.Raft
	if rf == "" {
		rf = fmt.Sprintf("%s:%d", host, port+1000)
	}

	var tcp = n.Config.Tcp
	if tcp == "" {
		tcp = fmt.Sprintf("%s:%d", host, port+2000)
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
		var _, leader = n.Store.Raft().LeaderWithID()
		list = append(list, &message.Address{
			Server: RaftAddr2Addr(string(s.Address)),
			Master: string(s.Address) == string(leader),
		})
	}
	return list
}

func (n *node) GetMasterAddr() *message.Address {
	var servers = n.Store.Raft().GetConfiguration().Configuration().Servers
	for _, s := range servers {
		var _, leader = n.Store.Raft().LeaderWithID()
		if string(s.Address) == string(leader) {
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
