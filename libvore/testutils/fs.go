package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func GetTestingFilename(t *testing.T, touchFile bool) string {
	t.Helper()
	filename, err := PseudoUuid(t)
	CheckNoError(t, err)

	fullFilename := filename + ".txt"

	if touchFile {
		if _, err := os.Stat(fullFilename); err == nil {
			t.Errorf("Randomized testing file '%s' already exists. I don't want to overwrite it :(", fullFilename)
		}

		file, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		CheckNoError(t, err)

		err = file.Close()
		CheckNoError(t, err)
	}

	return fullFilename
}

func RemoveTestingFile(t *testing.T, filename string) {
	t.Helper()
	err := os.Remove(filename)
	CheckNoError(t, err)
}

func BuildTestingFilesystem(t *testing.T, paths ...string) string {
	// t.Helper()
	directoryName, err := PseudoUuid(t)
	CheckNoError(t, err)

	if _, err := os.Stat(directoryName); err == nil {
		t.Errorf("Randomized testing directory '%s' already exists. I don't want to overwrite it :(", directoryName)
	}

	err = os.Mkdir(directoryName, os.FileMode(0777))
	CheckNoError(t, err)

	fixedDirectoryName := strings.TrimSuffix(directoryName, "/")

	for _, path := range paths {
		fixedPath := fixedDirectoryName + "/" + strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
		parentPath := filepath.Dir(fixedPath)
		err := os.MkdirAll(parentPath, os.FileMode(0777))
		CheckNoError(t, err)

		file, err := os.OpenFile(fixedPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		CheckNoError(t, err)

		err = file.Close()
		CheckNoError(t, err)
	}

	return fixedDirectoryName
}

func RemoveTestingFilesystem(t *testing.T, path string) {
	t.Helper()

	err := os.RemoveAll(path)
	CheckNoError(t, err)
}
