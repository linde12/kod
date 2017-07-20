package main

import "io"

type Buffer struct {
	*LineArray
	Cursor Cursor
}

func NewBuffer(in io.Reader) *Buffer {
	la := NewLineArray(in)
	b := Buffer{
		LineArray: la,
	}

	return &b
}

func (b *Buffer) CurLine() *Line {
	return b.lines[b.Cursor.Y]
}

func (b *Buffer) CursorPos() Pos {
	return Pos{b.Cursor.X, b.Cursor.Y}
}

func (b *Buffer) CursorUp() {
	if b.Cursor.Y > 0 {
		b.Cursor.Y -= 1

		if b.Cursor.X > len(b.CurLine().data) {
			b.CursorEnd()
		}
	}
}

func (b *Buffer) CursorDown() {
	if b.Cursor.Y+1 < len(b.lines) {
		b.Cursor.Y += 1

		if b.Cursor.X > len(b.CurLine().data) {
			b.CursorEnd()
		}
	}
}

func (b *Buffer) CursorLeft() {
	if b.Cursor.X > 0 {
		b.Cursor.X -= 1
	} else if b.Cursor.Y > 0 {
		// Move up one line
		b.CursorUp()
		b.CursorEnd()
	}
}

func (b *Buffer) CursorRight() {
	if b.Cursor.X < len(b.CurLine().data) {
		b.Cursor.X += 1
	} else if b.Cursor.Y+1 < len(b.lines) {
		// Move down one line
		b.CursorDown()
		b.CursorBegin()
	}
}

func (b *Buffer) CursorEnd() {
	line := b.CurLine()
	b.Cursor.X = len(line.data)
}

func (b *Buffer) CursorBegin() {
	b.Cursor.X = 0
}
