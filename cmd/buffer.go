package main

import (
	"log"

	"github.com/gdamore/tcell"
)

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

func (b *Buffer) HandleEvent(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		log.Println("Keypress:", e.Name())
		if e.Key() == tcell.KeyRune {
			line := buf.CurLine()
			line.InsertRune(e.Rune(), buf.Cursor.x)
			buf.CursorRight()
		} else {
			switch e.Key() {
			case tcell.KeyBackspace2, tcell.KeyBackspace:
				line := buf.CurLine()
				if len(line.data) > 0 {
					if b.Cursor.x > 0 {
						line.RemoveRune(buf.Cursor.x - 1)
						buf.CursorLeft()
					} else if buf.Cursor.y != 0 {
						lineAbove := buf.lines[buf.Cursor.y-1]
						end := len(lineAbove.data)
						buf.JoinLines(buf.Cursor.y-1, buf.Cursor.y)
						buf.CursorUp()
						buf.Cursor.SetX(end)
					}
				} else if buf.Cursor.y != 0 {
					buf.RemoveLine(buf.Cursor.y)
					buf.CursorUp()
					buf.CursorEnd()
				}
			case tcell.KeyEnter:
				buf.Split(buf.Cursor.x, buf.Cursor.y)
				buf.CursorDown()
				buf.CursorBegin()
			case tcell.KeyLeft:
				buf.CursorLeft()
			case tcell.KeyUp:
				buf.CursorUp()
			case tcell.KeyRight:
				buf.CursorRight()
			case tcell.KeyDown:
				buf.CursorDown()
			}
		}
	}
}
