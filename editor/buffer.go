package editor

import (
	"io"
	"os"
)

// Buffer is where a content of a file belongs. It is used to represent what should be rendered to the screen.
type Buffer struct {
	*LineArray
	Cursor *Cursor
	Path   string
}

// NewBuffer creates a new buffer by reading from the passed `in` argument.
func NewBuffer(path string) *Buffer {
	la := NewLineArray()
	b := &Buffer{}
	b.LineArray = la
	b.Cursor = &Cursor{buf: b}
	b.Path = path

	return b
}

// Save saves the buffer to the path arguments passed to `NewBuffer`.
func (b *Buffer) Save() error {
	return b.SaveAs(b.Path)
}

// SaveAs saves the buffer to the passed path.
func (b *Buffer) SaveAs(path string) error {
	r := NewBufferReader(b)
	f, err := os.Create(path)
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

// CurLine returns the current Y position of the main cursor.
func (b *Buffer) CurLine() *Line {
	return b.GetLine(b.Cursor.Y)
}

// GetLine returns the current line under the main cursor.
func (b *Buffer) GetLine(y int) *Line {
	return b.lines[y]
}

// CursorPos returns the current position of the cursor in form of a `Pos`.
func (b *Buffer) CursorPos() Pos {
	return Pos{b.Cursor.X, b.Cursor.Y}
}

// DeleteRuneBackward deletes a rune backward starting at the main cursor's position.
// If there is no more rune behind the cursor's position, the line will be joined with the line above.
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
