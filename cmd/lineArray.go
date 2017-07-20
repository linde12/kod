package main

import (
	"strings"
	"unicode/utf8"
)

type Line struct {
	data []byte
}

type LineArray struct {
	lines []*Line
}

// converts a position of a rune(on the screen) to it's position in the byte array
// graciously stolen from micro
func runeToByteIndex(txt []byte, n int) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
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
	dst.Insert(src.data[x:], 0)
	src.DeleteToEnd(x)
}

func (la *LineArray) JoinLines(a int, b int) {
	line := la.lines[a]
	line.Insert(la.lines[b].data, len(line.data))
	la.RemoveLine(b)
}

func (l *Line) Insert(value []byte, pos int) {
	x := runeToByteIndex(l.data, pos)
	for i := 0; i < len(value); i++ {
		l.insertByte(value[i], x)
		x++
	}
}

func (l *Line) insertByte(value byte, pos int) {
	l.data = append(l.data, 0)
	// shift everything to the right
	copy(l.data[pos+1:], l.data[pos:])
	l.data[pos] = value
}

func (l *Line) DeleteToEnd(pos int) {
	l.data = l.data[:pos]
}

func (l *Line) RemoveRune(pos int) {
	startX := runeToByteIndex(l.data, pos)
	endX := runeToByteIndex(l.data, pos+1)
	l.data = append(l.data[:startX], l.data[endX:]...)
}
