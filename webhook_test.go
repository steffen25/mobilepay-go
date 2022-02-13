package mobilepay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"testing"
)

func TestWebhooks_List(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/list_webhooks.json")
	if err != nil {
		t.Fatal(err)
	}

	gock.New(TestBaseUrl).
		Get("/v1/webhooks").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)

	ctx := context.TODO()

	data, err := client.Webhook.List(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(data.Webhooks))
}

func TestWebhooks_Create(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/create_webhook.json")
	if err != nil {
		t.Fatal(err)
	}
	events := []WebhookEvent{PaymentReserved.Name(), PaymentExpired.Name()}
	eventsJson, err := json.Marshal(events)
	if err != nil {
		t.Fatal(err)
	}
	eventsData := fmt.Sprintf("%v", string(eventsJson))

	testdata = bytes.Replace(testdata, []byte("EVENTS"), []byte(eventsData), 1)
	testdata = bytes.Replace(testdata, []byte("WEBHOOK_URL"), []byte("https://my-api.com/webhooks"), 1)

	gock.New(TestBaseUrl).
		Post("/v1/webhooks").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)

	params := &WebhookCreateParams{
		Events: events,
		Url:    "https://my-api.com/webhooks",
	}
	ctx := context.TODO()

	webhook, err := client.Webhook.Create(ctx, params)
	assert.Nil(t, err)
	assert.Equal(t, "https://my-api.com/webhooks", webhook.Url)
	assert.Equal(t, events, webhook.Events)
}

func TestWebhooks_Get(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/get_webhook.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("WEBHOOK_ID"), []byte("e4a2e195-74f6-42e1-a172-83291c9d2a41"), 1)
	testdata = bytes.Replace(testdata, []byte("WEBHOOK_URL"), []byte("http://localhost/callback"), 1)
	testdata = bytes.Replace(testdata, []byte("SIGNATURE_KEY"), []byte("secret"), 1)

	gock.New(TestBaseUrl).
		Get("/v1/webhooks/e4a2e195-74f6-42e1-a172-83291c9d2a41").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)
	ctx := context.TODO()

	webhook, err := client.Webhook.Get(ctx, "e4a2e195-74f6-42e1-a172-83291c9d2a41")

	assert.Nil(t, err)
	assert.Equal(t, "e4a2e195-74f6-42e1-a172-83291c9d2a41", webhook.WebhookId)
	assert.Equal(t, "http://localhost/callback", webhook.Url)
	assert.Equal(t, "secret", webhook.SignatureKey)
}

func TestWebhooks_Update(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/get_webhook.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("WEBHOOK_ID"), []byte("e4a2e195-74f6-42e1-a172-83291c9d2a41"), 1)
	testdata = bytes.Replace(testdata, []byte("WEBHOOK_URL"), []byte("http://localhost/webhook"), 1)
	testdata = bytes.Replace(testdata, []byte("SIGNATURE_KEY"), []byte("secret"), 1)

	gock.New(TestBaseUrl).
		Put("/v1/webhooks/e4a2e195-74f6-42e1-a172-83291c9d2a41").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)
	ctx := context.TODO()

	params := &WebhookUpdateParams{
		Events: []WebhookEvent{PaymentReserved.Name()},
		Url:    "http://localhost/webhook",
	}

	webhook, err := client.Webhook.Update(ctx, "e4a2e195-74f6-42e1-a172-83291c9d2a41", params)

	assert.Nil(t, err)
	assert.Equal(t, "e4a2e195-74f6-42e1-a172-83291c9d2a41", webhook.WebhookId)
	assert.Equal(t, "http://localhost/webhook", webhook.Url)
	assert.Equal(t, "secret", webhook.SignatureKey)
}

func TestWebhooks_Delete(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New(TestBaseUrl).
		Delete("/v1/webhooks/e4a2e195-74f6-42e1-a172-83291c9d2a41").
		Reply(204)

	client := NewClient("test", "test", config)
	ctx := context.TODO()

	err := client.Webhook.Delete(ctx, "e4a2e195-74f6-42e1-a172-83291c9d2a41")

	assert.Nil(t, err)
}
