package libvore

import (
	"io"
	"testing"
)

func TestWriteMemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	value := string(stream.contents)
	if value != "lets test!" {
		t.Errorf("Value was expected to be 'lets test!' but was '%s'", value)
	}
}

func TestSeekStartMemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	pos, err := stream.Seek(3, io.SeekStart)
	checkNoError(t, err)
	if pos != 3 {
		t.Errorf("Returned pos was expected to be 3 but was %d", pos)
	}

	if stream.pos != 3 {
		t.Errorf("Stream pos was expected to be 3 but was %d", pos)
	}
}

func TestSeekStartNegativeMemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	pos, err := stream.Seek(-5, io.SeekStart)
	if err == nil {
		t.Errorf("Got no error but expected an error :(")
	}

	if err.Error() != "negative result pos" {
		t.Errorf("Mismatch error message expected 'negative result pos' but go '%s'", err.Error())
	}

	if pos != 10 {
		t.Errorf("Returned pos was expected to be 10 but was %d", pos)
	}

	if stream.pos != 10 {
		t.Errorf("Stream pos was expected to be 10 but was %d", pos)
	}
}

func TestSeekCurrentMemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	value := string(stream.contents)
	if value != "lets test!" {
		t.Errorf("Value was expected to be 'lets test!' but was '%s'", value)
	}

	pos, err := stream.Seek(3, io.SeekCurrent)
	checkNoError(t, err)
	if pos != 13 {
		t.Errorf("Returned pos was expected to be 13 but was %d", pos)
	}

	if stream.pos != 13 {
		t.Errorf("Stream pos was expected to be 13 but was %d", pos)
	}
}

func TestSeekCurrent2MemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	value := string(stream.contents)
	if value != "lets test!" {
		t.Errorf("Value was expected to be 'lets test!' but was '%s'", value)
	}

	_, err := stream.Seek(3, io.SeekStart)
	checkNoError(t, err)

	pos, err := stream.Seek(4, io.SeekCurrent)
	checkNoError(t, err)
	if pos != 7 {
		t.Errorf("Returned pos was expected to be 7 but was %d", pos)
	}

	if stream.pos != 7 {
		t.Errorf("Stream pos was expected to be 7 but was %d", pos)
	}
}

func TestSeekEndMemoryStream(t *testing.T) {
	stream := NewMemoryStream()

	if stream.pos != 0 && len(stream.contents) != 0 {
		t.Errorf("MemoryStream cursor position not 0")
	}

	if len(stream.contents) != 0 {
		t.Errorf("MemoryStream contents not initialized to zero length")
	}

	stream.Write([]byte("lets test!")) // 10

	if stream.pos != 10 {
		t.Errorf("Position expected to be at 10 but was at %d", stream.pos)
	}

	value := string(stream.contents)
	if value != "lets test!" {
		t.Errorf("Value was expected to be 'lets test!' but was '%s'", value)
	}

	pos, err := stream.Seek(3, io.SeekEnd)
	checkNoError(t, err)
	if pos != 13 {
		t.Errorf("Returned pos was expected to be 13 but was %d", pos)
	}

	if stream.pos != 13 {
		t.Errorf("Stream pos was expected to be 13 but was %d", pos)
	}
}
