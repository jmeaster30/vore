package libvore

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Match struct {
	filename     string
	matchNumber  int
	fileOffset   Range
	lineNumber   Range
	columnNumber Range
	value        string
	variables    map[string]string
}

func (m Match) Print() {
	fmt.Println("============")
	fmt.Printf("Filename: %s\n", m.filename)
	fmt.Printf("MatchNumber: %d\n", m.matchNumber)
	fmt.Printf("Value: %s\n", m.value)
	fmt.Printf("FileOffset: %d %d\n", m.fileOffset.Start, m.fileOffset.End)
	fmt.Printf("Line: %d %d\n", m.lineNumber.Start, m.lineNumber.End)
	fmt.Printf("Column: %d %d\n", m.columnNumber.Start, m.columnNumber.End)
	fmt.Println("Variables:")
	fmt.Println("\t[key] = [value]")
	for key, value := range m.variables {
		fmt.Printf("\t%s = %s\n", key, value)
	}
	fmt.Println("============")
}

type Vore struct {
	tokens   []*Token
	commands []AstCommand
	bytecode []Command
}

func Compile(command string) Vore {
	return compile("source", strings.NewReader(command))
}

func CompileFile(source string) Vore {
	dat, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	return compile(source, dat)
}

func compile(filename string, reader io.Reader) Vore {
	lexer := initLexer(reader)

	tokens := lexer.getTokens()
	commands, parseError := parse(tokens)
	if parseError.isError {
		panic(fmt.Sprintf("\nERROR:  %s\nToken:  '%s'\nLine:   %d - %d\nColumn: %d - %d\n", parseError.message, parseError.token.lexeme, parseError.token.line.Start, parseError.token.line.End, parseError.token.column.Start, parseError.token.column.End))
	}

	bytecode := []Command{}
	for _, ast_comm := range commands {
		byte_comm := ast_comm.generate()
		bytecode = append(bytecode, byte_comm)
	}

	return Vore{tokens, commands, bytecode}
}

func (v *Vore) Run(filenames []string) []Match {
	result := []Match{}
	for _, command := range v.bytecode {
		command.print()
		for _, filename := range filenames {
			result = append(result, command.execute(filename)...)
		}
	}
	return result
}

func (v *Vore) PrintTokens() {
	for _, token := range v.tokens {
		fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
	}
}

func (v *Vore) PrintAST() {
	for _, command := range v.commands {
		command.print()
	}
	fmt.Println()
}
