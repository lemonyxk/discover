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
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
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
			console.Error(err)
		}
	})

	Node.Alive.DestroyConn()
	Node.Alive.DestroyData()
	Node.Register.Destroy()
	Node.Key.Destroy()

	console.Warning("LoseLeader addr:", Node.Addr.Addr, "leader:", leader.LeaderID, "master:", Node.IsMaster())
}

func NewLeader(leader raft.LeaderObservation) {

	// delete all register
	if Node.IsMaster() {

	}

	console.Warning("NewLeader addr:", Node.Addr.Addr, "leader:", leader.LeaderID, "master:", Node.IsMaster())
}

// OnPeerChange YOU GOT PEER CHANGE
func (n *node) OnPeerChange(peer raft.PeerObservation) {
	console.Warning("OnPeerChange peer:", peer.Peer, "remove:", peer.Removed)
}

// OnKeyChange YOU GOT KEY CHANGE
func (n *node) OnKeyChange(op *store.Command) {
	Node.Lock()
	defer Node.Unlock()

	var value, err = Node.Store.Get(op.Key)
	if err != nil {
		console.Error(err)
		return
	}

	var connections = Node.Key.Get(op.Key)
	for i := 0; i < len(connections); i++ {
		var err = connections[i].JsonEmit("/Key", message.Format{
			Status: "SUCCESS", Code: 200, Msg: message.Op{
				Key: op.Key, Value: value, Op: op.Op,
			},
		})
		if err != nil {
			console.Error(err)
		}
	}

	console.Info("OnKeyChange", op.Op, op.Key, op.Value)
}
