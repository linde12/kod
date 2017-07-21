package main

import "log"

type Cursor struct {
	Pos
	buf         *Buffer
	LastVisualX int
}

// TODO: retrieve from settings
const tabSize = 4

func (c *Cursor) GetVisualX() int {
	r := []rune(string(c.buf.GetLine(c.Y).data))
	if c.X > len(r) {
		c.X = len(r) - 1
	}

	return ByteWidth(string(r[:c.X]), tabSize)
}

func (c *Cursor) StoreVisualX() {
	log.Printf("VisualX=%v\n", c.GetVisualX())
	c.LastVisualX = c.GetVisualX()
}
