package files

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestParsePathSingleLeaf(t *testing.T) {
	path := ParsePath("test.txt")

	testutils.AssertLength(t, 1, path.entries)
	testutils.AssertEqual(t, "test.txt", path.entries[0])
}

func TestParsePathWildcardLeaf(t *testing.T) {
	path := ParsePath("test.*")

	testutils.AssertLength(t, 1, path.entries)
	testutils.AssertEqual(t, "test.*", path.entries[0])
}

func TestParsePathBranch(t *testing.T) {
	path := ParsePath("test/a.txt")

	testutils.AssertLength(t, 2, path.entries)
	testutils.AssertEqual(t, "test", path.entries[0])
	testutils.AssertEqual(t, "a.txt", path.entries[1])
}

func TestParsePathWildcardLeafAll(t *testing.T) {
	path := ParsePath("test/*")

	testutils.AssertLength(t, 2, path.entries)
	testutils.AssertEqual(t, "test", path.entries[0])
	testutils.AssertEqual(t, "*", path.entries[1])
}

func TestParsePathWildcardBranch(t *testing.T) {
	path := ParsePath("test/*/a.txt")

	testutils.AssertLength(t, 3, path.entries)
	testutils.AssertEqual(t, "test", path.entries[0])
	testutils.AssertEqual(t, "*", path.entries[1])
	testutils.AssertEqual(t, "a.txt", path.entries[2])
}

func TestParsePathFromRoot(t *testing.T) {
	path := ParsePath("/home/vore")

	testutils.AssertLength(t, 3, path.entries)
	testutils.AssertEqual(t, "/", path.entries[0])
	testutils.AssertEqual(t, "home", path.entries[1])
	testutils.AssertEqual(t, "vore", path.entries[2])
}

func TestGetFileListExactFile(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"b.txt",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("a.txt")

	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 1, fileList)
	testutils.AssertEqual(t, fs+"/a.txt", fileList[0])
}

func TestGetFileListAllTextFiles(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"b.txt",
		"c.json",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("*.txt")

	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 2, fileList)
	testutils.AssertEqual(t, fs+"/a.txt", fileList[0])
	testutils.AssertEqual(t, fs+"/b.txt", fileList[1])
}

func TestGetFileListAllTextFilesInFolder(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"test/a.txt",
		"test/b.txt",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("test/*.txt")

	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 2, fileList)
	testutils.AssertEqual(t, fs+"/test/a.txt", fileList[0])
	testutils.AssertEqual(t, fs+"/test/b.txt", fileList[1])
}

func TestGetFileListAllTextFilesInAllFolders(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"folder1/a.txt",
		"folder1/b.txt",
		"folder2/a.txt",
		"folder2/b.txt",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("*/*.txt")

	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 5, fileList)
	testutils.AssertEqual(t, fs+"/a.txt", fileList[0])
	testutils.AssertEqual(t, fs+"/folder1/a.txt", fileList[1])
	testutils.AssertEqual(t, fs+"/folder1/b.txt", fileList[2])
	testutils.AssertEqual(t, fs+"/folder2/a.txt", fileList[3])
	testutils.AssertEqual(t, fs+"/folder2/b.txt", fileList[4])
}

func TestGetFileListSingleTextFileInParentDirectory(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"folder1/a.txt",
		"folder2/a.txt",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("a.txt")

	fileList := path.GetFileList(fs + "/folder1/..")
	testutils.AssertLength(t, 1, fileList)
	testutils.AssertEqual(t, fs+"/a.txt", fileList[0])
}

func TestGetFileListSingelLetterThenWildcard(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"a.txt",
		"abe.txt",
		"cad.mov",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("a*")
	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 2, fileList)
	testutils.AssertEqual(t, fs+"/a.txt", fileList[0])
	testutils.AssertEqual(t, fs+"/abe.txt", fileList[1])
}

func TestGetFileListWildcardThenSingleLetter(t *testing.T) {
	fs := testutils.BuildTestingFilesystem(t,
		"test.a",
		"wow.txt",
		"a.png",
		"plaza",
	)
	defer testutils.RemoveTestingFilesystem(t, fs)
	path := ParsePath("*a")
	fileList := path.GetFileList(fs)
	testutils.AssertLength(t, 2, fileList)
	testutils.AssertEqual(t, fs+"/plaza", fileList[0])
	testutils.AssertEqual(t, fs+"/test.a", fileList[1])
}
