package gen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessPaymentResponse_UnmarshalJSON(t *testing.T) {
	contents := `{"b": "b-value"}`

	var p ProcessPaymentResponse
	err := json.Unmarshal([]byte(contents), &p)
	require.NoError(t, err)

	res, _ := p.ProcessPayment_Response_OneOf.AsResponseB()

	assert.Equal(t, "b-value", *res.B)
}

func TestProcessPaymentResponse_MarshalJSON(t *testing.T) {
	oneOf := &ProcessPayment_Response_OneOf{}
	_ = oneOf.FromResponseB(ResponseB{
		B: func() *string {
			s := "b-value"
			return &s
		}(),
	})
	p := ProcessPaymentResponse{
		ProcessPayment_Response_OneOf: oneOf,
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	assert.JSONEq(t, `{"b": "b-value"}`, string(data))
}
