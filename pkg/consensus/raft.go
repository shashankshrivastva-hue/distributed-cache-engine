package consensus

import (
	"sync"
	"time"
)

type Role string

const (
	Follower  Role = "FOLLOWER"
	Candidate Role = "CANDIDATE"
	Leader    Role = "LEADER"
)

type RaftNode struct {
	mu          sync.Mutex
	NodeID      string
	CurrentTerm int
	Role        Role
	Peers       []string
	Log         []LogEntry
	CommitIndex int
}

type LogEntry struct {
	Term    int
	Command string
	Key     string
	Value   []byte
}

func NewRaftNode(nodeID string, peers []string) *RaftNode {
	return &RaftNode{
		NodeID:      nodeID,
		CurrentTerm: 0,
		Role:        Follower,
		Peers:       peers,
		Log:         make([]LogEntry, 0),
		CommitIndex: 0,
	}
}

func (rn *RaftNode) StartElection() {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	rn.Role = Candidate
	rn.CurrentTerm++
}

func (rn *RaftNode) AppendEntry(key string, val []byte) LogEntry {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	entry := LogEntry{
		Term:    rn.CurrentTerm,
		Command: "SET",
		Key:     key,
		Value:   val,
	}
	rn.Log = append(rn.Log, entry)
	return entry
}
