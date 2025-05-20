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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (err ExecError) Error() string {
	message := fmt.Sprintf("ExecError: %s\n", err.message)
	message += "Instructions ---\n"
	startIdx := err.processState.instructionPointer
	endIdx := max(startIdx-5, 0)
	for idx, inst := range err.processState.instructions[endIdx : startIdx+1] {
		if idx+endIdx == err.processState.instructionPointer {
			message += fmt.Sprintf("   @@[%d] %v\n", idx+endIdx, inst)
		} else {
			message += fmt.Sprintf("     [%d] %v\n", idx+endIdx, inst)
		}
	}

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

func (err ExecError) Message() string {
	return err.message
}

func NewExecError(message string, instruction bytecode.ProcInstruction, processState ProcessState) ExecError {
	return ExecError{
		processState,
		instruction,
		message,
	}
}
