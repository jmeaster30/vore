package engine

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
)

type ProcessState struct {
	instructions       []bytecode.ProcInstruction
	instructionPointer int
	shouldReturn       bool
	environment        bytecode.MapValue
	stack              *ds.Stack[bytecode.Value]
	labels             map[string]int
}

func scanForLabels(insts []bytecode.ProcInstruction) map[string]int {
	result := make(map[string]int)
	for idx, inst := range insts {
		var in any = inst
		switch i := in.(type) {
		case bytecode.Label:
			result[i.Name] = idx
		}
	}
	return result
}

func executeProcessInstructions(insts []bytecode.ProcInstruction, environment bytecode.MapValue) (ds.Optional[bytecode.Value], error) {
	currentState := &ProcessState{
		instructions:       insts,
		instructionPointer: 0,
		shouldReturn:       false,
		environment:        environment,
		stack:              ds.NewStack[bytecode.Value](),
		labels:             scanForLabels(insts),
	}

	for currentState.instructionPointer < len(insts) {
		currentInstructionPointer := currentState.instructionPointer
		currentInstruction := insts[currentState.instructionPointer]
		err := executeProcessInstruction(&currentInstruction, currentState)
		if err != nil {
			return ds.None[bytecode.Value](), err
		}

		if currentInstructionPointer == currentState.instructionPointer {
			currentState.instructionPointer += 1
		}

		if currentState.shouldReturn {
			break
		}
	}

	return currentState.stack.Pop(), nil
}

func executeProcessInstruction(i *bytecode.ProcInstruction, state *ProcessState) error {
	var ii any = *i
	switch inst := ii.(type) {
	case bytecode.Jump:
		return executeJump(inst, state)
	case bytecode.Store:
		return executeStore(inst, state)
	case bytecode.Load:
		return executeLoad(inst, state)
	case bytecode.Push:
		return executePush(inst, state)
	case bytecode.ConditionalJump:
		return executeConditionalJump(inst, state)
	case bytecode.Debug:
		return executeDebug(inst, state)
	case bytecode.Return:
		return executeReturn(inst, state)
	case bytecode.Not:
		return executeNot(inst, state)
	case bytecode.Head:
		return executeHead(inst, state)
	case bytecode.Tail:
		return executeTail(inst, state)
	case bytecode.And:
		return executeAnd(inst, state)
	case bytecode.Or:
		return executeOr(inst, state)
	case bytecode.Add:
		return executeAdd(inst, state)
	case bytecode.Subtract:
		return executeSubtract(inst, state)
	case bytecode.Multiply:
		return executeMultiply(inst, state)
	case bytecode.Divide:
		return executeDivide(inst, state)
	case bytecode.Modulo:
		return executeModulo(inst, state)
	case bytecode.Equal:
		return executeEqual(inst, state)
	case bytecode.NotEqual:
		return executeNotEqual(inst, state)
	case bytecode.LessThan:
		return executeLessThan(inst, state)
	case bytecode.LessThanEqual:
		return executeLessThanEqual(inst, state)
	case bytecode.GreaterThan:
		return executeGreaterThan(inst, state)
	case bytecode.GreaterThanEqual:
		return executeGreaterThanEqual(inst, state)
	case bytecode.LabelJump:
		return executeLabelJump(inst, state)
	case bytecode.Label:
		return nil
	}
	return NewExecError("Unknown process instruction", *i, *state)
}

func executeJump(inst bytecode.Jump, state *ProcessState) error {
	state.instructionPointer = inst.NewProgramCounter
	return nil
}

func executeStore(inst bytecode.Store, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for store instruction", inst, *state)
	}

	storedValue := state.stack.Pop().GetValue()
	state.environment.Set(inst.VariableName, storedValue)
	return nil
}

func executeLoad(inst bytecode.Load, state *ProcessState) error {
	value, ok := state.environment.Get(inst.VariableName)
	if ok {
		state.stack.Push(value)
	} else {
		state.stack.Push(bytecode.NewString(""))
	}
	return nil
}

func executePush(inst bytecode.Push, state *ProcessState) error {
	state.stack.Push(inst.Value)
	return nil
}

func executeConditionalJump(inst bytecode.ConditionalJump, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for conditional jump", inst, *state)
	}

	condition := state.stack.Pop().GetValue()
	if !condition.Boolean() {
		state.instructionPointer = inst.NewProgramCounter
	}
	return nil
}

func executeDebug(inst bytecode.Debug, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for debug print", inst, *state)
	}

	value := state.stack.Pop().GetValue()
	fmt.Printf("%v\n", value)
	return nil
}

func executeReturn(inst bytecode.Return, state *ProcessState) error {
	state.shouldReturn = true
	return nil
}

func executeNot(inst bytecode.Not, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for not operation", inst, *state)
	}

	value := state.stack.Pop().GetValue()
	flipped := !value.Boolean()
	state.stack.Push(bytecode.NewBoolean(flipped))
	return nil
}

