package libvore

import "testing"

func TestVWriterTest(t *testing.T) {
	filename := getTestingFilename(t, true)

	writer := VWriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer was initialized as nil :(")
	}

	writer.WriteAt(0, "my data :)")
	writer.Close()

	reader := VReaderFromFile(filename)
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
	removeTestingFile(t, filename)
}

func TestVWriterPanicTest(t *testing.T) {
	filename := getTestingFilename(t, true)

	writer := VWriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer was initialized as nil :(")
	}

	mustPanic(t, "Expected panic when writing 'my data :)' to -7", func(t *testing.T) {
		writer.WriteAt(-7, "my data :)")
	})

	writer.Close()

	// TODO I can't figure out how to hit the other panics in the VWriter

	mustPanic(t, "Expected panic when closing an already closed writer", func(t *testing.T) {
		writer.Close()
	})

	removeTestingFile(t, filename)
}

func TestVReaderTest(t *testing.T) {
	filename := getTestingFilename(t, true)
	writer := VWriterFromFile(filename)
	if writer == nil {
		t.Errorf("Writer initialized as nil :(")
	}
	writer.WriteAt(0, "hello world")
	writer.Close()

	reader := VReaderFromFile(filename)

	mustPanic(t, "Expected panic when seeking to -1", func(t *testing.T) {
		reader.Seek(-1)
	})

	mustPanic(t, "Expected panic when reading a negative length string", func(t *testing.T) {
		reader.Read(-1)
	})

	value := reader.ReadAt(5, 8)
	if value != "" {
		t.Errorf("Expected to read an empty string since we read past the end of the buffer but we got '%s'", value)
	}

	reader.Seek(0)

	reader.Close()

	mustPanic(t, "Expected panic when read from an already closed reader", func(t *testing.T) {
		value := reader.Read(5)
		t.Errorf("Value was returned as '%s' (offset: %d, size: %d)", value, reader.offset, reader.size)
	})

	mustPanic(t, "Expected panic when read at from an already closed reader", func(t *testing.T) {
		value := reader.ReadAt(5, 0)
		t.Errorf("Value was returned as '%s' (offset: %d, size: %d)", value, reader.offset, reader.size)
	})

	mustPanic(t, "Expected panic when closing an already closed reader", func(t *testing.T) {
		reader.Close()
	})

	removeTestingFile(t, filename)
}

func TestVReaderPanicTest(t *testing.T) {
	filename := getTestingFilename(t, false)

	mustPanic(t, "Expected to panic when opening a file that is supposed to not exist", func(t *testing.T) {
		_ = VReaderFromFileToMemory(filename)
	})

	mustPanic(t, "Expected to panic when opening a file that is supposed to not exist", func(t *testing.T) {
		_ = VReaderFromFile(filename)
	})
}
