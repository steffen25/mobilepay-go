package mobilepay

import (
	"context"
	"fmt"
	"net/http"
)

// subscribe/create webhook

const webhooksBasePath = "v1/webhooks"

type GetWebhookRequest struct {
}

type WebhookEventEnum int
type WebhookEvent string

const (
	Unknown WebhookEventEnum = iota
	PaymentReserved
	PaymentExpired
	PaymentPointActivated
)

func (webhookEvent WebhookEventEnum) Name() WebhookEvent {
	return [...]WebhookEvent{"Unknown", "payment.reserved", "payment.expired", "paymentpoint.activated"}[webhookEvent]
}

type WebhookService interface {
	List(context.Context) (*WebhooksRoot, error)
	Create(context.Context, *WebhookCreateParams) (*Webhook, error)
	Get(context.Context, string) (*Webhook, error)
	Update(context.Context, string, *WebhookUpdateParams) (*Webhook, error)
	Delete(context.Context, string) error
}

type WebhookServiceOp struct {
	client *Client
}

var _ WebhookService = &WebhookServiceOp{}

type GetWebhookResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// WebhookCreateParams represents a request to create a payment.
type WebhookCreateParams struct {
	// List of subscribed events.
	Events []WebhookEvent `json:"events"`
	// URL to where webhook requests will be sent. Must be HTTPS. Scheme and host will be converted to lower case. Result can be seen in the response.
	Url string `json:"url"`
}

// WebhookUpdateParams represents a request to update a webhook record.
type WebhookUpdateParams struct {
	Url    string         `json:"url"`
	Events []WebhookEvent `json:"events"`
}

type Webhook struct {
	WebhookId    string         `json:"webhookId"`
	SignatureKey string         `json:"signatureKey"`
	Url          string         `json:"url"`
	Events       []WebhookEvent `json:"events"`
}

// webhooksRoot represents a response from the MobilePay App Payment API
type WebhooksRoot struct {
	Webhooks []Webhook `json:"webhooks"`
}

// List all webhooks.
func (s WebhookServiceOp) List(ctx context.Context) (*WebhooksRoot, error) {
	path := webhooksBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(WebhooksRoot)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// Create webhook
func (s *WebhookServiceOp) Create(ctx context.Context, createRequest *WebhookCreateParams) (*Webhook, error) {
	if createRequest == nil {
		return nil, NewArgError("createRequest", "cannot be nil")
	}

	path := webhooksBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, err
	}

	root := new(Webhook)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

// Get individual webhook. It requires a non-empty webhook id.
func (s *WebhookServiceOp) Get(ctx context.Context, webhookId string) (*Webhook, error) {
	path := fmt.Sprintf("%s/%s", webhooksBasePath, webhookId)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(Webhook)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (s WebhookServiceOp) Update(ctx context.Context, webhookId string, request *WebhookUpdateParams) (*Webhook, error) {
	path := fmt.Sprintf("%s/%s", webhooksBasePath, webhookId)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return nil, err
	}

	root := new(Webhook)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (s *WebhookServiceOp) Delete(ctx context.Context, webhookId string) error {
	path := fmt.Sprintf("%s/%s", webhooksBasePath, webhookId)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
