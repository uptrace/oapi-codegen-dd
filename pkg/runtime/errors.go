package runtime

import "errors"

// ErrValidationEmail is the sentinel error returned when an email fails validation
var (
	ErrValidationEmail = errors.New("email: failed to pass regex validation")
)

// ClientAPIError represents type for client API errors.
type ClientAPIError struct {
	err error
}

// Error implements the error interface.
func (e *ClientAPIError) Error() string {
	if e.err == nil {
		return "client api error"
	}
	return e.err.Error()
}

// Unwrap returns the underlying error.
func (e *ClientAPIError) Unwrap() error {
	return e.err
}

// NewClientAPIError creates a new ClientAPIError from the given error.
func NewClientAPIError(err error) error {
	return &ClientAPIError{err: err}
}
