package files

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestWriterTest(t *testing.T) {
	filename := testutils.GetTestingFilename(t, true)

	writer := WriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer was initialized as nil :(")
	}

	writer.WriteAt(0, "my data :)")
	writer.Close()

	reader := ReaderFromFile(filename)
	if reader == nil {
		t.Errorf("Reader was initialized as nil :(")
	}

	value := reader.Read(10)
	if len(value) != 10 {
		t.Errorf("Expected to read 10 characters but actual read %d", len(value))
	}

	if value != "my data :)" {
		t.Errorf("Expected to read 'my data :)' from file but actually read '%s'", value)
	}

	reader.Close()
	testutils.RemoveTestingFile(t, filename)
}

func TestWriterPanicTest(t *testing.T) {
	filename := testutils.GetTestingFilename(t, true)

	writer := WriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer was initialized as nil :(")
	}

	testutils.MustPanic(t, "Expected panic when writing 'my data :)' to -7", func(t *testing.T) {
		writer.WriteAt(-7, "my data :)")
	})

	writer.Close()

	// TODO I can't figure out how to hit the other panics in the VWriter

	testutils.MustPanic(t, "Expected panic when closing an already closed writer", func(t *testing.T) {
		writer.Close()
	})

	testutils.RemoveTestingFile(t, filename)
}

func TestReaderTest(t *testing.T) {
	filename := testutils.GetTestingFilename(t, true)
	writer := WriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer initialized as nil :(")
	}
	writer.WriteAt(0, "hello world")
	writer.Close()

	reader := ReaderFromFile(filename)

	testutils.MustPanic(t, "Expected panic when seeking to -1", func(t *testing.T) {
		reader.Seek(-1)
	})

	testutils.MustPanic(t, "Expected panic when reading a negative length string", func(t *testing.T) {
		reader.Read(-1)
	})

	value := reader.ReadAt(5, 8)
	if value != "" {
		t.Errorf("Expected to read an empty string since we read past the end of the buffer but we got '%s'", value)
	}

	reader.Seek(0)

	reader.Close()

	testutils.MustPanic(t, "Expected panic when read from an already closed reader", func(t *testing.T) {
		value := reader.Read(5)
		t.Errorf("Value was returned as '%s' (offset: %d, size: %d)", value, reader.offset, reader.size)
	})

	testutils.MustPanic(t, "Expected panic when read at from an already closed reader", func(t *testing.T) {
		value := reader.ReadAt(5, 0)
		t.Errorf("Value was returned as '%s' (offset: %d, size: %d)", value, reader.offset, reader.size)
	})

	testutils.MustPanic(t, "Expected panic when closing an already closed reader", func(t *testing.T) {
		reader.Close()
	})

	testutils.RemoveTestingFile(t, filename)
}

func TestReaderPanicTest(t *testing.T) {
	filename := testutils.GetTestingFilename(t, false)

	testutils.MustPanic(t, "Expected to panic when opening a file that is supposed to not exist", func(t *testing.T) {
		_ = ReaderFromFileToMemory(filename)
	})

	testutils.MustPanic(t, "Expected to panic when opening a file that is supposed to not exist", func(t *testing.T) {
		_ = ReaderFromFile(filename)
	})
}
