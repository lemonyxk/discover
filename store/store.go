// Package store provides a simple distributed key-value store. The keys and
// associated values are changed via distributed consensus, meaning that the
// values are changed only when a majority of nodes in the cluster agree on
// the new value.
//
// Distributed consensus is provided via the Raft algorithm, specifically the
// Hashicorp implementation.
package store

import (
	"fmt"
	json "github.com/lemonyxk/kitty/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/lemonyxk/console"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
	loseLeaderTimeout   = 15 * time.Second
)

// Store is a simple key-value store, where all changes are made via Raft consensus.
type Store struct {
	RaftDir  string
	RaftAddr string

	inMem bool // use mem or file

	mux  sync.Mutex
	data map[string][]byte // The key-value store for the system.

	raft *raft.Raft // The consensus mechanism

	onKeyChange func(op *Message)

	OnLeaderChange func(leader raft.LeaderObservation)
	OnPeerChange   func(leader raft.PeerObservation)
	OnKeyChange    func(op *Message)
	Ready          chan bool
	IsReady        bool
}

func (s *Store) Raft() *raft.Raft {
	return s.raft
}

func (s *Store) Shutdown() raft.Future {
	return s.raft.Shutdown()
}

// New returns a new Store.
func New(dataDir, raftAddr string) *Store {
	return &Store{
		RaftDir:  dataDir,
		RaftAddr: raftAddr,
		Ready:    make(chan bool, 1),
		data:     make(map[string][]byte),
		inMem:    false, // not support inMem
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (s *Store) Open() error {
	// create dir
	_ = os.MkdirAll(s.RaftDir, 0700)

	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(s.RaftAddr)

	var log = hclog.New(&hclog.LoggerOptions{
		Name:       "raft",
		Level:      hclog.NoLevel,
		Output:     config.LogOutput,
		JSONFormat: true,
	})

	config.Logger = log

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftAddr)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.RaftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	var allCount uint64
	if s.inMem {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		boltDB, err := NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
		if err != nil {
			return fmt.Errorf("new bolt store: %s", err)
		}
		logStore = boltDB
		stableStore = boltDB
		allCount = boltDB.Count()
		// first init
		if allCount == 0 {
			s.IsReady = true
		}
	}

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}

	var kCount uint64
	s.onKeyChange = func(op *Message) {
		// change event
		if s.OnKeyChange != nil && s.IsReady {
			s.OnKeyChange(op)
		}

		// finish all key
		kCount++
		if kCount == allCount {
			s.IsReady = true
		}
	}

	s.raft = ra

	ra.RegisterObserver(raft.NewObserver(nil, false, func(o *raft.Observation) bool {
		switch v := o.Data.(type) {
		case raft.LeaderObservation:
			if s.OnLeaderChange != nil {
				s.OnLeaderChange(v)
			}
		case raft.PeerObservation:
			if s.OnPeerChange != nil {
				s.OnPeerChange(v)
			}
		case raft.FailedHeartbeatObservation:
			// var sub = time.Now().Sub(v.LastContact)
			// if sub > loseLeaderTimeout {
			// 	// leave leader
			// 	var err = s.Leave(string(v.PeerID))
			// 	if err != nil {
			// 		console.Error(err)
			// 	} else {
			// 		console.Info("raft leave leader", v.PeerID, time.Now().Sub(v.LastContact))
			// 	}
			// } else {
			// 	console.Info("raft failed heartbeat", v.PeerID, time.Now().Sub(v.LastContact))
			// }
			console.Info("raft failed heartbeat", v.PeerID, time.Now().Sub(v.LastContact))
		case raft.ResumedHeartbeatObservation:
			console.Info("raft resumed heartbeat", v.PeerID)
		case raft.RequestVoteRequest:
			console.Info("raft request vote request", string(v.ID))
		case raft.RaftState:
			console.Info("raft state", v.String())
		default:
			console.Infof("raft other %+v\n", v)
		}
		return true
	}))

	// everything is ok
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if ra.State() == raft.Follower || ra.State() == raft.Leader {
				s.Ready <- true
				break
			}
		}
	}()

	console.Info("raft server start at", s.RaftAddr, "state", ra.State())

	return nil
}

