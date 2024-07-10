// Package flakes wrapper around snowflake
package flakes

import (
	"github.com/bwmarrin/snowflake"
)

func init() {
	// based on discord epoch because it's handy
	// for the existing convert tools
	const epoch = 1420070400000
	snowflake.Epoch = epoch
}

// Generator wrapper around snowflake
type Generator struct {
	node *snowflake.Node
}

// New sets up a new flakes Generator
func New(n int64) (Generator, error) {
	node, err := snowflake.NewNode(n)
	if err != nil {
		return Generator{}, err
	}

	return Generator{
		node: node,
	}, nil
}

// ID returns a snowflake as int64
func (g Generator) ID() int64 {
	return g.node.Generate().Int64()
}
