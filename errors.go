package mobilepay

import (
	"errors"
	"fmt"
)

var (
	ErrMissingVerifierProperties = errors.New("missing required verifier properties signature or webhook url")
)

// ArgError is an error that represents an error with an input to mobilepay app payment. It
// identifies the argument and the cause (if possible).
type ArgError struct {
	arg    string
	reason string
}

var _ error = &ArgError{}

// newArgError creates an InputError.
func newArgError(arg, reason string) *ArgError {
	return &ArgError{
		arg:    arg,
		reason: reason,
	}
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("%s is invalid because %s", e.arg, e.reason)
}
