package files

import (
	"errors"
	"io"
)

type MemoryStream struct {
	contents []byte
	pos      int // TODO need to update this to int64
}

func NewMemoryStream() *MemoryStream {
	return &MemoryStream{
		contents: []byte{},
		pos:      0,
	}
}

func (ms *MemoryStream) Write(buf []byte) (n int, err error) {
	minCap := ms.pos + len(buf)
	if minCap > cap(ms.contents) { // Make sure buf has enough capacity:
		buf2 := make([]byte, len(ms.contents), 2*(ms.pos+len(buf))) // add some extra
		copy(buf2, ms.contents)
		ms.contents = buf2
	}
	if minCap > len(ms.contents) {
		ms.contents = ms.contents[:minCap]
	}
	copy(ms.contents[ms.pos:], buf)
	ms.pos += len(buf)
	return len(buf), nil
}

func (ms *MemoryStream) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = ms.pos + offs
	case io.SeekEnd:
		newPos = len(ms.contents) + offs
	}
	if newPos < 0 {
		return int64(ms.pos), errors.New("negative result pos")
	}
	ms.pos = newPos
	return int64(newPos), nil
}

func (ms *MemoryStream) Close() error {
	return nil // Don't need to do anything here :)
}
