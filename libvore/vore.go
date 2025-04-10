package libvore

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/files"
)

type Vore struct {
	ast      *ast.Ast
	bytecode []bytecode.Command
}

func Compile(command string) (*Vore, error) {
	return compile(strings.NewReader(command))
}

func CompileFile(source string) (*Vore, error) {
	source_file, err := os.Open(source)
	if err != nil {
		return nil, NewVoreFileError(err)
	}
	return compile(source_file)
}

func compile(reader io.Reader) (*Vore, error) {
	commands, err := ast.ParseReader(reader)
	if err != nil {
		return nil, err
	}

	bytecode, err := bytecode.GenerateBytecode(commands)
	if err != nil {
		return nil, err
	}

	return &Vore{commands, bytecode}, nil
}

func (v *Vore) RunFiles(filenames []string, mode ReplaceMode, processFilenames bool) Matches {
	actualMode := mode
	if processFilenames {
		actualMode = NOTHING
	}
	result := Matches{}
	for _, command := range v.bytecode {
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
				foundMatches := command.execute(actualFilename, reader, actualMode)
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

func (v *Vore) Run(searchText string) Matches {
	result := Matches{}
	for _, command := range v.bytecode {
		reader := files.ReaderFromString(searchText)
		result = append(result, command.execute("text", reader, NOTHING)...)
	}
	return result
}

// func (v *Vore) PrintTokens() {
// 	for _, token := range v.tokens {
// 		fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
// 	}
// }

func (v *Vore) PrintAST() {
	for _, command := range v.commands {
		command.print()
	}
	fmt.Println()
}

func (v *Vore) PrintBytecode() {
	for _, command := range v.bytecode {
		fmt.Printf("%s\n", command)
	}
}
