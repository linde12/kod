package editor

import (
	"io"
	"os"
)

type Buffer struct {
	*LineArray
	Cursor *Cursor
	Path   string
}

func NewBuffer(in io.Reader, path string) *Buffer {
	la := NewLineArray(in)
	b := &Buffer{}
	b.LineArray = la
	b.Cursor = &Cursor{buf: b}
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

func (b *Buffer) DeleteRuneBackward() {
	line := b.CurLine()
	if len(line.data) > 0 {
		if b.Cursor.X > 0 {
			b.RemoveRune(Pos{b.Cursor.X - 1, b.Cursor.Y})
			b.Cursor.Left()
		} else if b.Cursor.Y != 0 {
			lineAbove := b.lines[b.Cursor.Y-1]
			end := len(lineAbove.data)
			b.JoinLines(b.Cursor.Y-1, b.Cursor.Y)
			b.Cursor.Up()
			b.Cursor.X = end
		}
	}
}
