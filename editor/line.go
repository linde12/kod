package editor

import "github.com/gdamore/tcell"

type Line struct {
	Text    string
	Cursors []int
	styles  map[int]tcell.Color
}

func NewLine(text string, cursors []int, styles []int) *Line {
	line := &Line{}
	line.Text = text
	line.Cursors = cursors
	return line
}

// TODO: Implement syntax highlight
//func (l *Line) SetStyles(styles []int) {
//for i := 0; i < len(styles); i+=3 {
//start := offset + styles[i]
//end := start + styles[i+1]
//styleId := styles[i+2]

//style := &tcell.Style{}
//}
//}
