package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/shashankshrivastva-hue/distributed-cache-engine/pkg/protocol"
	"github.com/shashankshrivastva-hue/distributed-cache-engine/pkg/storage"
)

func main() {
	cache := storage.NewLRUCache(100000)
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to bind TCP port 6379: %v", err)
	}
	defer listener.Close()

	fmt.Println("⚡ Distributed Cache Engine listening on port :6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, cache)
	}
}

func handleClient(conn net.Conn, cache *storage.LRUCache) {
	defer conn.Close()
	parser := protocol.NewParser(conn)

	for {
		cmd, err := parser.ParseCommand()
		if err != nil {
			return
		}
		if len(cmd) == 0 {
			continue
		}

		action := strings.ToUpper(cmd[0])
		switch action {
		case "PING":
			conn.Write([]byte(protocol.FormatSimpleString("PONG")))
		case "SET":
			if len(cmd) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for SET\r\n"))
				continue
			}
			cache.Set(cmd[1], []byte(cmd[2]), 0)
			conn.Write([]byte(protocol.FormatSimpleString("OK")))
		case "GET":
			if len(cmd) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for GET\r\n"))
				continue
			}
			val, ok := cache.Get(cmd[1])
			if !ok {
				conn.Write([]byte(protocol.FormatNullBulk()))
			} else {
				conn.Write([]byte(protocol.FormatBulkString(string(val))))
			}
		case "DEL":
			if len(cmd) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for DEL\r\n"))
				continue
			}
			cache.Delete(cmd[1])
			conn.Write([]byte(":1\r\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", action)))
		}
	}
}
