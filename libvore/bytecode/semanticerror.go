package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type SemanticError struct {
	astNode ast.AstNode
	message string
}

func (s SemanticError) Error() string {
	return fmt.Sprintf("SemanticError: %s at %s", s.message, s.astNode.NodeString())
}

func NewSemanticError(node ast.AstNode, message string) *SemanticError {
	return &SemanticError{node, message}
}
