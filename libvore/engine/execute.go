package engine

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/bytecode"
)

type ProcessStatus int

const (
	NEXT ProcessStatus = iota
	BREAKLOOP
	CONTINUELOOP
	RETURNING
)

type ProcessState struct {
	currentValue bytecode.Value
	environment  map[string]bytecode.Value
	status       ProcessStatus
}

func executeStatement(s *ast.AstProcessStatement, state ProcessState) ProcessState {
	var si any = *s
	switch stmt := si.(type) {
	case *ast.AstProcessSet:
		return executeSet(stmt, state)
	case *ast.AstProcessIf:
		return executeIf(stmt, state)
	case *ast.AstProcessLoop:
		return executeLoop(stmt, state)
	case ast.AstProcessBreak:
		return executeBreak(state)
	case ast.AstProcessContinue:
		return executeContinue(state)
	case *ast.AstProcessReturn:
		return executeReturn(stmt, state)
	case *ast.AstProcessDebug:
		return executeDebug(stmt, state)
	}
	panic(fmt.Sprintf("Unknown process statement %T", si))
}

func executeSet(s *ast.AstProcessSet, state ProcessState) ProcessState {
	expr_state := executeExpression(&s.Expr, state)
	expr_state.environment[s.Name] = expr_state.currentValue
	return expr_state
}

func executeReturn(s *ast.AstProcessReturn, state ProcessState) ProcessState {
	expr_state := executeExpression(&s.Expr, state)
	expr_state.status = RETURNING
	return expr_state
}

