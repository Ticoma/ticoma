package main

import (
	"fmt"
	"os"

	nodes "ticoma/packages/nodes"
)

func main() {
	fmt.Printf("Hello from main")
	nodes.GameNode()
	nodes.PlayerNode()
	os.Exit(0)
}
