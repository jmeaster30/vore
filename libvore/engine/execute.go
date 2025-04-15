package engine

import (
	"fmt"
	"strconv"

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

type ProcessValue interface {
	getType() bytecode.ProcessType
	getString() string
	getNumber() int
	getBoolean() bool
}

type ProcessValueString struct {
	value string
}

func (v ProcessValueString) getString() string {
	return v.value
}

func (v ProcessValueString) getNumber() int {
	intval, err := strconv.Atoi(v.value)
	if err != nil {
		return 0
	}
	return intval
}

func (v ProcessValueString) getBoolean() bool {
	return len(v.value) != 0
}

func (v ProcessValueString) getType() bytecode.ProcessType {
	return bytecode.PTSTRING
}

type ProcessValueNumber struct {
	value int
}

func (v ProcessValueNumber) getString() string {
	return strconv.Itoa(v.value)
}

func (v ProcessValueNumber) getNumber() int {
	return v.value
}

func (v ProcessValueNumber) getBoolean() bool {
	return v.value != 0
}

func (v ProcessValueNumber) getType() bytecode.ProcessType {
	return bytecode.PTNUMBER
}

type ProcessValueBoolean struct {
	value bool
}

func (v ProcessValueBoolean) getString() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v ProcessValueBoolean) getNumber() int {
	if v.value {
		return 1
	}
	return 0
}

func (v ProcessValueBoolean) getBoolean() bool {
	return v.value
}

func (v ProcessValueBoolean) getType() bytecode.ProcessType {
	return bytecode.PTBOOLEAN
}

type ProcessState struct {
	currentValue ProcessValue
	environment  map[string]ProcessValue
	status       ProcessStatus
}

func executeStatement(s ast.AstProcessStatement, state ProcessState) ProcessState {
	var si any = s
	switch si.(type) {
	case ast.AstProcessSet:
		return executeSet(si.(ast.AstProcessSet), state)
	case ast.AstProcessIf:
		return executeIf(si.(ast.AstProcessIf), state)
	case ast.AstProcessLoop:
		return executeLoop(si.(ast.AstProcessLoop), state)
	case ast.AstProcessBreak:
		return executeBreak(si.(ast.AstProcessBreak), state)
	case ast.AstProcessContinue:
		return executeContinue(si.(ast.AstProcessContinue), state)
	case ast.AstProcessReturn:
		return executeReturn(si.(ast.AstProcessReturn), state)
	case ast.AstProcessDebug:
		return executeDebug(si.(ast.AstProcessDebug), state)
	}
	panic(fmt.Sprintf("Unknown process statement %T", si))
}

func executeSet(s ast.AstProcessSet, state ProcessState) ProcessState {
	expr_state := executeExpression(s.Expr, state)
	expr_state.environment[s.Name] = expr_state.currentValue
	return expr_state
}

func executeReturn(s ast.AstProcessReturn, state ProcessState) ProcessState {
	expr_state := executeExpression(s.Expr, state)
	expr_state.status = RETURNING
	return expr_state
}

