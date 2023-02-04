package libvore

import (
	"io"
	"os"
)

type VWriter struct {
	contents WriteSeekCloser
}

func VWriterFromFile(filename string) *VWriter {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		panic(err)
	}

	return &VWriter{
		contents: file,
	}
}

func VWriterFromMemory() *VWriter {
	return &VWriter{
		contents: NewMemoryStream(),
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

func (vw *VWriter) Close() {
	err := vw.contents.Close()
	if err != nil {
		panic(err)
	}
}
