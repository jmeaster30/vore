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
	s.astNode.PrintNode() // TODO turn PrintNode() into NodeString()
	return fmt.Sprintf("SemanticError: %s", s.message)
}

func NewSemanticError(node ast.AstNode, message string) *SemanticError {
	return &SemanticError{node, message}
}
