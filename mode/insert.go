package mode

import (
	"github.com/gdamore/tcell"
	"github.com/linde12/kod/commands"
	"github.com/linde12/kod/editor"
)

type InsertMode struct {
	editor *editor.Editor
	count  string
}

func NewInsertMode(e *editor.Editor) *InsertMode {
	return &InsertMode{e, ""}
}

func (m *InsertMode) OnKey(ev *tcell.EventKey) {
	editor := m.editor

	if ev.Key() == tcell.KeyRune {
		editor.Commands <- commands.InsertRune{ev.Rune()}
	} else {
		switch ev.Key() {
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			editor.Commands <- commands.DeleteRuneBackward{}
		case tcell.KeyTAB:
			editor.Commands <- commands.InsertRune{'\t'}
		case tcell.KeyEnter:
			editor.Commands <- commands.InsertRune{'\n'}
		case tcell.KeyEscape:
			editor.SetMode(NewNormalMode(editor))
		}
	}
}
