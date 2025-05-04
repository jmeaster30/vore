package libvore

import (
	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/engine"
)

type (
	LexError      ast.LexError
	ParseError    ast.ParseError
	GenError      bytecode.GenError
	SemanticError bytecode.SemanticError
	ExecError     engine.ExecError
)

func ToLexError(err error) ds.Optional[LexError] {
	switch a := err.(type) {
	case *ast.LexError:
		return ds.Some(LexError(*a))
	default:
		return ds.None[LexError]()
	}
}

func ToParseError(err error) ds.Optional[ParseError] {
	switch a := err.(type) {
	case *ast.ParseError:
		return ds.Some(ParseError(*a))
	default:
		return ds.None[ParseError]()
	}
}

func ToGenError(err error) ds.Optional[GenError] {
	switch a := err.(type) {
	case *bytecode.GenError:
		return ds.Some(GenError(*a))
	default:
		return ds.None[GenError]()
	}
}

func ToSemanticError(err error) ds.Optional[SemanticError] {
	switch a := err.(type) {
	case *bytecode.SemanticError:
		return ds.Some(SemanticError(*a))
	default:
		return ds.None[SemanticError]()
	}
}

func ToExecError(err error) ds.Optional[ExecError] {
	switch a := err.(type) {
	case *engine.ExecError:
		return ds.Some(ExecError(*a))
	default:
		return ds.None[ExecError]()
	}
}