func executeIf(s *ast.AstProcessIf, state ProcessState) ProcessState {
	expr_state := executeExpression(&s.Condition, state)
	if expr_state.currentValue.GetBoolean() {
		for _, stmt := range s.TrueBody {
			expr_state = executeStatement(&stmt, expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	} else {
		for _, stmt := range s.FalseBody {
			expr_state = executeStatement(&stmt, expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	}

	return expr_state
}

func executeDebug(s *ast.AstProcessDebug, state ProcessState) ProcessState {
	expr_state := executeExpression(&s.Expr, state)
	fmt.Println(expr_state.currentValue.GetString())
	return expr_state
}

func executeLoop(s *ast.AstProcessLoop, state ProcessState) ProcessState {
	expr_state := state
	for {
		for _, stmt := range s.Body {
			expr_state = executeStatement(&stmt, expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
		if expr_state.status == RETURNING || expr_state.status == BREAKLOOP {
			if expr_state.status == BREAKLOOP {
				expr_state.status = NEXT
			}
			break
		} else if expr_state.status == CONTINUELOOP {
			expr_state.status = NEXT
			continue
		}
	}

	return expr_state
}

func executeContinue(state ProcessState) ProcessState {
	state.status = CONTINUELOOP
	return state
}

func executeBreak(state ProcessState) ProcessState {
	state.status = BREAKLOOP
	return state
}

func executeExpression(s *ast.AstProcessExpression, state ProcessState) ProcessState {
	var si any = *s
	switch exp := si.(type) {
	case ast.AstProcessBinaryExpression:
		return executeBinaryExpr(&exp, state)
	case ast.AstProcessUnaryExpression:
		return executeUnaryExpression(&exp, state)
	case ast.AstProcessString:
		return executeString(&exp, state)
	case ast.AstProcessBoolean:
		return executeBoolean(&exp, state)
	case ast.AstProcessNumber:
		return executeNumber(&exp, state)
	case ast.AstProcessVariable:
		return executeVariable(&exp, state)
	}
	panic(fmt.Sprintf("unknown process expression type %T", si))
}

func executeBinaryExpr(s *ast.AstProcessBinaryExpression, state ProcessState) ProcessState {
	lhs_state := executeExpression(&s.Lhs, state)
	rhs_state := executeExpression(&s.Rhs, state)

	final_state := lhs_state
	if lhs_state.currentValue.GetType() == bytecode.ValueType_String {
		if s.Op == ast.PLUS {
			final := lhs_state.currentValue.GetString() + rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.StringValue{Value: final}
		} else if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.GetString() == rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.GetString() != rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.GetString() < rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.GetString() > rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.GetString() <= rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.GetString() >= rhs_state.currentValue.GetString()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if rhs_state.currentValue.GetType() == bytecode.ValueType_Number && s.Op == ast.MINUS {
			final := lhs_state.currentValue.GetNumber() - rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if rhs_state.currentValue.GetType() == bytecode.ValueType_Number && s.Op == ast.MULT {
			final := lhs_state.currentValue.GetNumber() * rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if rhs_state.currentValue.GetType() == bytecode.ValueType_Number && s.Op == ast.DIV {
			final := lhs_state.currentValue.GetNumber() / rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if rhs_state.currentValue.GetType() == bytecode.ValueType_Number && s.Op == ast.MOD {
			final := lhs_state.currentValue.GetNumber() % rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else {
			panic("SHOULDN'T GET HERE (string) :(")
		}
	} else if lhs_state.currentValue.GetType() == bytecode.ValueType_Boolean {
		if s.Op == ast.AND {
			final := lhs_state.currentValue.GetBoolean() && rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.OR {
			final := lhs_state.currentValue.GetBoolean() || rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.GetBoolean() == rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.GetBoolean() != rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.GetNumber() < rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.GetNumber() > rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.GetNumber() <= rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.GetNumber() >= rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else {
			panic("SHOULDN'T GET HERE (bool) :(")
		}
	} else if lhs_state.currentValue.GetType() == bytecode.ValueType_Number {
		if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.GetBoolean() == rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.GetBoolean() != rhs_state.currentValue.GetBoolean()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.GetNumber() < rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.GetNumber() > rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.GetNumber() <= rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.GetNumber() >= rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.BooleanValue{Value: final}
		} else if s.Op == ast.PLUS {
			final := lhs_state.currentValue.GetNumber() + rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if s.Op == ast.MINUS {
			final := lhs_state.currentValue.GetNumber() - rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if s.Op == ast.MULT {
			final := lhs_state.currentValue.GetNumber() * rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if s.Op == ast.DIV {
			final := lhs_state.currentValue.GetNumber() / rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else if s.Op == ast.MOD {
			final := lhs_state.currentValue.GetNumber() % rhs_state.currentValue.GetNumber()
			final_state.currentValue = bytecode.NumberValue{Value: final}
		} else {
			panic("SHOULDN'T GET HERE (number) :(")
		}
	}
	return final_state
}

func executeUnaryExpression(s *ast.AstProcessUnaryExpression, state ProcessState) ProcessState {
	expr_state := executeExpression(&s.Expr, state)
	if s.Op == ast.NOT {
		expr_state.currentValue = bytecode.BooleanValue{Value: !expr_state.currentValue.GetBoolean()}
	} else if s.Op == ast.HEAD {
		if len(expr_state.currentValue.GetString()) <= 0 {
			expr_state.currentValue = bytecode.StringValue{Value: ""}
		} else {
			expr_state.currentValue = bytecode.StringValue{Value: expr_state.currentValue.GetString()[0:1]}
		}
	} else if s.Op == ast.TAIL {
		if len(expr_state.currentValue.GetString()) <= 1 {
			expr_state.currentValue = bytecode.StringValue{Value: ""}
		} else {
			expr_state.currentValue = bytecode.StringValue{Value: expr_state.currentValue.GetString()[1:]}
		}
	}
	return expr_state
}

func executeString(s *ast.AstProcessString, state ProcessState) ProcessState {
	state.currentValue = bytecode.StringValue{Value: s.Value}
	return state
}

func executeBoolean(s *ast.AstProcessBoolean, state ProcessState) ProcessState {
	state.currentValue = bytecode.BooleanValue{Value: s.Value}
	return state
}

func executeNumber(s *ast.AstProcessNumber, state ProcessState) ProcessState {
	state.currentValue = bytecode.NumberValue{Value: s.Value}
	return state
}

func executeVariable(s *ast.AstProcessVariable, state ProcessState) ProcessState {
	val, prs := state.environment[s.Name]
	if prs {
		state.currentValue = val
	} else {
		state.currentValue = bytecode.StringValue{Value: ""}
	}
	return state
}
