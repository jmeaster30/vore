package engine

import (
	"os"

	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/files"
)

func Run(bytecode *bytecode.Bytecode, searchText string) Matches {
	result := Matches{}
	for _, command := range bytecode.Bytecode {
		reader := files.ReaderFromString(searchText)
		result = append(result, search(&command, "text", reader, NOTHING)...)
	}
	return result
}

func RunFiles(bytecode *bytecode.Bytecode, filenames []string, mode ReplaceMode, processFilenames bool) Matches {
	actualMode := mode
	if processFilenames {
		actualMode = NOTHING
	}
	result := Matches{}
	for _, command := range bytecode.Bytecode {
		// command.print()
		for _, filename := range filenames {
			actualFiles := []string{}
			info, err := os.Stat(filename)
			if err != nil {
				panic(err)
			}
			fixedFilename := filename
			if info.IsDir() {
				if filename[len(filename)-1] != '/' || filename[len(filename)-1] != '\\' {
					fixedFilename += "/"
				}
				entries, err := os.ReadDir(filename)
				if err != nil {
					panic(err)
				}
				for _, entry := range entries {
					actualFiles = append(actualFiles, fixedFilename+entry.Name())
				}
			} else {
				actualFiles = append(actualFiles, fixedFilename)
			}
			for _, actualFilename := range actualFiles {
				var reader *files.Reader
				if processFilenames {
					reader = files.ReaderFromString(actualFilename)
				} else {
					reader = files.ReaderFromFile(actualFilename)
				}
				foundMatches := search(&command, actualFilename, reader, actualMode)
				result = append(result, foundMatches...)
				if processFilenames && len(foundMatches) != 0 && len(foundMatches[0].Replacement.GetValueOrDefault("")) != 0 {
					err := os.Rename(actualFilename, foundMatches[0].Replacement.GetValueOrDefault(""))
					if err != nil {
						os.Stderr.WriteString("Failed to rename file '" + actualFilename + "' to '" + foundMatches[0].Replacement.GetValueOrDefault("") + "'\n")
					}
				}
			}
		}
	}
	return result
}
