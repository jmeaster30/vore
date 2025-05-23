package files

import (
	"errors"
	"io"
	"os"
)

type BufferedFile struct {
	closed        bool
	file          *os.File
	fileSize      int64
	buffer        []byte
	bufferSize    int64
	minOffset     int64
	maxOffset     int64
	currentOffset int64
}

func NewBufferedFile(file *os.File, fileSize int64) *BufferedFile {
	bufferSize := int64(4096)

	bufferedFile := &BufferedFile{
		closed:        false,
		file:          file,
		fileSize:      fileSize,
		buffer:        make([]byte, bufferSize),
		bufferSize:    bufferSize,
		minOffset:     0,
		currentOffset: 0,
	}

	bytesRead, err := bufferedFile.file.Read(bufferedFile.buffer)
	if err != nil {
		panic(err)
	}
	bufferedFile.maxOffset = int64(bytesRead)

	return bufferedFile
}

func (v *BufferedFile) Read(p []byte) (int, error) {
	if v.closed {
		return 0, io.ErrClosedPipe
	}
	// there is probably a fancier way to do this
	outputOffset := 0
	outputSize := len(p)
	for {
		for v.currentOffset < v.maxOffset && outputOffset < outputSize {
			p[outputOffset] = v.buffer[v.currentOffset-v.minOffset]
			v.currentOffset += 1
			outputOffset += 1
		}

		if outputOffset == outputSize {
			break
		}

		_, err := v.Seek(0, io.SeekCurrent) // we seek to where we are now to recenter the buffer
		if err != nil {
			return outputOffset, err
		}
	}

	return outputOffset, nil
}

func (v *BufferedFile) Seek(offset int64, whence int) (int64, error) {
	if v.closed {
		return 0, io.ErrClosedPipe
	}
	newOffset := v.currentOffset
	if whence == io.SeekStart {
		newOffset = offset
	} else if whence == io.SeekCurrent {
		newOffset += offset
	} else if whence == io.SeekEnd {
		return v.currentOffset, errors.New("TODO seek from end of file not implemented")
	}

	if newOffset == -1 {
		return v.currentOffset, errors.New("seeking to negative file offset")
	}

	if newOffset < v.minOffset || newOffset >= v.maxOffset {
		newStart := newOffset - (v.bufferSize / 2)
		if newStart < 0 {
			newStart = 0
		}

		fileBound := v.fileSize - 4096
		if fileBound < 0 {
			fileBound = v.fileSize
		}
		if newStart >= fileBound {
			newStart = v.fileSize - 4096
			if newStart < 0 {
				newStart = 0
			}
		}

		bytesRead, err := v.file.ReadAt(v.buffer, newStart)
		// it is actually expected to have an EOF error here when we are working with a file that is less than 4096 bytes
		if err != nil && err != io.EOF {
			return v.currentOffset, err
		}
		v.minOffset = newStart
		v.maxOffset = newStart + int64(bytesRead)
	}

	v.currentOffset = newOffset
	return v.currentOffset, nil
}

func (v *BufferedFile) Close() error {
	if v.closed {
		return io.ErrClosedPipe
	}
	v.closed = true
	return v.file.Close()
}
