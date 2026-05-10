package chains

// ChainResult is the output from a chain invocation.
type ChainResult map[string]any

// BaseChain defines the interface for all chains.
type BaseChain interface {
	InputKeys() []string
	OutputKeys() []string
	Invoke(inputs map[string]any) (ChainResult, error)
	ValidateInputs(inputs map[string]any) error
	Name() string
}

// BaseChainImpl provides default validation logic for chains.
type BaseChainImpl struct{}

func (b *BaseChainImpl) ValidateInputs(inputs map[string]any, keys []string) error {
	for _, key := range keys {
		if _, ok := inputs[key]; !ok {
			return NewChainError(ErrMissingInput, "missing input: "+key, nil)
		}
	}
	return nil
}
