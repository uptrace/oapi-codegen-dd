package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceErrorImplementsError(t *testing.T) {
	var _ error = ServiceError{}
}

func TestServiceError_Error(t *testing.T) {
	t.Run("returns message from nested path", func(t *testing.T) {
		msg := "something went wrong"
		err := ServiceError{
			ErrorData: &ErrorDetails{
				Message: &msg,
			},
		}
		assert.Equal(t, msg, err.Error())
	})

	t.Run("returns unknown error when ErrorData is nil", func(t *testing.T) {
		err := ServiceError{}
		assert.Equal(t, "unknown error", err.Error())
	})

	t.Run("returns unknown error when Message is nil", func(t *testing.T) {
		err := ServiceError{
			ErrorData: &ErrorDetails{},
		}
		assert.Equal(t, "unknown error", err.Error())
	})
}
