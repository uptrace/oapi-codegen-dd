package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceErrorImplementsError(t *testing.T) {
	var _ error = ServiceError{}
}

func TestServiceError_Error(t *testing.T) {
	t.Run("returns first message from nested arrays", func(t *testing.T) {
		err := ServiceError{
			Errors: []ErrorDetails{
				{Messages: []string{"first error", "second error"}},
				{Messages: []string{"third error"}},
			},
		}
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("returns unknown error when Errors is empty", func(t *testing.T) {
		err := ServiceError{}
		assert.Equal(t, "unknown error", err.Error())
	})

	t.Run("returns unknown error when Messages is empty", func(t *testing.T) {
		err := ServiceError{
			Errors: []ErrorDetails{
				{Messages: []string{}},
			},
		}
		assert.Equal(t, "unknown error", err.Error())
	})
}
