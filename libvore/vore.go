package libvore

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Vore struct {
	tokens []*Token
}

func Compile(command string) Vore {
	return compile(strings.NewReader(command))
}

func CompileFile(source string) Vore {
	dat, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	return compile(dat)
}

func compile(reader io.Reader) Vore {
	lexer := initLexer(reader)

	tokens := lexer.getTokens()

	return Vore{tokens}
}

func (v Vore) PrintTokens() {
	for _, token := range v.tokens {
		fmt.Printf("[%s] '%s'\n", token.tokenType.pp(), token.lexeme)
	}
}
