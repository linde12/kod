package commands

import (
	"github.com/linde12/kod/editor"
)

type InsertRune struct {
	Rune rune
}

func (r InsertRune) Apply(e *editor.Editor) {
	v := e.CurView()
	v.Insert([]byte(string(r.Rune)), v.CursorPos())
	v.Cursor.Right()
}

type DeleteRuneBackward struct{}

func (r DeleteRuneBackward) Apply(e *editor.Editor) {
	v := e.CurView()
	v.DeleteRuneBackward()
}
