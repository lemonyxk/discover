/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-03 19:25
**/

package app

import (
	"github.com/hashicorp/raft"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

// OnLeaderChange YOU GOT LEADER
// when lose leader, the addr is empty until there is new master
func (n *node) OnLeaderChange(leader raft.LeaderObservation) {

	// LOSE LEADER
	if leader.Leader == "" {
		LoseLeader(leader)
	}

	// GOT NEW MASTER
	if leader.Leader != "" {
		NewLeader(leader)
	}
}

func LoseLeader(leader raft.LeaderObservation) {

	// CLOSE all client
	if Node.Server == nil {
		return
	}

	Node.Server.Range(func(conn server.Conn) {
		var err = conn.Close()
		if err != nil {
			console.Error.Logf("node close conn error: %s", err)
		}
	})

	Node.Alive.DestroyConn()
	Node.Alive.DestroyData()
	Node.Register.Destroy()
	Node.Key.Destroy()

	console.Warn.Logf("LoseLeader addr: %s leader: %s master: %v", Node.Addr.Addr, leader.Leader, Node.IsMaster())
}

func NewLeader(leader raft.LeaderObservation) {

	// delete all register
	if Node.IsMaster() {

	}

	console.Warn.Logf("NewLeader addr: %s leader: %s master: %v", Node.Addr.Addr, leader.Leader, Node.IsMaster())
}

// OnPeerChange YOU GOT PEER CHANGE
func (n *node) OnPeerChange(peer raft.PeerObservation) {
	console.Warn.Logf("OnPeerChange peer: %s remove: %v", peer.Peer, peer.Removed)
}

// OnKeyChange YOU GOT KEY CHANGE
func (n *node) OnKeyChange(op *store.Message) {
	Node.Lock()
	defer Node.Unlock()

	var connections = Node.Key.Get(op.Key)
	for i := 0; i < len(connections); i++ {
		var sender = socket.NewSender(connections[i].Conn())
		sender.SetCode(200)
		var err = sender.Emit("/Key", store.Build(op))
		if err != nil {
			console.Error.Logf("OnKeyChange error: %s", err)
		}
	}

	console.Info.Logf("OnKeyChange %v %s %d", op.Op, op.Key, len(op.Value))
}
