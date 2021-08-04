/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-03 19:25
**/

package app

import (
	"github.com/hashicorp/raft"
	"github.com/lemoyxk/console"
	"github.com/lemoyxk/kitty/socket"

	"discover/store"
)

// YOU GOT LEADER
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

	for conn := range Node.Server.GetConnections() {
		var err = conn.Close()
		if err != nil {
			console.Error(err)
		}
	}

	Node.Alive.DestroyConn()
	Node.Alive.DestroyData()
	Node.Register.Destroy()
	Node.Listen.Destroy()

	console.Error("local addr:", Node.Addr.Raft, "leader addr:", leader.Leader, Node.IsMaster())
}

func NewLeader(leader raft.LeaderObservation) {

	// delete all register
	if Node.IsMaster() {

	}

	console.Warning("local addr:", Node.Addr.Raft, "leader addr:", leader.Leader, Node.IsMaster())
}

// YOU GOT PEER CHANGE
func (n *node) OnPeerChange(peer raft.PeerObservation) {
	console.Warning("peer:", peer.Peer, "remove:", peer.Removed)
}

// YOU GOT KEY CHANGE
func (n *node) OnKeyChange(op *store.Command) {
	Node.Lock()
	defer Node.Unlock()

	var value, err = Node.Store.Get(op.Key)
	if err != nil {
		console.Error(err)
		return
	}

	var connections = Node.Listen.Get(op.Key)
	for i := 0; i < len(connections); i++ {
		var err = connections[i].Emit(socket.Pack{
			Event: "/OnListen",
			Data:  []byte(value),
		})
		if err != nil {
			console.Error(err)
		}
	}

	console.Info(op.Op, op.Key, op.Value)
}
