package libvore

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/engine"
)

type Vore struct {
	ast      *ast.Ast
	bytecode *bytecode.Bytecode
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

func (v *Vore) Run(searchText string) engine.Matches {
	return engine.Run(v.bytecode, searchText)
}

func (v *Vore) RunFiles(filenames []string, mode engine.ReplaceMode, processFilenames bool) engine.Matches {
	return engine.RunFiles(v.bytecode, filenames, mode, processFilenames)
}

func (v *Vore) PrintAST() {
	for _, command := range v.ast.Commands() {
		command.Print()
	}
	fmt.Println()
}

func (v *Vore) PrintBytecode() {
	for _, command := range v.bytecode.Bytecode {
		fmt.Printf("%s\n", command)
	}
}
