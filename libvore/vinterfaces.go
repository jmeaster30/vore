package libvore

import "strings"

type WriteSeekCloser interface {
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

type ReadSeekCloser interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

type StringReadSeekCloser struct {
	contents *strings.Reader
}

func NewStringReadCloser(value string) *StringReadSeekCloser {
	return &StringReadSeekCloser{
		contents: strings.NewReader(value),
	}
}

func (srsc *StringReadSeekCloser) Read(p []byte) (int, error) {
	return srsc.contents.Read(p)
}

func (srsc *StringReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return srsc.contents.Seek(offset, whence)
}

func (srsc *StringReadSeekCloser) Close() error {
	return nil // don't need to do anything here :)
}
