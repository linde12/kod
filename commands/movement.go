package commands

import (
	"github.com/linde12/kod/editor"
)

const (
	MoveUp        = 0
	MoveDown      = 1
	MoveLeft      = 2
	MoveRight     = 3
	MoveWordRight = 4
)

type MoveEOL struct{}

func (m MoveEOL) Apply(e *editor.Editor) {
	v := e.CurView()
	c := v.Cursor
	c.End()
}

type MoveRune struct {
	Dir int
}

func (m MoveRune) Apply(e *editor.Editor) {
	v := e.CurView()
	switch m.Dir {
	case MoveUp:
		v.Cursor.Up()
	case MoveDown:
		v.Cursor.Down()
	case MoveLeft:
		v.Cursor.Left()
	case MoveRight:
		v.Cursor.Right()
	}
}
