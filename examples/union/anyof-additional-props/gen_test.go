package union

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotification_UnmarshalJSON_Email(t *testing.T) {
	input := `{"email": "test@example.com", "subject": "Hello", "tracking": "abc123"}`

	var n Notification
	err := json.Unmarshal([]byte(input), &n)
	require.NoError(t, err)

	// Union variant
	require.NotNil(t, n.Notification_AnyOf)
	email, err := n.Notification_AnyOf.AsEmailNotification()
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Email)
	assert.Equal(t, "Hello", *email.Subject)

	// AdditionalProperties - only tracking
	assert.Equal(t, map[string]string{"tracking": "abc123"}, n.AdditionalProperties)
}

func TestNotification_UnmarshalJSON_SMS(t *testing.T) {
	input := `{"phone": "+1234567890", "campaign": "promo"}`

	var n Notification
	err := json.Unmarshal([]byte(input), &n)
	require.NoError(t, err)

	sms, err := n.Notification_AnyOf.AsSMSNotification()
	require.NoError(t, err)
	assert.Equal(t, "+1234567890", sms.Phone)
	assert.Equal(t, map[string]string{"campaign": "promo"}, n.AdditionalProperties)
}

func TestNotification_UnmarshalJSON_Push(t *testing.T) {
	input := `{"device_token": "token123", "badge": 5, "source": "mobile"}`

	var n Notification
	err := json.Unmarshal([]byte(input), &n)
	require.NoError(t, err)

	push, err := n.Notification_AnyOf.AsPushNotification()
	require.NoError(t, err)
	assert.Equal(t, "token123", push.DeviceToken)
	assert.Equal(t, 5, *push.Badge)
	assert.Equal(t, map[string]string{"source": "mobile"}, n.AdditionalProperties)
}

func TestNotification_UnmarshalJSON_NoAdditionalProperties(t *testing.T) {
	input := `{"email": "test@example.com"}`

	var n Notification
	err := json.Unmarshal([]byte(input), &n)
	require.NoError(t, err)

	assert.Nil(t, n.AdditionalProperties)
}

func TestNotification_MarshalJSON(t *testing.T) {
	n := Notification{
		Notification_AnyOf:   &Notification_AnyOf{},
		AdditionalProperties: map[string]string{"tracking": "xyz"},
	}
	require.NoError(t, n.Notification_AnyOf.FromSMSNotification(SMSNotification{Phone: "+1234567890"}))

	data, err := json.Marshal(n)
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "+1234567890", result["phone"])
	assert.Equal(t, "xyz", result["tracking"])
}

func TestNotification_Roundtrip(t *testing.T) {
	original := `{"email":"test@example.com","custom":"value"}`

	var n Notification
	require.NoError(t, json.Unmarshal([]byte(original), &n))

	data, err := json.Marshal(n)
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "test@example.com", result["email"])
	assert.Equal(t, "value", result["custom"])
}
