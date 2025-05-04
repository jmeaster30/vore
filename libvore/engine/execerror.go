package engine

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
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
	for idx, value := range ds.Subslice[bytecode.Value](err.processState.stack, 0, 5) {
		message += fmt.Sprintf("     [%d] (%s) %s\n", idx, value.Type(), value.String())
	}

	message += "Environment ----\n"
	for _, entry := range err.processState.environment.Entries() {
		key, value := entry.Values()
		message += fmt.Sprintf("     [%s] (%s) '%s'\n", key, value.Type(), value.String())
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
