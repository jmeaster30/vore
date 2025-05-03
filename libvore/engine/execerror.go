package engine

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/bytecode"
)

type ExecError struct {
	processState ProcessState
	instruction  bytecode.ProcInstruction
	message      string
}

func (err ExecError) Error() string {
	message := fmt.Sprintf("ExecError: %s\n", err.message)
	message += fmt.Sprintf("Inst [%d] %T %#v\n", err.processState.instructionPointer, err.instruction, err.instruction)

	message += "Stack ----------\n"
	for idx, value := range err.processState.stack.Subslice(0, 5) {
		message += fmt.Sprintf("     [%d] (%s) %s\n", idx, value.GetType(), value.GetString())
	}

	message += "Environment ----\n"
	for key, value := range err.processState.environment {
		message += fmt.Sprintf("     [%s] (%s) '%s'\n", key, value.GetType(), value.GetString())
	}

	return message
}

func NewExecError(message string, instruction bytecode.ProcInstruction, processState ProcessState) ExecError {
	return ExecError{
		processState,
		instruction,
		message,
	}
}
