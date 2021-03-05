/**
* @program: discover
*
* @description:
*
* @author: lemo
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
	"github.com/lemoyxk/console"
	"github.com/lemoyxk/exception"
	"github.com/lemoyxk/kitty"
	client2 "github.com/lemoyxk/kitty/socket/udp/client"
	"github.com/lemoyxk/kitty/socket/websocket/server"
	"github.com/lemoyxk/utils"

	"discover/store"
	"discover/structs"
)

var Node = &node{
	StartTime: time.Now(),
}

type node struct {
	Store     *store.Store
	ServerMap *serverMap
	Addr      *structs.Address
	Config    *Config
	Client    *client2.Client
	Server    *server.Server
	Register  *register
	Alive     *alive
	Listen    *listen

	// start time
	StartTime time.Time

	lock sync.Mutex
}

func (n *node) GetMaster() structs.Address {
	return RaftAddr2Addr(string(n.Store.Raft().Leader()))
}

func (n *node) IsReady() bool {
	return n.Store.Raft().State() == raft.Leader || n.Store.Raft().State() == raft.Follower
}

func (n *node) IsMaster() bool {
	return n.Store.Raft().State() == raft.Leader
}

func (n *node) InitRegister() {
	n.Register = &register{data: make(map[int64]*structs.Register)}
}

func (n *node) InitAlive() {
	n.Alive = &alive{data: make(map[string][]structs.ServerInfo), conn: make(map[string][]*server.Conn)}
}

func (n *node) InitListen() {
	n.Listen = &listen{conn: make(map[string][]*server.Conn)}
}

func (n *node) InitStore() {
	n.Store = store.New(Node.Config.Dir, Node.Addr.Raft)
	n.Store.OnKeyChange = n.OnKeyChange
	n.Store.OnLeaderChange = n.OnLeaderChange
	n.Store.OnPeerChange = n.OnPeerChange
	exception.AssertError(n.Store.Open())
}

func (n *node) InitServerMap() {
	n.ServerMap = &serverMap{
		servers: make(map[string]structs.WhoIsMaster),
	}
}

func (n *node) InitAddr() {

	addr, err := net.ResolveTCPAddr("tcp", n.Config.Addr)
	exception.AssertError(err)

	n.Addr = &structs.Address{
		Addr: addr.String(),
		Http: addr.String(),
		Raft: fmt.Sprintf("%s:%d", addr.IP, addr.Port+1000),
		Tcp:  fmt.Sprintf("%s:%d", addr.IP, addr.Port+2000),
		Udp:  udpAddress,
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

func (n *node) Join(masterAddr string, addr string) {

	if masterAddr != Node.Config.Addr {

		var isMasterRes = utils.HttpClient.Get(fmt.Sprintf("http://%s/IsMaster", masterAddr)).Query(kitty.M{"addr": addr}).Send()
		console.AssertError(isMasterRes.LastError())
		if isMasterRes.String() != "OK" {
			console.Warning(masterAddr, "is master:", isMasterRes.String())
			time.Sleep(time.Millisecond * 100)
			n.Join(masterAddr, addr)
			return
		}

		var joinRes = utils.HttpClient.
			Post(fmt.Sprintf("http://%s/Join", masterAddr)).
			Form(kitty.M{"addr": addr}).
			Send()
		console.AssertError(joinRes.LastError())
		if joinRes.String() != "OK" {
			console.Error(joinRes.String())
			return
		}
	}
}

func (n *node) GetServerList() []structs.WhoIsMaster {
	var servers = n.Store.Raft().GetConfiguration().Configuration().Servers
	var list []structs.WhoIsMaster
	for _, s := range servers {
		list = append(list, structs.WhoIsMaster{
			Addr:      RaftAddr2Addr(string(s.Address)),
			Timestamp: 0,
			IsMaster:  s.Address == n.Store.Raft().Leader(),
		})
	}
	return list
}

func (n *node) GetMasterAddr() structs.WhoIsMaster {
	var servers = n.Store.Raft().GetConfiguration().Configuration().Servers
	for _, s := range servers {
		if s.Address == n.Store.Raft().Leader() {
			return structs.WhoIsMaster{
				Addr:      RaftAddr2Addr(string(s.Address)),
				Timestamp: 0,
				IsMaster:  s.Address == n.Store.Raft().Leader(),
			}
		}
	}
	return structs.WhoIsMaster{}
}

func (n *node) Lock() {
	n.lock.Lock()
}

func (n *node) Unlock() {
	n.lock.Unlock()
}
