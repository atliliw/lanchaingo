package retrieval

import "fmt"

// RetrieverError represents errors in retrieval operations.
type RetrieverError struct {
	Kind    RetrieverErrorKind
	Message string
	Cause   error
}

type RetrieverErrorKind int

const (
	ErrStore  RetrieverErrorKind = iota
	ErrEmbedding
	ErrNoResults
)

func (e *RetrieverError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("retriever: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("retriever: %s", e.Message)
}

func (e *RetrieverError) Unwrap() error { return e.Cause }

func NewRetrieverError(kind RetrieverErrorKind, msg string, cause error) *RetrieverError {
	return &RetrieverError{Kind: kind, Message: msg, Cause: cause}
}