func executeHead(inst bytecode.Head, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for head operation", inst, *state)
	}

	value := state.stack.Pop().GetValue()
	str := value.String()

	if len(str) < 1 {
		state.stack.Push(bytecode.NewString(""))
	} else {
		state.stack.Push(bytecode.NewString(string(str[0])))
	}

	return nil
}

func executeTail(inst bytecode.Tail, state *ProcessState) error {
	if state.stack.IsEmpty() {
		return NewExecError("Empty stack for tail operation", inst, *state)
	}

	strvalue := state.stack.Pop().GetValue().String()
	if len(strvalue) < 1 {
		state.stack.Push(bytecode.NewString(""))
	} else {
		state.stack.Push(bytecode.NewString(string(strvalue[1:])))
	}

	return nil
}

func executeAnd(inst bytecode.And, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for and operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Boolean()
	a := state.stack.Pop().GetValue().Boolean()
	state.stack.Push(bytecode.NewBoolean(a && b))
	return nil
}

func executeOr(inst bytecode.Or, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for or operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Boolean()
	a := state.stack.Pop().GetValue().Boolean()
	state.stack.Push(bytecode.NewBoolean(a || b))
	return nil
}

func executeAdd(inst bytecode.Add, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for add operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	switch a.Type() {
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewString(a.String() + b.String()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewNumber(a.Number() + b.Number()))
	default:
		return NewExecError(fmt.Sprintf("Unknown operation + for type (%s, %s) :(", a.Type(), b.Type()), inst, *state)
	}

	return nil
}

func executeSubtract(inst bytecode.Subtract, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for subtract operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Number()
	a := state.stack.Pop().GetValue().Number()
	state.stack.Push(bytecode.NewNumber(a - b))
	return nil
}

func executeMultiply(inst bytecode.Multiply, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for multiply operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Number()
	a := state.stack.Pop().GetValue().Number()
	state.stack.Push(bytecode.NewNumber(a * b))
	return nil
}

func executeDivide(inst bytecode.Divide, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for divide operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Number()
	a := state.stack.Pop().GetValue().Number()
	state.stack.Push(bytecode.NewNumber(a / b))
	return nil
}

func executeModulo(inst bytecode.Modulo, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for add operation", inst, *state)
	}

	b := state.stack.Pop().GetValue().Number()
	a := state.stack.Pop().GetValue().Number()
	state.stack.Push(bytecode.NewNumber(a % b))
	return nil
}

func executeEqual(inst bytecode.Equal, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for equal operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(false))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean(a.Boolean() == b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() == b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(a.String() == b.String()))
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the equality check
	}
	return nil
}

func executeNotEqual(inst bytecode.NotEqual, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for not equal operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(true))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean(a.Boolean() != b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() != b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(a.String() != b.String()))
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(true)) // FIXME actually do the equality check
	}
	return nil
}

func executeLessThan(inst bytecode.LessThan, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for less than operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(false))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean(!a.Boolean() && b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() < b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the comparison
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the equality check
	}
	return nil
}

func executeLessThanEqual(inst bytecode.LessThanEqual, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for less than equal operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(false))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean((!a.Boolean() && b.Boolean()) || a.Boolean() == b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() <= b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the comparison
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the equality check
	}
	return nil
}

func executeGreaterThan(inst bytecode.GreaterThan, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for greater than operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(false))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean(a.Boolean() && !b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() > b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the comparison
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the equality check
	}
	return nil
}

func executeGreaterThanEqual(inst bytecode.GreaterThanEqual, state *ProcessState) error {
	if state.stack.Size() < 2 {
		return NewExecError("Empty stack for greater than equal operation", inst, *state)
	}

	b := state.stack.Pop().GetValue()
	a := state.stack.Pop().GetValue()

	if a.Type() != b.Type() {
		state.stack.Push(bytecode.NewBoolean(false))
		return nil
	}

	switch a.Type() {
	case bytecode.ValueType_Boolean:
		state.stack.Push(bytecode.NewBoolean((!a.Boolean() && b.Boolean()) || a.Boolean() == b.Boolean()))
	case bytecode.ValueType_Number:
		state.stack.Push(bytecode.NewBoolean(a.Number() >= b.Number()))
	case bytecode.ValueType_String:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the comparison
	case bytecode.ValueType_Map:
		state.stack.Push(bytecode.NewBoolean(false)) // FIXME actually do the equality check
	}
	return nil
}

func executeLabelJump(inst bytecode.LabelJump, state *ProcessState) error {
	newIP, ok := state.labels[inst.Label]
	if ok {
		state.instructionPointer = newIP
		return nil
	}
	return NewExecError("Unknown jump label", inst, *state)
}
