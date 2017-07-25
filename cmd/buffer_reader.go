package main

import (
	"io"
)

type BufferReader struct {
	line    *Line
	last    *Line
	buf     *Buffer
	lineIdx int
	offset  int
}

func NewBufferReader(buf *Buffer) *BufferReader {
	return &BufferReader{
		line: buf.lines[0],
		last: buf.lines[len(buf.lines)-1],
		buf:  buf,
	}
}

func (b *BufferReader) NextLine() (*Line, error) {
	b.lineIdx++
	if b.lineIdx < len(b.buf.lines) {
		return b.buf.lines[b.lineIdx], nil
	}
	return nil, io.EOF
}

func (b *BufferReader) Read(data []byte) (n int, err error) {
	read := 0
	for len(data) > 0 {
		left := len(b.line.data) - b.offset

		if len(data) <= left {
			// line is bigger than buffer, copy as much as we can of the line
			n := copy(data, b.line.data[b.offset:])
			read += n
			b.offset += n
			break
		}

		n := copy(data, b.line.data[b.offset:])
		read += n
		data = data[n:]

		if len(data) > 0 && b.lineIdx != len(b.buf.lines)-1 {
			data[0] = '\n'
			data = data[1:]
			read++
		}

		b.line, err = b.NextLine()
		if err != nil {
			return read, err
		}
		b.offset = 0
	}

	return read, nil
}
