package libvore

import "github.com/jmeaster30/vore/libvore/ast"

type VoreFileError struct {
	err error
}

func (err *VoreFileError) Error() string {
	return err.err.Error()
}

func (err *VoreFileError) Token() *ast.Token {
	return nil
}

func (err *VoreFileError) Message() string {
	return err.Error()
}

func NewVoreFileError(err error) *VoreFileError {
	return &VoreFileError{err}
}
