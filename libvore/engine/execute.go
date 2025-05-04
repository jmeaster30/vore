package engine

import (
	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
)

type ProcessState struct {
	instructionPointer int
	environment        bytecode.MapValue
	stack              *ds.Stack[bytecode.Value]
}

func executeProcessInstructions(insts []bytecode.ProcInstruction, environment bytecode.MapValue) (ds.Optional[bytecode.Value], error) {
	currentState := ProcessState{
		instructionPointer: 0,
		environment:        environment,
		stack:              ds.NewStack[bytecode.Value](),
	}

	for currentState.instructionPointer < len(insts) {
		currentInstruction := insts[currentState.instructionPointer]
		newState, err := executeProcessInstruction(&currentInstruction, currentState)
		if err != nil {
			return ds.None[bytecode.Value](), err
		}

		oldIP := currentState.instructionPointer

		currentState = newState
		if oldIP == newState.instructionPointer {
			currentState.instructionPointer += 1
		}
	}

	return currentState.stack.Pop(), nil
}

func executeProcessInstruction(i *bytecode.ProcInstruction, state ProcessState) (ProcessState, error) {
	var ii any = *i
	switch inst := ii.(type) {
	case *bytecode.Jump:
		return executeJump(inst, state)
	case *bytecode.Store:
		return executeStore(inst, state)
	case *bytecode.Load:
		return executeLoad(inst, state)
	case *bytecode.Push:
		return executePush(inst, state)
	case *bytecode.ConditionalJump:
		return executeConditionalJump(inst, state)
	case *bytecode.Debug:
		return executeDebug(inst, state)
	case *bytecode.Return:
		return executeReturn(inst, state)
	case *bytecode.Not:
		return executeNot(state)
	case *bytecode.Head:
		return executeHead(state)
	case *bytecode.Tail:
		return executeTail(state)
	case *bytecode.And:
		return executeAnd(state)
	case *bytecode.Or:
		return executeOr(state)
	case *bytecode.Add:
		return executeAdd(state)
	case *bytecode.Subtract:
		return executeSubtract(state)
	case *bytecode.Multiply:
		return executeMultiply(state)
	case *bytecode.Divide:
		return executeDivide(state)
	case *bytecode.Modulo:
		return executeModulo(state)
	case *bytecode.Equal:
		return executeEqual(state)
	case *bytecode.NotEqual:
		return executeNotEqual(state)
	case *bytecode.LessThan:
		return executeLessThan(state)
	case *bytecode.LessThanEqual:
		return executeLessThanEqual(state)
	case *bytecode.GreaterThan:
		return executeGreaterThan(state)
	case *bytecode.GreaterThanEqual:
		return executeGreaterThanEqual(state)
	}
	return state, NewExecError("Unknown process instruction", *i, state)
}

func executeJump(inst *bytecode.Jump, state ProcessState) (ProcessState, error) {
	state.instructionPointer = inst.NewProgramCounter
	return state, nil
}

func executeStore(inst *bytecode.Store, state ProcessState) (ProcessState, error) {
	if state.stack.IsEmpty() {
		return state, NewExecError("Empty stack for store instruction", *inst, state)
	}

	storedValue := state.stack.Pop().GetValue()
	state.environment.Set(inst.VariableName, storedValue)
	return state, nil
}

func executeLoad(inst *bytecode.Load, state ProcessState) (ProcessState, error) {
	value, ok := state.environment.Get(inst.VariableName)
	if ok {
		state.stack.Push(bytecode.NewString(""))
	} else {
		state.stack.Push(value)
	}
	return state, nil
}

func executePush(inst *bytecode.Push, state ProcessState) (ProcessState, error) {
	state.stack.Push(inst.Value)
	return state, nil
}

func executeConditionalJump(inst *bytecode.ConditionalJump, state ProcessState) (ProcessState, error) {
	if state.stack.IsEmpty() {
		return state, NewExecError("Empty stack for conditional jump", *inst, state)
	}

	condition := state.stack.Pop().GetValue()
	if condition.Boolean() {
		state.instructionPointer = inst.NewProgramCounter
	}
	return state, nil
}

func executeDebug(inst *bytecode.Debug, state ProcessState) (ProcessState, error) {
	if state.stack.IsEmpty() {
		return state, NewExecError("Empty stack for debug print", *inst, state)
	}

	value := state.stack.Pop().GetValue()
	print(value.String())
	return state, nil
}

func executeReturn(inst *bytecode.Return, state ProcessState) (ProcessState, error) {
	return state, NewExecError("TODO I didn't implement returns yet lol", *inst, state)
}

func executeNot(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeHead(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeTail(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeAnd(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeOr(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeAdd(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeSubtract(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeMultiply(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeDivide(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeModulo(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeEqual(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeNotEqual(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeLessThan(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeLessThanEqual(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeGreaterThan(state ProcessState) (ProcessState, error) {
	return state, nil
}

func executeGreaterThanEqual(state ProcessState) (ProcessState, error) {
	return state, nil
}