// BootstrapCluster MUST SET ONE
func (s *Store) BootstrapCluster(ok bool) {
	if !ok {
		return
	}

	// when add new server
	// MAKE SURE the leader not change
	s.raft.VerifyLeader()

	s.raft.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(s.RaftAddr),
				Address: raft.ServerAddress(s.RaftAddr),
			},
		},
	})
}

// Get returns the value for the given key.
func (s *Store) Get(key string) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.data[key], nil
}

func (s *Store) All() []*KV {
	s.mux.Lock()
	defer s.mux.Unlock()
	var res []*KV
	for k, v := range s.data {
		res = append(res, &KV{k, v})
	}
	return res
}

// Set sets the value for the given key.
func (s *Store) Set(key string, value []byte) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &Message{
		Op:    Set,
		Key:   key,
		Value: value,
	}

	f := s.raft.Apply(Build(c), raftTimeout)
	return f.Error()
}

// Delete deletes the given key.
func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &Message{
		Op:  Delete,
		Key: key,
	}

	f := s.raft.Apply(Build(c), raftTimeout)
	return f.Error()
}

// Clear clears all KV pairs in the store.
func (s *Store) Clear() error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	c := &Message{
		Op: Clear,
	}

	f := s.raft.Apply(Build(c), raftTimeout)
	return f.Error()
}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
// func (s *Store) Join(nodeID, addr string) error {
// 	s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)
//
// 	configFuture := s.raft.GetConfiguration()
// 	if err := configFuture.Error(); err != nil {
// 		s.logger.Printf("failed to get raft configuration: %v", err)
// 		return err
// 	}
//
// 	for _, srv := range configFuture.Configuration().Servers {
// 		// If a node already exists with either the joining node's ID or address,
// 		// that node may need to be removed from the config first.
// 		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
// 			// However if *both* the ID and the address are the same, then nothing -- not even
// 			// a join operation -- is needed.
// 			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
// 				s.logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
// 				return nil
// 			}
//
// 			future := s.raft.RemoveServer(srv.ID, 0, 0)
// 			if err := future.Error(); err != nil {
// 				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
// 			}
// 		}
// 	}
//
// 	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
// 	if f.Error() != nil {
// 		return f.Error()
// 	}
// 	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)
// 	return nil
// }

func (s *Store) Join(addr string) error {

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		console.Info("failed to get raft configuration:", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.Address == raft.ServerAddress(addr) {
			console.Info("node at", addr, "already member of cluster, ignoring join request")
			return nil
		}
	}

	f := s.raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}

	console.Info("node at", addr, "joined successfully")

	return nil
}

func (s *Store) Leave(addr string) error {

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		console.Info("failed to get raft configuration:", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.Address == raft.ServerAddress(addr) {
			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node at %s: %s", addr, err)
			}

			console.Info("node at", addr, "removed successfully")
			return nil
		}
	}

	console.Info("node at", addr, "not found")

	return nil
}

type fsm Store

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c, err = Parse(l.Data)
	if err != nil {
		panic(err)
	}

	defer f.onKeyChange(c)

	switch c.Op {
	case Set:
		return f.applySet(c.Key, c.Value)
	case Delete:
		return f.applyDelete(c.Key)
	case Clear:
		return f.applyClear()
	default:
		panic(fmt.Sprintf("unrecognized command op: %d", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	// Clone the map.
	o := make(map[string][]byte)
	for k, v := range f.data {
		o[k] = v
	}
	return &fsmSnapshot{store: o}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string][]byte)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	f.data = o
	return nil
}

func (f *fsm) applySet(key string, value []byte) interface{} {
	f.mux.Lock()
	defer f.mux.Unlock()
	f.data[key] = value
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mux.Lock()
	defer f.mux.Unlock()
	delete(f.data, key)
	return nil
}

func (f *fsm) applyClear() interface{} {
	f.mux.Lock()
	defer f.mux.Unlock()
	f.data = make(map[string][]byte)
	return nil
}

type fsmSnapshot struct {
	store map[string][]byte
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		_ = sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
