package main

type Line struct {
	data []byte
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
