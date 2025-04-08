package libvore

import (
	"os"
	"strings"

	"github.com/jmeaster30/vore/libvore/algo"
)

type PathEntryType int

const (
	Directory PathEntryType = iota
	WildcardDirectory
	File
	WildcardFile
)

type PathEntry struct {
	entryType PathEntryType
	value     string
}

type Path struct {
	entries []PathEntry
}

func ParsePath(path string) *Path {
	var entries []PathEntry
	if path[0] == '/' {
		entries = append(entries, PathEntry{entryType: Directory, value: "/"})
		path = path[1:]
	}
	splitPath := strings.Split(path, "/")
	for idx, pathPart := range splitPath {
		if idx == len(splitPath) {
			if strings.ContainsRune(pathPart, '*') {
				entries = append(entries, PathEntry{entryType: WildcardFile, value: pathPart})
			} else {
				entries = append(entries, PathEntry{entryType: File, value: pathPart})
			}
		} else {
			if strings.ContainsRune(pathPart, '*') {
				entries = append(entries, PathEntry{entryType: WildcardDirectory, value: pathPart})
			} else {
				entries = append(entries, PathEntry{entryType: Directory, value: pathPart})
			}
		}
	}
	return &Path{entries}
}

func pathMatches(target string, matches string) bool {
	if !strings.ContainsRune(matches, '*') {
		return target == matches
	}

	matchParts := algo.Window(algo.SplitKeep(matches, "*"), 2)

	result := true
	for _, part := range matchParts {
		if len(part) == 1 {
			if part[0] != "*" && target != part[0] {
				result = false
			}
			break
		} else if part[0] == "*" {
			splitStart := strings.Index(target, part[1])
			if splitStart == -1 {
				target = ""
			} else {
				target = target[splitStart:]
			}
		} else if strings.HasPrefix(target, part[0]) {
			target = strings.TrimPrefix(target, part[0])
		} else {
			result = false
			break
		}
	}
	return result
}

func directoryExists(entries []os.DirEntry, name string) bool {
	for _, e := range entries {
		if e.Name() == name {
			return true
		}
	}
	return false
}

func (path *Path) shrink() *Path {
	return &Path{entries: path.entries[1:]}
}

func (path *Path) GetFileList(currentDirectory string) []string {
	if len(path.entries) == 1 {
		entries, err := os.ReadDir(currentDirectory)
		if err != nil {
			return []string{}
		}

		var results []string
		for _, e := range entries {
			if !e.IsDir() && pathMatches(e.Name(), path.entries[0].value) {
				results = append(results, currentDirectory+"/"+e.Name())
			}
		}
		return results
	}

	if path.entries[0].value == "/" {
		return path.shrink().GetFileList("/")
	}

	entries, err := os.ReadDir(currentDirectory)
	if err != nil {
		return []string{}
	}

	var results []string

	if strings.Trim(path.entries[0].value, "*") == "" {
		results = append(results, path.shrink().GetFileList(currentDirectory)...)
	}

	if path.entries[0].entryType == WildcardDirectory {
		for _, e := range entries {
			if pathMatches(e.Name(), path.entries[0].value) {
				results = append(results, path.shrink().GetFileList(currentDirectory+"/"+e.Name())...)
			}
		}
	} else if path.entries[0].entryType == Directory && directoryExists(entries, path.entries[0].value) {
		results = append(results, path.shrink().GetFileList(currentDirectory+"/"+path.entries[0].value)...)
	}

	return results
}
