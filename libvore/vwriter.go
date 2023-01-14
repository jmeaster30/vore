package libvore

import (
	"io"
	"os"
)

type VWriter struct {
	contents io.WriteSeeker
	offset   int
}

func NewVWriter(filename string) *VWriter {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		panic(err)
	}

	return &VWriter{
		contents: file,
		offset:   0,
	}
}

func DummyVWriter() *VWriter {
	return &VWriter{
		contents: NewMemoryStream(),
		offset:   0,
	}
}

func (vw *VWriter) WriteAt(offset int, data string) {
	_, serr := vw.contents.Seek(int64(offset), io.SeekStart)
	if serr != nil {
		panic(serr)
	}
	_, werr := vw.contents.Write([]byte(data))
	if werr != nil {
		panic(werr)
	}
}
