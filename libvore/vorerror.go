package libvore

import ast "github.com/jmeaster30/vore/libvore/ast"

type VoreError interface {
	Token() *ast.Token
	Message() string
}
