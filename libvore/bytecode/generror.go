package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type GenError struct {
	err error
}

func (g *GenError) Error() string {
	return g.err.Error()
}

func (g *GenError) Token() *ast.Token {
	return nil
}

func (g *GenError) Message() string {
	return g.Error()
}

func NewGenError(msg string) *GenError {
	return &GenError{fmt.Errorf("GenError: %s", msg)}
}
