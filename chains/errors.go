package chains

import "fmt"

// ChainError represents errors in chain operations.
type ChainError struct {
	Kind    ChainErrorKind
	Message string
	Cause   error
}

type ChainErrorKind int

const (
	ErrMissingInput   ChainErrorKind = iota
	ErrOutput
	ErrExecution
	ErrOther
)

func (e *ChainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("chain: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("chain: %s", e.Message)
}

func (e *ChainError) Unwrap() error {
	return e.Cause
}

func NewChainError(kind ChainErrorKind, msg string, cause error) *ChainError {
	return &ChainError{Kind: kind, Message: msg, Cause: cause}
}
