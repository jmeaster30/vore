package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type GenError struct {
	astNode ast.AstNode
	message string
}

func (g *GenError) Error() string {
	g.astNode.PrintNode()
	return fmt.Sprintf("GenError: %s", g.message)
}

func (g *GenError) Message() string {
	return g.Error()
}

func NewGenError(node ast.AstNode, msg string) *GenError {
	return &GenError{node, msg}
}
