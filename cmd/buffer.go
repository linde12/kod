package main

import (
	"io"
	"os"
)

type Buffer struct {
	*LineArray
	Cursor Cursor
	Path   string
}

func NewBuffer(in io.Reader, path string) *Buffer {
	la := NewLineArray(in)
	b := &Buffer{}
	b.LineArray = la
	b.Cursor = Cursor{buf: b}
	b.Path = path

	return b
}

func (b *Buffer) Save() error {
	return b.SaveAs(b.Path)
}

func (b *Buffer) SaveAs(filename string) error {
	r := NewBufferReader(b)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}

func (b *Buffer) CurLine() *Line {
	return b.GetLine(b.Cursor.Y)
}

func (b *Buffer) GetLine(y int) *Line {
	return b.lines[y]
}

func (b *Buffer) CursorPos() Pos {
	return Pos{b.Cursor.X, b.Cursor.Y}
}

func (b *Buffer) CursorUp() {
	if b.Cursor.Y > 0 {
		b.Cursor.Y -= 1

		if b.Cursor.X > Count(b.CurLine().data) {
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
	b.Cursor.StoreVisualX()
}

func (b *Buffer) CursorRight() {
	if b.Cursor.X < Count(b.CurLine().data) {
		b.Cursor.X += 1
	} else if b.Cursor.Y+1 < len(b.lines) {
		// Move down one line
		b.CursorDown()
		b.CursorBegin()
	}
	b.Cursor.StoreVisualX()
}

func (b *Buffer) CursorEnd() {
	line := b.CurLine()
	n := Count(line.data)
	if n > 0 {
		b.Cursor.X = n
	} else {
		b.Cursor.X = 0
	}
}

func (b *Buffer) CursorBegin() {
	b.Cursor.X = 0
}
