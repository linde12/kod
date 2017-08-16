package editor

import (
	"math"

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

func (lc *LineCache) ApplyUpdate(msg *rpc.Message) {
	switch msg.Value.(type) {
	case *rpc.Update:
		update := msg.Value.(*rpc.Update)
		// TODO: Make capacity = height of buffer
		newLines := make([]*Line, 0, 10)
		newInvalidBefore := 0
		newInvalidAfter := 0
		index := 0

		for _, op := range update.Update.Ops {
			switch op.Op {
			case "copy":
				// TODO: Respect lc.invalidBefore
				if index < lc.invalidBefore {
					invalid := int(math.Min(float64(op.N), float64(lc.invalidBefore-index)))
					lc.addInvalid(newLines, &newInvalidBefore, &newInvalidAfter, invalid)
					op.N -= invalid
					index += invalid
				}
				for op.N > 0 && index < lc.invalidBefore+len(lc.lines) {
					newLines = lc.addLine(newLines, &newInvalidBefore, &newInvalidAfter, lc.lines[index-lc.invalidBefore])
					op.N--
					index++
				}
				lc.addInvalid(newLines, &newInvalidBefore, &newInvalidAfter, op.N)
				index += op.N
			case "skip":
				index += op.N
			case "invalidate":
				lc.addInvalid(newLines, &newInvalidBefore, &newInvalidAfter, op.N)
			case "ins":
				for _, line := range op.Lines {
					newline := NewLine(line.Text, line.Cursor, line.Styles)
					newLines = lc.addLine(newLines, &newInvalidBefore, &newInvalidAfter, newline)
				}
			case "update":
				for _, line := range op.Lines {
					lineToUpdate := lc.lines[index-lc.invalidBefore]
					lineToUpdate.Cursors = line.Cursor
					newLines = lc.addLine(newLines, &newInvalidBefore, &newInvalidAfter, lineToUpdate)
					index++
				}
			}
		}
		lc.lines = newLines
		lc.invalidBefore = newInvalidBefore
		lc.invalidAfter = newInvalidAfter
	}
}
