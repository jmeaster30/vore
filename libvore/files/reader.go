package files

import (
	"io"
	"os"
)

type Reader struct {
	contents ReadSeekCloser
	offset   int
	size     int
}

func ReaderFromFileToMemory(filename string) *Reader {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return &Reader{
		contents: NewStringReadCloser(string(contents)),
		offset:   0,
		size:     len(contents),
	}
}

func ReaderFromFile(filename string) *Reader {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	return &Reader{
		contents: NewBufferedFile(file, fileinfo.Size()),
		offset:   0,
		size:     int(fileinfo.Size()),
	}
}

func ReaderFromString(contents string) *Reader {
	return &Reader{
		contents: NewStringReadCloser(contents),
		offset:   0,
		size:     len(contents),
	}
}

func (v *Reader) Size() int {
	return v.size
}

func (v *Reader) Seek(offset int) {
	v.offset = offset
	_, err := v.contents.Seek(int64(offset), io.SeekStart)
	if err != nil {
		panic(err)
	}
}

func (v *Reader) Read(length int) string {
	if v.offset+length-1 >= v.size {
		return ""
	}
	currentString := make([]byte, length)
	n, err := v.contents.Read(currentString)
	if err != nil {
		panic(err)
	}
	if n != length {
		return ""
	}
	return string(currentString)
}

func (v *Reader) ReadAt(length int, offset int) string {
	if offset+length-1 >= v.size {
		return ""
	}
	currentString := make([]byte, length)
	v.Seek(offset)
	n, err := v.contents.Read(currentString)
	if err != nil {
		panic(err)
	}
	if n != length {
		return ""
	}
	return string(currentString)
}

func (v *Reader) Close() {
	err := v.contents.Close()
	if err != nil {
		panic(err)
	}
}
