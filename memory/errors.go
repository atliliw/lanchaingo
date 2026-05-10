package memory

import "fmt"

// MemoryError represents errors in memory operations.
type MemoryError struct {
	Kind    MemoryErrorKind
	Message string
	Cause   error
}

type MemoryErrorKind int

const (
	ErrLoad  MemoryErrorKind = iota
	ErrSave
	ErrClear
	ErrOther
)

func (e *MemoryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("memory: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("memory: %s", e.Message)
}

func (e *MemoryError) Unwrap() error {
	return e.Cause
}

func NewMemoryError(kind MemoryErrorKind, msg string, cause error) *MemoryError {
	return &MemoryError{Kind: kind, Message: msg, Cause: cause}
}
