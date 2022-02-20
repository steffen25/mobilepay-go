package mobilepay

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

const (
	validUrl      = "https://webhook.site/080a55d2-ff87-4494-a05c-e0e3beb78134"
	validSecret   = "4aa30d41-4368-47a0-b4ef-6a83dc8be5d6"
	invalidSecret = "invalid-signature-key"
	validBody     = `{"notificationId":"4352f1ae-59c3-430c-a402-d74641dd8555","eventType":"test.notification","eventDate":"2022-02-20T16:35:28Z","data":{"type":"test","id":"57ff4ddf-575f-4c4a-99c8-b190a1e1f316"}}`
	invalidBody   = `{"notificationId":"1234f1ae-59c3-430c-a402-d74641dd8555","eventType":"test.notification","eventDate":"2022-02-20T16:35:28Z","data":{"type":"test","id":"57ff4ddf-575f-4c4a-99c8-b190a1e1f316"}}`
)

func newHeader(valid bool) http.Header {
	h := http.Header{}
	if valid {
		h.Set("x-mobilepay-signature", "HIcf0Ivp0HwjB2qVIwU1vIdf/60=")
	} else {
		h.Set("x-mobilepay-signature", "")
	}

	return h
}

func TestNewWebhooksVerifier_Empty_Webhook_Url(t *testing.T) {
	header := newHeader(false)
	_, err := NewWebhooksVerifier(header, "", validSecret)
	assert.Error(t, err)
}

func TestWebhooksVerifier_Ensure_Valid_Secret_Valid_Body(t *testing.T) {
	header := newHeader(true)
	verifier, err := NewWebhooksVerifier(header, validUrl, validSecret)
	assert.Nil(t, err)
	_, err = io.WriteString(&verifier, validBody)
	assert.Nil(t, err)
	err = verifier.Ensure()
	assert.Nil(t, err)
}

func TestWebhooksVerifier_Ensure_Valid_Secret_Invalid_Body(t *testing.T) {
	header := newHeader(true)
	verifier, err := NewWebhooksVerifier(header, validUrl, validSecret)
	assert.Nil(t, err)
	_, err = io.WriteString(&verifier, invalidBody)
	assert.Nil(t, err)
	err = verifier.Ensure()
	assert.Error(t, err)
}

func TestWebhooksVerifier_Ensure_Invalid_Secret_Valid_Body(t *testing.T) {
	header := newHeader(true)
	verifier, err := NewWebhooksVerifier(header, validUrl, invalidSecret)
	assert.Nil(t, err)
	_, err = io.WriteString(&verifier, validBody)
	assert.Nil(t, err)
	err = verifier.Ensure()
	assert.Error(t, err)
}
