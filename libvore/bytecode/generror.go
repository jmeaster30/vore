package bytecode

import "fmt"

type GenError struct {
	err error
}

func (g *GenError) Error() string {
	return g.err.Error()
}

func NewGenError(msg string) *GenError {
	return &GenError{fmt.Errorf("GEN ERROR: %s\n", msg)}
}
