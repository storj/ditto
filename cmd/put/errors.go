package put

import (
	"fmt"
)

type InvalidArgsLenError struct {
	argsLen int
}

func NewInvalidArgsError(argsLen int) *InvalidArgsLenError {
	return &InvalidArgsLenError{argsLen}
}

func (e *InvalidArgsLenError) Error() string {
	return fmt.Sprintf("invalid args len, provided %d, should be 2", e.argsLen)
}