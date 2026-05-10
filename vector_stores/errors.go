package vector_stores

import "fmt"

// VectorStoreError represents errors in vector store operations.
type VectorStoreError struct {
	Kind    VectorStoreErrorKind
	Message string
	Cause   error
}

type VectorStoreErrorKind int

const (
	ErrDocumentNotFound VectorStoreErrorKind = iota
	ErrEmbedding
	ErrStorage
	ErrConnection
)

func (e *VectorStoreError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("vector_store: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("vector_store: %s", e.Message)
}

func (e *VectorStoreError) Unwrap() error { return e.Cause }

func NewVectorStoreError(kind VectorStoreErrorKind, msg string, cause error) *VectorStoreError {
	return &VectorStoreError{Kind: kind, Message: msg, Cause: cause}
}
