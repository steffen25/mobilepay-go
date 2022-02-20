package mobilepay

import (
	"crypto/hmac"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"hash"
	"net/http"
)

// Mobilepay signature header
const mpSignature = "x-mobilepay-signature"

type WebhooksVerifier struct {
	webhookUrl []byte
	signature  []byte
	hmac       hash.Hash
}

// NewWebhooksVerifier is a helper function to verify incoming webhooks from Mobilepay.
// See https://mobilepaydev.github.io/MobilePay-Payments-API/docs/webhooks-api
// webhookUrl is your webhook url that you used to create the webhook.
// webhookSignatureKey is returned by Mobilepay when you create a webhook.
func NewWebhooksVerifier(header http.Header, webhookUrl, webhookSignatureKey string) (wv WebhooksVerifier, err error) {
	signature := header.Get(mpSignature)
	hash := hmac.New(sha1.New, []byte(webhookSignatureKey))

	if webhookUrl == "" || signature == "" {
		return WebhooksVerifier{}, ErrMissingVerifierProperties
	}

	if _, err := hash.Write([]byte(webhookUrl)); err != nil {
		return WebhooksVerifier{}, err
	}

	return WebhooksVerifier{
		webhookUrl: []byte(webhookUrl),
		signature:  []byte(signature),
		hmac:       hash,
	}, nil
}

func (v *WebhooksVerifier) Write(body []byte) (n int, err error) {
	return v.hmac.Write(body)
}

func (v WebhooksVerifier) Ensure() error {
	computed := v.hmac.Sum(nil)
	sEnc := b64.StdEncoding.EncodeToString(computed)

	// constant time compare in order to prevent leaking information
	if hmac.Equal([]byte(sEnc), v.signature) {
		return nil
	}

	return fmt.Errorf("Computed unexpected signature of: %s", sEnc)
}
