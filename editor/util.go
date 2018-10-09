package editor

import (
	"unicode/utf8"

	runewidth "github.com/mattn/go-runewidth"
)

func ByteWidth(str string, tabSize int) int {
	w := runewidth.StringWidth(str)
	lineIdx := 0
	for _, ch := range str {
		switch ch {
		case '\t':
			ts := tabSize - (lineIdx % tabSize)
			w += ts
			lineIdx += ts
		case '\n':
			lineIdx = 0
		default:
			lineIdx++
		}
	}
	return w
}

func GetCursorVisualX(x int, line string) int {
	if x > len(line) {
		x = len(line) - 1
	}

	return ByteWidth(string(line[:x]), tabSize)
}

func Count(b []byte) int {
	return utf8.RuneCountInString(string(b))
}
