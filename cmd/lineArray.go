package main

import "strings"

type Line struct {
	data []byte
}

type LineArray struct {
	lines []*Line
}

func NewLineArray(str string) *LineArray {
	la := LineArray{}

	// alloc 1000 lines by default
	la.lines = make([]*Line, 0, 1000)

	lines := strings.Split(str, "\n")

	for _, txt := range lines {
		//TODO: Allocate more if over 1000 lines
		line := &Line{
			data: []byte(txt[:len(txt)]),
		}
		la.lines = append(la.lines, line)
	}

	return &la
}

func (la *LineArray) AppendLine() {
	la.lines = append(la.lines, &Line{data: []byte{}})
}

func (la *LineArray) InsertLineAfter(pos int) {
	// insert after current line
	pos++
	la.AppendLine()
	// shift everything to the right
	copy(la.lines[pos+1:], la.lines[pos:])
	la.lines[pos] = &Line{}
}

func (la *LineArray) RemoveLine(pos int) {
	la.lines = append(la.lines[:pos], la.lines[pos+1:]...)
}

func (la *LineArray) Split(x int, y int) {
	la.InsertLineAfter(y)
	src := la.lines[y]
	dst := la.lines[y+1]

	dstData := src.data[x:]
	dst.data = make([]byte, len(dstData))
	copy(dst.data, dstData)
	src.data = src.data[:x]
}

func (la *LineArray) JoinLines(a int, b int) {
	line := la.lines[a]
	line.data = append(line.data, la.lines[b].data...)
	la.RemoveLine(b)
}

func (l *Line) InsertRune(r rune, pos int) {
	l.data = append(l.data, 0)
	// shift everything to the right
	copy(l.data[pos+1:], l.data[pos:])
	l.data[pos] = byte(r)
}

func (l *Line) RemoveRune(pos int) {
	l.data = append(l.data[:pos], l.data[pos+1:]...)
}
