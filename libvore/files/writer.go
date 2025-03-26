package files

import (
	"io"
	"os"
)

type Writer struct {
	contents WriteSeekCloser
}

func WriterFromFile(filename string) *Writer {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		panic(err)
	}

	return &Writer{
		contents: file,
	}
}

func WriterFromMemory() *Writer {
	return &Writer{
		contents: NewMemoryStream(),
	}
}

func (vw *Writer) WriteAt(offset int, data string) {
	_, serr := vw.contents.Seek(int64(offset), io.SeekStart)
	if serr != nil {
		panic(serr)
	}
	_, werr := vw.contents.Write([]byte(data))
	if werr != nil {
		panic(werr)
	}
}

func (vw *Writer) Close() {
	err := vw.contents.Close()
	if err != nil {
		panic(err)
	}
}
