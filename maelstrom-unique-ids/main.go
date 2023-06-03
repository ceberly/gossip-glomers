package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	var counter int64 = 0

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resp := make(map[string]any)

		resp["type"] = "generate_ok"
		resp["id"] = fmt.Sprintf("%s:%d", n.ID(), counter)
		counter += 1

		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting...")
}
