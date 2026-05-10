package chains

import "fmt"

// ChainStep represents one step in a SequentialChain.
type ChainStep struct {
	Chain         BaseChain
	InputMapping  map[string]string
	OutputMapping map[string]string
}

// SequentialChain executes multiple chains in sequence.
type SequentialChain struct {
	steps []ChainStep
	name  string
}

func NewSequentialChain() *SequentialChain {
	return &SequentialChain{
		steps: make([]ChainStep, 0),
		name:  "sequential_chain",
	}
}

func (c *SequentialChain) WithName(name string) *SequentialChain {
	c.name = name
	return c
}

func (c *SequentialChain) addChain(chain BaseChain, inputKeys, outputKeys []string) *SequentialChain {
	inputMapping := make(map[string]string)
	for _, k := range inputKeys {
		inputMapping[k] = k
	}
	outputMapping := make(map[string]string)
	for _, k := range outputKeys {
		outputMapping[k] = k
	}
	c.steps = append(c.steps, ChainStep{
		Chain:         chain,
		InputMapping:  inputMapping,
		OutputMapping: outputMapping,
	})
	return c
}

func (c *SequentialChain) AddChain(chain BaseChain, inputKeys, outputKeys []string) *SequentialChain {
	return c.addChain(chain, inputKeys, outputKeys)
}

func (c *SequentialChain) AddChainWithMapping(chain BaseChain, inputMapping, outputMapping map[string]string) *SequentialChain {
	c.steps = append(c.steps, ChainStep{
		Chain:         chain,
		InputMapping:  inputMapping,
		OutputMapping: outputMapping,
	})
	return c
}

func (c *SequentialChain) InputKeys() []string {
	if len(c.steps) == 0 {
		return nil
	}
	keys := make([]string, 0, len(c.steps[0].InputMapping))
	for _, v := range c.steps[0].InputMapping {
		keys = append(keys, v)
	}
	return keys
}

func (c *SequentialChain) OutputKeys() []string {
	if len(c.steps) == 0 {
		return nil
	}
	last := c.steps[len(c.steps)-1]
	keys := make([]string, 0, len(last.OutputMapping))
	for _, v := range last.OutputMapping {
		keys = append(keys, v)
	}
	return keys
}

func (c *SequentialChain) ValidateInputs(inputs map[string]any) error {
	for _, key := range c.InputKeys() {
		if _, ok := inputs[key]; !ok {
			return NewChainError(ErrMissingInput, "missing input: "+key, nil)
		}
	}
	return nil
}

func (c *SequentialChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	state := make(map[string]any)
	for k, v := range inputs {
		state[k] = v
	}

	finalOutput := make(ChainResult)

	for i, step := range c.steps {
		stepInputs := make(map[string]any)
		for chainKey, globalKey := range step.InputMapping {
			if val, ok := state[globalKey]; ok {
				stepInputs[chainKey] = val
			} else {
				return nil, NewChainError(ErrMissingInput,
					fmt.Sprintf("step %d: missing input '%s' (mapped from '%s')", i, chainKey, globalKey), nil)
			}
		}

		output, err := step.Chain.Invoke(stepInputs)
		if err != nil {
			return nil, NewChainError(ErrExecution,
				fmt.Sprintf("step %d (%s) failed", i, step.Chain.Name()), err)
		}

		for chainKey, globalKey := range step.OutputMapping {
			if val, ok := output[chainKey]; ok {
				state[globalKey] = val
				finalOutput[globalKey] = val
			}
		}
	}

	return finalOutput, nil
}

func (c *SequentialChain) Name() string {
	return c.name
}