func executeIf(s ast.AstProcessIf, state ProcessState) ProcessState {
	expr_state := executeExpression(s.Condition, state)
	if expr_state.currentValue.getBoolean() {
		for _, stmt := range s.TrueBody {
			expr_state = executeStatement(stmt, expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	} else {
		for _, stmt := range s.FalseBody {
			expr_state = executeStatement(stmt, expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	}

	return expr_state
}

func executeDebug(s ast.AstProcessDebug, state ProcessState) ProcessState {
	expr_state := executeExpression(s.Expr, state)
	fmt.Println(expr_state.currentValue.getString())
	return expr_state
}

func executeLoop(s ast.AstProcessLoop, state ProcessState) ProcessState {
	expr_state := state
	for {
		for _, stmt := range s.Body {
			expr_state = executeStatement(stmt, expr_state)
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

func executeContinue(s ast.AstProcessContinue, state ProcessState) ProcessState {
	state.status = CONTINUELOOP
	return state
}

func executeBreak(s ast.AstProcessBreak, state ProcessState) ProcessState {
	state.status = BREAKLOOP
	return state
}

func executeExpression(s ast.AstProcessExpression, state ProcessState) ProcessState {
	var si any = s
	switch si.(type) {
	case ast.AstProcessBinaryExpression:
		return executeBinaryExpr(si.(ast.AstProcessBinaryExpression), state)
	case ast.AstProcessUnaryExpression:
		return executeUnaryExpression(si.(ast.AstProcessUnaryExpression), state)
	case ast.AstProcessString:
		return executeString(si.(ast.AstProcessString), state)
	case ast.AstProcessBoolean:
		return executeBoolean(si.(ast.AstProcessBoolean), state)
	case ast.AstProcessVariable:
		return executeVariable(si.(ast.AstProcessVariable), state)
	}
	panic(fmt.Sprintf("unknown process expression type %T", si))
}

func executeBinaryExpr(s ast.AstProcessBinaryExpression, state ProcessState) ProcessState {
	lhs_state := executeExpression(s.Lhs, state)
	rhs_state := executeExpression(s.Rhs, state)

	final_state := lhs_state
	if lhs_state.currentValue.getType() == bytecode.PTSTRING {
		if s.Op == ast.PLUS {
			final := lhs_state.currentValue.getString() + rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueString{final}
		} else if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.getString() == rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.getString() != rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.getString() < rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.getString() > rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.getString() <= rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.getString() >= rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if rhs_state.currentValue.getType() == bytecode.PTNUMBER && s.Op == ast.MINUS {
			final := lhs_state.currentValue.getNumber() - rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == bytecode.PTNUMBER && s.Op == ast.MULT {
			final := lhs_state.currentValue.getNumber() * rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == bytecode.PTNUMBER && s.Op == ast.DIV {
			final := lhs_state.currentValue.getNumber() / rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == bytecode.PTNUMBER && s.Op == ast.MOD {
			final := lhs_state.currentValue.getNumber() % rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else {
			panic("SHOULDN'T GET HERE (string) :(")
		}
	} else if lhs_state.currentValue.getType() == bytecode.PTBOOLEAN {
		if s.Op == ast.AND {
			final := lhs_state.currentValue.getBoolean() && rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.OR {
			final := lhs_state.currentValue.getBoolean() || rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.getBoolean() == rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.getBoolean() != rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.getNumber() < rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.getNumber() > rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.getNumber() <= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.getNumber() >= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else {
			panic("SHOULDN'T GET HERE (bool) :(")
		}
	} else if lhs_state.currentValue.getType() == bytecode.PTNUMBER {
		if s.Op == ast.DEQUAL {
			final := lhs_state.currentValue.getBoolean() == rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.NEQUAL {
			final := lhs_state.currentValue.getBoolean() != rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESS {
			final := lhs_state.currentValue.getNumber() < rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATER {
			final := lhs_state.currentValue.getNumber() > rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.LESSEQ {
			final := lhs_state.currentValue.getNumber() <= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.GREATEREQ {
			final := lhs_state.currentValue.getNumber() >= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.Op == ast.PLUS {
			final := lhs_state.currentValue.getNumber() + rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.Op == ast.MINUS {
			final := lhs_state.currentValue.getNumber() - rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.Op == ast.MULT {
			final := lhs_state.currentValue.getNumber() * rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.Op == ast.DIV {
			final := lhs_state.currentValue.getNumber() / rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.Op == ast.MOD {
			final := lhs_state.currentValue.getNumber() % rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else {
			panic("SHOULDN'T GET HERE (number) :(")
		}
	}
	return final_state
}

func executeUnaryExpression(s ast.AstProcessUnaryExpression, state ProcessState) ProcessState {
	expr_state := executeExpression(s.Expr, state)
	if s.Op == ast.NOT {
		expr_state.currentValue = ProcessValueBoolean{!expr_state.currentValue.getBoolean()}
	} else if s.Op == ast.HEAD {
		if len(expr_state.currentValue.getString()) <= 0 {
			expr_state.currentValue = ProcessValueString{""}
		} else {
			expr_state.currentValue = ProcessValueString{expr_state.currentValue.getString()[0:1]}
		}
	} else if s.Op == ast.TAIL {
		if len(expr_state.currentValue.getString()) <= 1 {
			expr_state.currentValue = ProcessValueString{""}
		} else {
			expr_state.currentValue = ProcessValueString{expr_state.currentValue.getString()[1:]}
		}
	}
	return expr_state
}

func executeString(s ast.AstProcessString, state ProcessState) ProcessState {
	state.currentValue = ProcessValueString{s.Value}
	return state
}

func executeBoolean(s ast.AstProcessBoolean, state ProcessState) ProcessState {
	state.currentValue = ProcessValueBoolean{s.Value}
	return state
}

func executeNumber(s ast.AstProcessNumber, state ProcessState) ProcessState {
	state.currentValue = ProcessValueNumber{s.Value}
	return state
}

func executeVariable(s ast.AstProcessVariable, state ProcessState) ProcessState {
	val, prs := state.environment[s.Name]
	if prs {
		state.currentValue = val
	} else {
		state.currentValue = ProcessValueString{""}
	}
	return state
}
