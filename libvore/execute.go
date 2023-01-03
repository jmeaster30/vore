package libvore

import (
	"fmt"
	"strconv"
)

type ProcessStatus int

const (
	NEXT ProcessStatus = iota
	BREAKLOOP
	CONTINUELOOP
	RETURNING
)

type ProcessValue interface {
	getType() ProcessType
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
func (v ProcessValueString) getType() ProcessType {
	return PTSTRING
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
func (v ProcessValueNumber) getType() ProcessType {
	return PTNUMBER
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
func (v ProcessValueBoolean) getType() ProcessType {
	return PTBOOLEAN
}

type ProcessState struct {
	currentValue ProcessValue
	environment  map[string]ProcessValue
	status       ProcessStatus
}

func (s AstProcessSet) execute(state ProcessState) ProcessState {
	expr_state := s.expr.execute(state)
	expr_state.environment[s.name] = expr_state.currentValue
	return expr_state
}

func (s AstProcessReturn) execute(state ProcessState) ProcessState {
	expr_state := s.expr.execute(state)
	expr_state.status = RETURNING
	return expr_state
}

func (s AstProcessIf) execute(state ProcessState) ProcessState {
	expr_state := s.condition.execute(state)
	if expr_state.currentValue.getBoolean() {
		for _, stmt := range s.trueBody {
			expr_state = stmt.execute(expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	} else {
		for _, stmt := range s.falseBody {
			expr_state = stmt.execute(expr_state)
			if expr_state.status != NEXT {
				break
			}
		}
	}

	return expr_state
}

func (s AstProcessDebug) execute(state ProcessState) ProcessState {
	expr_state := s.expr.execute(state)
	fmt.Println(expr_state.currentValue.getString())
	return expr_state
}

func (s AstProcessLoop) execute(state ProcessState) ProcessState {
	expr_state := state
	for {
		for _, stmt := range s.body {
			expr_state = stmt.execute(expr_state)
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

func (s AstProcessContinue) execute(state ProcessState) ProcessState {
	state.status = CONTINUELOOP
	return state
}

func (s AstProcessBreak) execute(state ProcessState) ProcessState {
	state.status = BREAKLOOP
	return state
}

func (s AstProcessBinaryExpression) execute(state ProcessState) ProcessState {
	lhs_state := s.lhs.execute(state)
	rhs_state := s.rhs.execute(state)

	final_state := lhs_state
	if lhs_state.currentValue.getType() == PTSTRING {
		if s.op == PLUS {
			final := lhs_state.currentValue.getString() + rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueString{final}
		} else if s.op == DEQUAL {
			final := lhs_state.currentValue.getString() == rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == NEQUAL {
			final := lhs_state.currentValue.getString() != rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESS {
			final := lhs_state.currentValue.getString() < rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATER {
			final := lhs_state.currentValue.getString() > rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESSEQ {
			final := lhs_state.currentValue.getString() <= rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATEREQ {
			final := lhs_state.currentValue.getString() >= rhs_state.currentValue.getString()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if rhs_state.currentValue.getType() == PTNUMBER && s.op == MINUS {
			final := lhs_state.currentValue.getNumber() - rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == PTNUMBER && s.op == MULT {
			final := lhs_state.currentValue.getNumber() * rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == PTNUMBER && s.op == DIV {
			final := lhs_state.currentValue.getNumber() / rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if rhs_state.currentValue.getType() == PTNUMBER && s.op == MOD {
			final := lhs_state.currentValue.getNumber() % rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else {
			panic("SHOULDN'T GET HERE (string) :(")
		}
	} else if lhs_state.currentValue.getType() == PTBOOLEAN {
		if s.op == AND {
			final := lhs_state.currentValue.getBoolean() && rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == OR {
			final := lhs_state.currentValue.getBoolean() || rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == DEQUAL {
			final := lhs_state.currentValue.getBoolean() == rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == NEQUAL {
			final := lhs_state.currentValue.getBoolean() != rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESS {
			final := lhs_state.currentValue.getNumber() < rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATER {
			final := lhs_state.currentValue.getNumber() > rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESSEQ {
			final := lhs_state.currentValue.getNumber() <= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATEREQ {
			final := lhs_state.currentValue.getNumber() >= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else {
			panic("SHOULDN'T GET HERE (bool) :(")
		}
	} else if lhs_state.currentValue.getType() == PTNUMBER {
		if s.op == DEQUAL {
			final := lhs_state.currentValue.getBoolean() == rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == NEQUAL {
			final := lhs_state.currentValue.getBoolean() != rhs_state.currentValue.getBoolean()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESS {
			final := lhs_state.currentValue.getNumber() < rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATER {
			final := lhs_state.currentValue.getNumber() > rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == LESSEQ {
			final := lhs_state.currentValue.getNumber() <= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == GREATEREQ {
			final := lhs_state.currentValue.getNumber() >= rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueBoolean{final}
		} else if s.op == PLUS {
			final := lhs_state.currentValue.getNumber() + rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.op == MINUS {
			final := lhs_state.currentValue.getNumber() - rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.op == MULT {
			final := lhs_state.currentValue.getNumber() * rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.op == DIV {
			final := lhs_state.currentValue.getNumber() / rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else if s.op == MOD {
			final := lhs_state.currentValue.getNumber() % rhs_state.currentValue.getNumber()
			final_state.currentValue = ProcessValueNumber{final}
		} else {
			panic("SHOULDN'T GET HERE (number) :(")
		}
	}
	return final_state
}

func (s AstProcessUnaryExpression) execute(state ProcessState) ProcessState {
	expr_state := s.expr.execute(state)
	if s.op == NOT {
		expr_state.currentValue = ProcessValueBoolean{!expr_state.currentValue.getBoolean()}
	} else if s.op == HEAD {
		if len(expr_state.currentValue.getString()) <= 0 {
			expr_state.currentValue = ProcessValueString{""}
		} else {
			expr_state.currentValue = ProcessValueString{expr_state.currentValue.getString()[0:1]}
		}
	} else if s.op == TAIL {
		if len(expr_state.currentValue.getString()) <= 1 {
			expr_state.currentValue = ProcessValueString{""}
		} else {
			expr_state.currentValue = ProcessValueString{expr_state.currentValue.getString()[1:]}
		}
	}
	return expr_state
}

func (s AstProcessString) execute(state ProcessState) ProcessState {
	state.currentValue = ProcessValueString{s.value}
	return state
}

func (s AstProcessNumber) execute(state ProcessState) ProcessState {
	state.currentValue = ProcessValueNumber{s.value}
	return state
}

func (s AstProcessVariable) execute(state ProcessState) ProcessState {
	val, prs := state.environment[s.name]
	if prs {
		state.currentValue = val
	} else {
		state.currentValue = ProcessValueString{""}
	}
	return state
}
