package libvore

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Vore struct {
	tokens   []*Token
	commands []AstCommand
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

	return Vore{tokens, commands}
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
}
