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
	return fmt.Sprintf("GenError: %s at node %s", g.message, g.astNode.NodeString())
}

func (g *GenError) Message() string {
	return g.Error()
}

func NewGenError(node ast.AstNode, msg string) *GenError {
	return &GenError{node, msg}
}
