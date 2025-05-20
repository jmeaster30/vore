package bytecode

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestGenerateProcessBoolean(t *testing.T) {
	ast_boolean := ast.AstProcessBoolean{Value: true}
	insts, err := generateProcessBoolean(ast_boolean, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{BooleanValue{true}}}, insts)
}

func TestGenerateProcessString(t *testing.T) {
	ast_string := ast.AstProcessString{Value: "testString"}
	insts, err := generateProcessString(ast_string, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{StringValue{"testString"}}}, insts)
}

func TestGenerateProcessNumber(t *testing.T) {
	ast_number := ast.AstProcessNumber{Value: 12}
	insts, err := generateProcessNumber(ast_number, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{NumberValue{12}}}, insts)
}

func TestGenerateProcessVariable(t *testing.T) {
	ast_variable := ast.AstProcessVariable{Name: "myVariable"}
	insts, err := generateProcessVariable(ast_variable, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Load{ast_variable.Name}}, insts)
}

func TestGenerateProcessUnaryOp_Not(t *testing.T) {
	ast_unary := &ast.AstProcessUnaryExpression{Op: ast.NOT, Expr: &ast.AstProcessBoolean{Value: true}}
	insts, err := generateProcessUnaryExpression(ast_unary, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{BooleanValue{true}}, Not{}}, insts)
}

func TestGenerateProcessUnaryOp_Head(t *testing.T) {
	ast_unary := &ast.AstProcessUnaryExpression{Op: ast.HEAD, Expr: &ast.AstProcessString{Value: "oh wow"}}
	insts, err := generateProcessUnaryExpression(ast_unary, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{StringValue{"oh wow"}}, Head{}}, insts)
}

func TestGenerateProcessUnaryOp_Tail(t *testing.T) {
	ast_unary := &ast.AstProcessUnaryExpression{Op: ast.TAIL, Expr: &ast.AstProcessString{Value: "super cool"}}
	insts, err := generateProcessUnaryExpression(ast_unary, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{StringValue{"super cool"}}, Tail{}}, insts)
}

func TestGenerateProcessUnaryOp_Nested(t *testing.T) {
	ast_unary := &ast.AstProcessUnaryExpression{
		Op: ast.NOT,
		Expr: &ast.AstProcessUnaryExpression{
			Op:   ast.NOT,
			Expr: &ast.AstProcessBoolean{Value: false},
		},
	}
	insts, err := generateProcessUnaryExpression(ast_unary, nil)
	testutils.CheckNoError(t, err)
	testutils.AssertEqual(t, []ProcInstruction{Push{BooleanValue{false}}, Not{}, Not{}}, insts)
}
