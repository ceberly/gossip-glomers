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

func main() {
	var neighbors []interface{}

	var seen Seen
	seen.seenSet = make(map[float64]bool)
	seen.seen = make([]float64, 0, 10000)

	n := maelstrom.NewNode()

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

			var wg sync.WaitGroup
			for _, neighbor := range neighbors {
				wg.Add(1)
				go func(node string) {
					defer wg.Done()
					n.Send(node, body)
				}(neighbor.(string))
			}

			wg.Wait()
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
