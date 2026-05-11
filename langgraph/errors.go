package langgraph

import "fmt"

// GraphError represents errors in graph operations.
type GraphError struct {
	Kind    GraphErrorKind
	Message string
	Cause   error
}

type GraphErrorKind int

const (
	ErrValidation      GraphErrorKind = iota
	ErrExecution
	ErrRouting
	ErrRecursionLimit
	ErrNode
	ErrCheckpoint
	ErrState
	ErrInterrupted
	ErrResume
	ErrInfiniteCycle
	ErrOrphanNode
	ErrDuplicateEdge
	ErrSubgraphCycle
)

func (e *GraphError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("graph: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("graph: %s", e.Message)
}

func (e *GraphError) Unwrap() error { return e.Cause }

func NewGraphError(kind GraphErrorKind, msg string, cause error) *GraphError {
	return &GraphError{Kind: kind, Message: msg, Cause: cause}
}

// SubgraphInterruptError is returned by SubgraphNode.Execute when the embedded
// subgraph was interrupted during execution. The parent graph's execution loop
// recognizes this error and converts it into a GraphInvocationResult with
// Interrupted=true, enabling interrupt propagation across graph boundaries.
type SubgraphInterruptError struct {
	SubgraphName  string
	InterruptedAt string
	PartialState  StateSchema
}

func (e *SubgraphInterruptError) Error() string {
	return fmt.Sprintf("graph: subgraph '%s' interrupted at '%s'", e.SubgraphName, e.InterruptedAt)
}
