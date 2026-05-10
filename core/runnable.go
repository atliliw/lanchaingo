package core

import "context"

// Runnable is the unified execution interface for all LangChain components.
//
// Every component (LLM, Chain, Agent, Retriever, etc.) implements Runnable,
// enabling uniform composition through Invoke/Batch/Stream.
//
// Maps to Rust langchainrust::core::runnables::runnable_trait::Runnable.
type Runnable interface {
	// Invoke transforms a single input to output.
	Invoke(ctx context.Context, input map[string]any) (any, error)

	// Batch processes multiple inputs and returns their outputs.
	// Default implementation processes inputs sequentially.
	// Override for concurrent execution or batched API calls.
	Batch(ctx context.Context, inputs []map[string]any) ([]any, error)

	// Stream returns a channel that yields output values in real-time.
	// Types with native streaming (e.g., LLM chat) override this.
	Stream(ctx context.Context, input map[string]any) (<-chan any, error)
}

// RunnableFunc is an adapter that turns a plain function into a Runnable.
type RunnableFunc func(ctx context.Context, input map[string]any) (any, error)

func (f RunnableFunc) Invoke(ctx context.Context, input map[string]any) (any, error) {
	return f(ctx, input)
}

func (f RunnableFunc) Batch(ctx context.Context, inputs []map[string]any) ([]any, error) {
	results := make([]any, len(inputs))
	for i, input := range inputs {
		result, err := f(ctx, input)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (f RunnableFunc) Stream(ctx context.Context, input map[string]any) (<-chan any, error) {
	ch := make(chan any, 1)
	go func() {
		defer close(ch)
		result, err := f(ctx, input)
		if err != nil {
			return
		}
		ch <- result
	}()
	return ch, nil
}
