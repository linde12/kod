package main

type Buffer struct {
	*LineArray
	Cursor Cursor
}

func NewBuffer() *Buffer {
	b := Buffer{
		LineArray: &LineArray{},
	}
	b.AppendLine()

	return &b
}

func (b *Buffer) CursorPos() (int, int) {
	return b.Cursor.x, b.Cursor.y
}

func (b *Buffer) CurLine() *Line {
	_, y := b.CursorPos()
	return b.lines[y]
}

func (b *Buffer) CursorUp() {
	if b.Cursor.y > 0 {
		b.Cursor.SetY(b.Cursor.y - 1)

		if b.Cursor.x > len(b.CurLine().data) {
			b.CursorEnd()
		}
	}
}

func (b *Buffer) CursorDown() {
	if b.Cursor.y+1 < len(b.lines) {
		b.Cursor.SetY(b.Cursor.y + 1)

		if b.Cursor.x > len(b.CurLine().data) {
			b.CursorEnd()
		}
	}
}

func (b *Buffer) CursorLeft() {
	if b.Cursor.x > 0 {
		b.Cursor.SetX(b.Cursor.x - 1)
	} else if b.Cursor.y > 0 {
		// Move up one line
		b.CursorUp()
		b.CursorEnd()
	}
}

func (b *Buffer) CursorRight() {
	if b.Cursor.x < len(b.CurLine().data) {
		b.Cursor.SetX(b.Cursor.x + 1)
	} else if b.Cursor.y+1 < len(b.lines) {
		// Move down one line
		b.CursorDown()
		b.CursorBegin()
	}
}

func (b *Buffer) CursorEnd() {
	line := b.CurLine()
	b.Cursor.SetX(len(line.data))
}

func (b *Buffer) CursorBegin() {
	b.Cursor.SetX(0)
}
