package commands

import (
	"github.com/linde12/kod/editor"
)

const (
	MoveUp = iota
	MoveDown
	MoveLeft
	MoveRight
	MoveWordRight
)

type MoveEOL struct{}

func (m MoveEOL) Apply(e *editor.Editor) {
	v := e.CurView()
	c := v.Cursor
	c.End()
}

type MoveBOL struct{}

func (m MoveBOL) Apply(e *editor.Editor) {
	v := e.CurView()
	c := v.Cursor
	c.Begin()
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

type MoveStartOfBuffer struct{}

func (m MoveStartOfBuffer) Apply(e *editor.Editor) {
	v := e.CurView()
	c := v.Cursor
	c.StartOfBuffer()
}

type MoveEndOfBuffer struct{}

func (m MoveEndOfBuffer) Apply(e *editor.Editor) {
	v := e.CurView()
	c := v.Cursor
	c.EndOfBuffer()
}
