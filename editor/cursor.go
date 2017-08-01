package editor

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
	c.LastVisualX = c.GetVisualX()
}

func (c *Cursor) Up() {
	if c.Y > 0 {
		c.Y -= 1

		if c.X > Count(c.buf.CurLine().data) {
			c.End()
		}
	}
}

func (c *Cursor) Down() {
	if c.Y+1 < len(c.buf.lines) {
		c.Y += 1

		if c.X > len(c.buf.CurLine().data) {
			c.End()
		}
	}
}

func (c *Cursor) Left() {
	if c.X > 0 {
		c.X -= 1
	} else if c.Y > 0 {
		// Move up one line
		c.Up()
		c.End()
	}
	c.StoreVisualX()
}

func (c *Cursor) Right() {
	if c.X < Count(c.buf.CurLine().data) {
		c.X += 1
	} else if c.Y+1 < len(c.buf.lines) {
		// Move down one line
		c.Down()
		c.Begin()
	}
	c.StoreVisualX()
}

func (c *Cursor) End() {
	line := c.buf.CurLine()
	n := Count(line.data)
	if n > 0 {
		c.X = n
	} else {
		c.X = 0
	}
}

func (c *Cursor) Begin() {
	c.X = 0
}

func (c *Cursor) StartOfBuffer() {
	c.Y = 0
	c.X = 0
}

func (c *Cursor) EndOfBuffer() {
	c.Y = len(c.buf.lines) - 1
	c.End()
}
