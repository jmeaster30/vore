package bytecode

import "github.com/jmeaster30/vore/libvore/ast"

type GenerateProcessInfo struct{}

func generateProcessBytecode(statements []ast.AstProcessStatement) []ProcInstruction {
	result := []ProcInstruction{}
	for _, statement := range statements {
		insts := generateProcessStatement(statement)
		result = append(result, insts...)
	}
	return result
}

func generateProcessStatement(statement ast.AstProcessStatement) []ProcInstruction {
	var st any = statement
	switch stmt := st.(type) {
	case *ast.AstProcessSet:
		return generateProcessSet(stmt)
	case *ast.AstProcessBreak:
		return generateProcessBreak(stmt)
	case *ast.AstProcessContinue:
		return generateProcessContinue(stmt)
	case *ast.AstProcessDebug:
		return generateProcessDebug(stmt)
	case *ast.AstProcessIf:
		return generateProcessIf(stmt)
	case *ast.AstProcessLoop:
		return generateProcessLoop(stmt)
	case *ast.AstProcessReturn:
		return generateProcessReturn(stmt)
	}
	panic("Bad instruction statement")
}

func generateProcessSet(setStmt *ast.AstProcessSet) []ProcInstruction {
	expr := generateProcessExpression(setStmt.Expr)
	store := Store{setStmt.Name}
	return append(expr, store)
}

func generateProcessBreak(breakStmt *ast.AstProcessBreak) []ProcInstruction {
	return []ProcInstruction{}
}

func generateProcessContinue(continueStmt *ast.AstProcessContinue) []ProcInstruction {
	return []ProcInstruction{}
}

func generateProcessDebug(debug *ast.AstProcessDebug) []ProcInstruction {
	insts := generateProcessExpression(debug.Expr)
	return append(insts, Debug{})
}

func generateProcessIf(ifStmt *ast.AstProcessIf) []ProcInstruction {
	return []ProcInstruction{}
}

func generateProcessLoop(loop *ast.AstProcessLoop) []ProcInstruction {
	return []ProcInstruction{}
}

func generateProcessReturn(returnStmt *ast.AstProcessReturn) []ProcInstruction {
	return []ProcInstruction{}
}

func generateProcessExpression(expression ast.AstProcessExpression) []ProcInstruction {
	var ex any = expression
	switch expr := ex.(type) {
	case *ast.AstProcessUnaryExpression:
		return generateProcessUnaryExpression(expr)
	case *ast.AstProcessBinaryExpression:
		return generateProcessBinaryExpression(expr)
	case *ast.AstProcessBoolean:
		return generateProcessBoolean(expr)
	case *ast.AstProcessString:
		return generateProcessString(expr)
	case *ast.AstProcessNumber:
		return generateProcessNumber(expr)
	case *ast.AstProcessVariable:
		return generateProcessVariable(expr)
	}
	panic("Bad expression type")
}

func generateProcessUnaryExpression(unary *ast.AstProcessUnaryExpression) []ProcInstruction {
	val := generateProcessExpression(unary.Expr)
	var unaryInst ProcInstruction
	switch unary.Op {
	case ast.HEAD:
		unaryInst = Head{}
	case ast.TAIL:
		unaryInst = Tail{}
	case ast.NOT:
		unaryInst = Not{}
	default:
		panic("BAD UNARY INST")
	}

	return append(val, unaryInst)
}

func generateProcessBinaryExpression(binary *ast.AstProcessBinaryExpression) []ProcInstruction {
	left := generateProcessExpression(binary.Lhs)
	right := generateProcessExpression(binary.Rhs)
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
		panic("Bad binary op")
	}
	return append(append(left, right...), binaryInst)
}

func generateProcessBoolean(boolean *ast.AstProcessBoolean) []ProcInstruction {
	return []ProcInstruction{Push{BooleanValue{boolean.Value}}}
}

func generateProcessString(stringValue *ast.AstProcessString) []ProcInstruction {
	return []ProcInstruction{Push{StringValue{stringValue.Value}}}
}

func generateProcessNumber(number *ast.AstProcessNumber) []ProcInstruction {
	return []ProcInstruction{Push{NumberValue{number.Value}}}
}

func generateProcessVariable(variable *ast.AstProcessVariable) []ProcInstruction {
	return []ProcInstruction{Load{variable.Name}}
}
