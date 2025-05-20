package files

import (
	"os"
	"os/user"
	"strings"

	"github.com/jmeaster30/vore/libvore/algo"
)

type Path struct {
	entries []string
}

func ParsePath(path string) *Path {
	var entries []string
	if path[0] == '/' {
		entries = append(entries, "/")
		path = path[1:]
	}
	splitPath := strings.Split(path, "/")
	for _, pathPart := range splitPath {
		if strings.ContainsRune(pathPart, '*') {
			entries = append(entries, pathPart)
		} else if pathPart == "." {
			continue
		} else if pathPart == ".." {
			if len(entries) == 0 || entries[len(entries)-1] == ".." {
				entries = append(entries, pathPart)
			} else {
				entries = entries[:len(entries)-1]
			}
		} else {
			entries = append(entries, pathPart)
		}
	}
	return &Path{entries}
}

func matchIndexes(target string, match string) []int {
	result := []int{}
	for {
		splitOffset := strings.Index(target, match)
		if splitOffset == -1 {
			break
		}
		if len(result) == 0 {
			result = append(result, splitOffset)
		} else {
			result = append(result, result[len(result)-1]+1+splitOffset)
		}
		target = target[splitOffset+1:]
	}
	return result
}

func pathMatchesRecurse(target string, matchPartIdx int, matchParts [][]string) bool {
	part := matchParts[matchPartIdx]
	if len(part) == 1 {
		return part[0] == "*" || target == part[0]
	} else if part[0] == "*" {
		splits := matchIndexes(target, part[1])
		for _, split := range splits {
			if pathMatchesRecurse(target[split:], matchPartIdx+1, matchParts) {
				return true
			}
		}
	} else if strings.HasPrefix(target, part[0]) {
		target = strings.TrimPrefix(target, part[0])
		matchPartIdx++
		return pathMatchesRecurse(target, matchPartIdx, matchParts)
	}
	return false
}

func pathMatches(target string, matches string) bool {
	if !strings.ContainsRune(matches, '*') {
		return target == matches
	}

	matchParts := algo.Window(algo.SplitKeep(matches, "*"), 2)

	return pathMatchesRecurse(target, 0, matchParts)
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

func normalize(path string) string {
	if path == "" {
		return path
	}
	parts := strings.Split(path, "/")

	finalPath := []string{}
	for _, part := range parts {
		if part == ".." && len(finalPath) == 0 {
			finalPath = append(finalPath, part)
		} else if part == ".." && len(finalPath) >= 1 && finalPath[len(finalPath)-1] == "/" {
			continue
		} else if part == ".." {
			finalPath = finalPath[:len(finalPath)-1]
		} else if part == "." {
			continue
		} else if part == "~" {
			current, err := user.Current()
			if err != nil {
				panic(err)
			}
			finalPath = []string{"/", "home", current.Username}
		} else if part == "" {
			finalPath = append(finalPath, "/")
		} else {
			finalPath = append(finalPath, part)
		}
	}

	result := ""
	for idx, part := range finalPath {
		result += part
		if part != "/" && idx != len(finalPath)-1 {
			result += "/"
		}
	}
	return result
}

func (path *Path) GetFileList(currentDirectory string) []string {
	normCurrentDirectory := normalize(currentDirectory)

	if len(path.entries) == 1 {
		entries, err := os.ReadDir(normCurrentDirectory)
		if err != nil {
			return []string{}
		}

		var results []string
		for _, e := range entries {
			if !e.IsDir() && pathMatches(e.Name(), path.entries[0]) {
				results = append(results, normCurrentDirectory+"/"+e.Name())
			}
		}
		return results
	}

	entries, err := os.ReadDir(normCurrentDirectory)
	if err != nil {
		return []string{}
	}

	var results []string

	if strings.Trim(path.entries[0], "*") == "" {
		results = append(results, path.shrink().GetFileList(normCurrentDirectory)...)
	}

	isWildcard := strings.Contains(path.entries[0], "*")

	if isWildcard {
		for _, e := range entries {
			if pathMatches(e.Name(), path.entries[0]) {
				results = append(results, path.shrink().GetFileList(normCurrentDirectory+"/"+e.Name())...)
			}
		}
	} else if !isWildcard && directoryExists(entries, path.entries[0]) {
		results = append(results, path.shrink().GetFileList(normCurrentDirectory+"/"+path.entries[0])...)
	}

	return results
}
