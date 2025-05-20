package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/ds"
)

type LoopInfo struct {
	breakLabel    string
	continueLabel string
}

type GenerateProcessInfo struct {
	loopStack                *ds.Stack[LoopInfo]
	currentInstructionOffset int
}

func generateProcessBytecode(statements []ast.AstProcessStatement, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	result := []ProcInstruction{}
	currentInfo := info
	if currentInfo == nil {
		currentInfo = &GenerateProcessInfo{
			ds.NewStack[LoopInfo](),
			0,
		}
	}

	baseOffset := currentInfo.currentInstructionOffset

	for _, statement := range statements {
		currentInfo.currentInstructionOffset = baseOffset + len(result)
		insts, err := generateProcessStatement(statement, currentInfo)
		if err != nil {
			return nil, err
		}
		result = append(result, insts...)
	}
	return result, nil
}

func generateProcessStatement(statement ast.AstProcessStatement, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	var st any = statement
	switch stmt := st.(type) {
	case *ast.AstProcessSet:
		return generateProcessSet(stmt, info)
	case *ast.AstProcessBreak:
		return generateProcessBreak(stmt, info)
	case *ast.AstProcessContinue:
		return generateProcessContinue(stmt, info)
	case *ast.AstProcessDebug:
		return generateProcessDebug(stmt, info)
	case *ast.AstProcessIf:
		return generateProcessIf(stmt, info)
	case *ast.AstProcessLoop:
		return generateProcessLoop(stmt, info)
	case *ast.AstProcessReturn:
		return generateProcessReturn(stmt, info)
	}
	return nil, NewGenError(statement, "unknown process type")
}

func generateProcessSet(setStmt *ast.AstProcessSet, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	expr, err := generateProcessExpression(setStmt.Expr, info)
	if err != nil {
		return nil, err
	}
	store := Store{setStmt.Name}
	return append(expr, store), nil
}

func generateProcessBreak(breakStmt *ast.AstProcessBreak, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	loopInfo := info.loopStack.Peek()
	if !loopInfo.HasValue() {
		return nil, NewGenError(*breakStmt, "Break statement not allowed outside of loop stack")
	}

	return []ProcInstruction{LabelJump{loopInfo.GetValue().breakLabel}}, nil
}

func generateProcessContinue(continueStmt *ast.AstProcessContinue, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	loopInfo := info.loopStack.Peek()
	if !loopInfo.HasValue() {
		return nil, NewGenError(*continueStmt, "Continue statement not allowed outside of loop stack")
	}

	return []ProcInstruction{LabelJump{loopInfo.GetValue().continueLabel}}, nil
}

func generateProcessDebug(debug *ast.AstProcessDebug, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	insts, err := generateProcessExpression(debug.Expr, info)
	if err != nil {
		return nil, err
	}
	return append(insts, Debug{}), nil
}

func generateProcessIf(ifStmt *ast.AstProcessIf, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	condition, err := generateProcessExpression(ifStmt.Condition, info)
	if err != nil {
		return nil, err
	}

	info.currentInstructionOffset += len(condition) + 1
	trueblock, err := generateProcessBytecode(ifStmt.TrueBody, info)
	if err != nil {
		return nil, err
	}

	info.currentInstructionOffset += len(trueblock) + 1
	conditionalJump := ConditionalJump{info.currentInstructionOffset}
	falseBlock, err := generateProcessBytecode(ifStmt.FalseBody, info)
	if err != nil {
		return nil, err
	}
	info.currentInstructionOffset += len(falseBlock)
	endJump := Jump{info.currentInstructionOffset}

	result := append(condition, conditionalJump)
	result = append(result, trueblock...)
	result = append(result, endJump)
	result = append(result, falseBlock...)
	return result, nil
}

