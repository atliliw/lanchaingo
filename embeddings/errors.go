package embeddings

import "fmt"

// EmbeddingError represents errors in embedding operations.
type EmbeddingError struct {
	Kind    EmbeddingErrorKind
	Message string
	Cause   error
}

type EmbeddingErrorKind int

const (
	ErrHTTP    EmbeddingErrorKind = iota
	ErrAPI
	ErrParse
	ErrEmptyInput
)

func (e *EmbeddingError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("embedding: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("embedding: %s", e.Message)
}

func (e *EmbeddingError) Unwrap() error { return e.Cause }

func NewEmbeddingError(kind EmbeddingErrorKind, msg string, cause error) *EmbeddingError {
	return &EmbeddingError{Kind: kind, Message: msg, Cause: cause}
}
