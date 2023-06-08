package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// keep a set of seen values for uniqueness.
// keep the keys which are what we are interested in a separate slice
// so that the list doesn't have to be created on every read() message.
// This is a pretty inefficient way to store this, I think, but I'm
// trying to focus on the distributed systems part of these challenges :)
type Seen struct {
	sync.Mutex
	seenSet map[float64]bool
	seen    []float64
}

type GossipQueueItem struct {
	NodeId  string
	Message float64
}

func handleGossipQueue(node *maelstrom.Node, queue chan GossipQueueItem) {
	neighborMessages := make(map[string]chan float64)

	for r := range queue {
		// set up a new neighbor message queue
		q, ok := neighborMessages[r.NodeId]
		if !ok {
			neighborMessages[r.NodeId] = make(chan float64)
			defer close(neighborMessages[r.NodeId])
			go func(nodeId string) {
				for m := range neighborMessages[nodeId] {
					node.Send(nodeId, m)
				}
			}(r.NodeId)
		}

		q <- r.Message
	}
}

func main() {
	var neighbors []interface{}

	var seen Seen
	seen.seenSet = make(map[float64]bool)
	seen.seen = make([]float64, 0, 10000)

	n := maelstrom.NewNode()

	// maps node id -> slice of messages to retry
	gossipQueue := make(chan GossipQueueItem, 1)

	defer close(gossipQueue)
	go handleGossipQueue(n, gossipQueue)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		if neighbors == nil {
			log.Fatal("no neighbors set")
		}

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := body["message"].(float64)

		seen.Lock()
		if _, ok := seen.seenSet[message]; !ok {
			seen.seenSet[message] = true
			seen.seen = append(seen.seen, message)
			seen.Unlock()

			for _, neighbor := range neighbors {
				gossipQueue <- GossipQueueItem{neighbor.(string), message}
			}
		} else {
			seen.Unlock()
		}

		resp := make(map[string]any)

		resp["type"] = "broadcast_ok"

		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resp := make(map[string]any)

		resp["type"] = "read_ok"

		seen.Lock()
		resp["messages"] = seen.seen
		seen.Unlock()

		return n.Reply(msg, resp)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body["topology"].(map[string]interface{})
		neighbors = topology[n.ID()].([]interface{})

		resp := make(map[string]any)

		resp["type"] = "topology_ok"

		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting...")
}
