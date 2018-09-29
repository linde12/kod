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

func (lc *LineCache) ApplyUpdate(update *rpc.Update) {
	// TODO: Make capacity = height of buffer
	newLines := make([]*Line, 0, 10)
	newInvalidBefore := 0
	newInvalidAfter := 0

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
				op.N -= len(lc.lines)
				lc.lines = make([]*Line, 0, 10)
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
		}
	}

	lc.lines = newLines
	lc.invalidBefore = newInvalidBefore
	lc.invalidAfter = newInvalidAfter
}
