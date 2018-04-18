package editor

import (
	"github.com/linde12/kod/rpc"
)

type LineCache struct {
	lines         []*Line
	invalidBefore int
	invalidAfter  int
}

func NewLineCache() *LineCache {
	return &LineCache{
		// TODO: Make capacity = height of buffer
		lines: make([]*Line, 0, 10),
	}
}

func (lc *LineCache) addInvalid(newLines []*Line, newInvalidBefore *int, newInvalidAfter *int, n int) {
	if len(newLines) == 0 {
		*newInvalidBefore += n
	} else {
		*newInvalidAfter += n
	}
}

func (lc *LineCache) addLine(newLines []*Line, newInvalidBefore *int, newInvalidAfter *int, line *Line) []*Line {
	if line != nil {
		for i := 0; i < *newInvalidAfter; i++ {
			newLines = append(newLines, nil)
		}
		*newInvalidAfter = 0
		newLines = append(newLines, line)
	} else {
		lc.addInvalid(newLines, newInvalidBefore, newInvalidAfter, 1)
	}
	return newLines
}

func (lc *LineCache) ApplyUpdate(update *rpc.Update) {
	// TODO: Make capacity = height of buffer
	newLines := make([]*Line, 0, 10)
	newInvalidBefore := 0
	newInvalidAfter := 0
	index := 0

	for _, op := range update.Update.Ops {
		switch op.Op {
		case "copy":
			if lc.invalidBefore > op.N {
				lc.invalidBefore -= op.N
				newInvalidBefore += op.N
				continue
			} else if lc.invalidAfter > 0 {
				op.N -= lc.invalidBefore
				newInvalidBefore += lc.invalidBefore
				lc.invalidBefore = 0
			}

			if op.N < len(lc.lines) {
				firstNLines := append(lc.lines[:0], lc.lines[:op.N]...)
				lc.lines = lc.lines[op.N:]
				newLines = append(newLines, firstNLines...)
				continue
			} else {
				newLines = append(newLines, lc.lines...)
				lc.lines = make([]*Line, 0, 10)
				op.N -= len(lc.lines)
			}

			if lc.invalidAfter >= op.N {
				lc.invalidAfter -= op.N
				newInvalidAfter += op.N
				continue
			}

		case "skip":
			if lc.invalidBefore > op.N {
				lc.invalidBefore -= op.N
				continue
			} else if lc.invalidBefore > 0 {
				op.N = lc.invalidBefore
				lc.invalidBefore = 0
			}

			if op.N < len(lc.lines) {
				lc.lines = append(lc.lines[:0], lc.lines[op.N:]...)
				continue
			} else {
				lc.lines = append(lc.lines[:0], lc.lines[len(lc.lines)-1:]...)
				op.N -= len(lc.lines)
			}

			if lc.invalidAfter >= op.N {
				lc.invalidBefore -= op.N
				continue
			}

		case "invalidate":
			lc.addInvalid(newLines, &newInvalidBefore, &newInvalidAfter, op.N)
		case "ins":
			for _, line := range op.Lines {
				newline := NewLine(line.Text, line.Cursor, line.Styles)
				newLines = append(newLines, newline)
			}
		case "update":
			for _, line := range op.Lines {
				lineToUpdate := lc.lines[index-lc.invalidBefore]
				lineToUpdate.Cursors = line.Cursor
				lineToUpdate.SetStyles(line.Styles)
				newLines = lc.addLine(newLines, &newInvalidBefore, &newInvalidAfter, lineToUpdate)
				index++
			}
		}
	}

	lc.lines = newLines
	lc.invalidBefore = newInvalidBefore
	lc.invalidAfter = newInvalidAfter
}

