# ⚡ Distributed Cache Engine

[![Go Report Card](https://goreportcard.com/badge/github.com/shashankshrivastva-hue/distributed-cache-engine)](https://goreportcard.com/report/github.com/shashankshrivastva-hue/distributed-cache-engine)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)](https://golang.org)

High-performance, in-memory distributed key-value cache built in Go with Raft consensus, 32-way sharded locking, and native Redis RESP protocol compliance.

---

## 🏛️ Architecture

```mermaid
graph TD
    Client[Redis CLI / TCP Client] -->|RESP Commands| Network[TCP Server :6379]
    Network --> Parser[RESP Protocol Parser]
    Parser --> Router{Command Router}
    Router -->|Read/Write| Shard[32-Way Sharded Store]
    Router -->|LRU Eviction| LRU[LRU Cache & TTL Ticker]
    Router -->|Replication| Raft[Raft Consensus Node]
    Raft -->|Heartbeat| Peers[Cluster Raft Peers]
```

## 🚀 Key Features

- **RESP Protocol Compatibility**: Connect seamlessly via standard `redis-cli` or any language Redis driver.
- **Microsecond Latency**: 32-way memory sharding reduces thread contention across high core counts.
- **LRU Eviction & Precise TTL**: O(1) eviction using doubly-linked list with auto-expiring keys.
- **Raft Consensus Engine**: Built-in state machine for replicated log consistency.

## 🛠️ Usage

```bash
# Build and run
go run cmd/cache-server/main.go

# Connect via standard redis-cli
redis-cli -p 6379
127.0.0.1:6379> SET user:101 "Shashank Shrivastva"
OK
127.0.0.1:6379> GET user:101
"Shashank Shrivastva"
```

## 📜 License
MIT License. Built by [Shashank Shrivastva](https://github.com/shashankshrivastva-hue).