func generateProcessLoop(loop *ast.AstProcessLoop, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	info.loopStack.Push(LoopInfo{
		breakLabel:    fmt.Sprintf("loopAt%dBreak", info.currentInstructionOffset),
		continueLabel: fmt.Sprintf("loopAt%dContinue", info.currentInstructionOffset),
	})
	info.currentInstructionOffset = info.currentInstructionOffset + 1
	result, err := generateProcessBytecode(loop.Body, info)
	if err != nil {
		return nil, err
	}

	loopInfo := info.loopStack.Pop()

	result = append([]ProcInstruction{Label{loopInfo.GetValue().continueLabel}}, result...)
	result = append(result, LabelJump{loopInfo.GetValue().continueLabel})
	result = append(result, Label{loopInfo.GetValue().breakLabel})
	return result, nil
}

func generateProcessReturn(returnStmt *ast.AstProcessReturn, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	result, err := generateProcessExpression(returnStmt.Expr, info)
	if err != nil {
		return nil, err
	}

	return append(result, Return{}), nil
}

func generateProcessExpression(expression ast.AstProcessExpression, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	var ex any = expression
	switch expr := ex.(type) {
	case ast.AstProcessUnaryExpression:
		return generateProcessUnaryExpression(&expr, info)
	case ast.AstProcessBinaryExpression:
		return generateProcessBinaryExpression(&expr, info)
	case ast.AstProcessBoolean:
		return generateProcessBoolean(expr, info)
	case ast.AstProcessString:
		return generateProcessString(expr, info)
	case ast.AstProcessNumber:
		return generateProcessNumber(expr, info)
	case ast.AstProcessVariable:
		return generateProcessVariable(expr, info)
	}
	return nil, NewGenError(expression, "unknown expression type")
}

func generateProcessUnaryExpression(unary *ast.AstProcessUnaryExpression, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	val, err := generateProcessExpression(unary.Expr, info)
	if err != nil {
		return nil, err
	}
	var unaryInst ProcInstruction
	switch unary.Op {
	case ast.HEAD:
		unaryInst = Head{}
	case ast.TAIL:
		unaryInst = Tail{}
	case ast.NOT:
		unaryInst = Not{}
	default:
		return nil, NewGenError(*unary, "unknown unary expression type")
	}

	return append(val, unaryInst), nil
}

func generateProcessBinaryExpression(binary *ast.AstProcessBinaryExpression, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	left, err := generateProcessExpression(binary.Lhs, info)
	if err != nil {
		return nil, err
	}
	right, err := generateProcessExpression(binary.Rhs, info)
	if err != nil {
		return nil, err
	}
	var binaryInst ProcInstruction
	switch binary.Op {
	case ast.PLUS:
		binaryInst = Add{}
	case ast.MINUS:
		binaryInst = Subtract{}
	case ast.MULT:
		binaryInst = Multiply{}
	case ast.DIV:
		binaryInst = Divide{}
	case ast.MOD:
		binaryInst = Modulo{}
	case ast.LESS:
		binaryInst = LessThan{}
	case ast.LESSEQ:
		binaryInst = LessThanEqual{}
	case ast.GREATER:
		binaryInst = GreaterThan{}
	case ast.GREATEREQ:
		binaryInst = GreaterThanEqual{}
	case ast.DEQUAL:
		binaryInst = Equal{}
	case ast.NEQUAL:
		binaryInst = NotEqual{}
	default:
		return nil, NewGenError(*binary, "unknown binary instruction type")
	}
	return append(append(left, right...), binaryInst), nil
}

func generateProcessBoolean(boolean ast.AstProcessBoolean, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	return []ProcInstruction{Push{BooleanValue{boolean.Value}}}, nil
}

func generateProcessString(stringValue ast.AstProcessString, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	return []ProcInstruction{Push{StringValue{stringValue.Value}}}, nil
}

func generateProcessNumber(number ast.AstProcessNumber, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	return []ProcInstruction{Push{NumberValue{number.Value}}}, nil
}

func generateProcessVariable(variable ast.AstProcessVariable, info *GenerateProcessInfo) ([]ProcInstruction, error) {
	return []ProcInstruction{Load{variable.Name}}, nil
}
