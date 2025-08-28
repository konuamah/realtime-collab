package main

import (
	"sync"
)

type Operation struct {
	Type  string // "insert" or "delete"
	Char  string // character to insert
	Index int    // position in text
}

type CRDT struct {
	Text string
	mu   sync.Mutex
}

func (c *CRDT) ApplyOp(op Operation) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch op.Type {
	case "insert":
		if op.Index >= len(c.Text) {
			c.Text += op.Char
		} else {
			c.Text = c.Text[:op.Index] + op.Char + c.Text[op.Index:]
		}
	case "delete":
		if op.Index < len(c.Text) && len(c.Text) > 0 {
			c.Text = c.Text[:op.Index] + c.Text[op.Index+1:]
		}
	}
}

func (c *CRDT) GetText() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Text
}
