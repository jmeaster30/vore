package libvore

import (
	"io"
	"os"
	"strings"
)

type VReader struct {
	contents io.ReadSeeker
	offset   int
	size     int
}

func VReaderFromFileToMemory(filename string) *VReader {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return &VReader{
		contents: strings.NewReader(string(contents)),
		offset:   0,
		size:     len(contents),
	}
}

func VReaderFromFile(filename string) *VReader {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	return &VReader{
		contents: file,
		offset:   0,
		size:     int(fileinfo.Size()),
	}
}

func VReaderFromString(contents string) *VReader {
	return &VReader{
		contents: strings.NewReader(contents),
		offset:   0,
		size:     len(contents),
	}
}

func (v *VReader) Seek(offset int) {
	v.offset = offset
	_, err := v.contents.Seek(int64(offset), 0)
	if err != nil {
		panic(err)
	}
}

func (v *VReader) Read(length int) string {
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

func (v *VReader) ReadAt(length int, offset int) string {
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
