package main

type Cursor struct {
	x int
	y int
}

func (c *Cursor) SetX(x int) {
	c.x = x
}

func (c *Cursor) SetY(y int) {
	c.y = y
}
