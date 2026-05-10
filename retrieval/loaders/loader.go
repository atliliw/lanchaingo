package loaders

import "fmt"

// LoaderError represents errors in document loading.
type LoaderError struct {
	Message string
	Cause   error
}

func (e *LoaderError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("loader: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("loader: %s", e.Message)
}

func (e *LoaderError) Unwrap() error { return e.Cause }

func NewLoaderError(msg string, cause error) *LoaderError {
	return &LoaderError{Message: msg, Cause: cause}
}
