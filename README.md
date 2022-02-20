# mobilepay-go

The unofficial MobilePay Go client library.

[![Build Status](https://github.com/steffen25/mobilepay-go/workflows/golangci-lint/badge.svg)](https://github.com/steffen25/mobilepay-go/actions?query=workflow%3Agolangci-lint)
[![Test suite Status](https://github.com/steffen25/mobilepay-go/workflows/test-suite/badge.svg)](https://github.com/steffen25/mobilepay-go/actions?query=workflow%3Agolangci-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/steffen25/mobilepay-go)](https://goreportcard.com/report/github.com/steffen25/mobilepay-go)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/steffen25/mobilepay-go)
[![codecov](https://codecov.io/gh/steffen25/mobilepay-go/branch/master/graph/badge.svg)](https://codecov.io/gh/steffen25/mobilepay-go)

## Installation
```shell
$ go get -u github.com/steffen25/mobilepay-go
```

Then, reference mobilepay-go in a Go program with `import`:

``` go
import (
    "github.com/steffen25/mobilepay-go"
)
```

## Documentation

Below are a few simple examples:

### Initialize a Mobilepay client

```go
cfg := &mobilepay.Config{
    HTTPClient: nil,
    Logger:     nil,
    URL:        mobilepay.DefaultBaseURL,
}

mp := mobilepay.New("client_id", "api_key", cfg)
```
The config properties allows you to use custom configured values if you would like.
Otherwise the client will fall back to use a default value. 

E.g. you might want to use a HTTP client that automatically retries request if a server error occurs.
https://github.com/hashicorp/go-retryablehttp

The library will default to the production base URL if not set.
You can use the constants defined by the mobilepay package if you would like to try out the sandbox environment (highly recommended).
Use either the `mobilepay.DefaultBaseURL` or `mobilepay.TestBaseUrl`.

All the examples below will use the reference `mp` as a reference to the client.

### Payments
 Get all payments

```go
opts := mobilepay.ListOptions{
    PageSize:   10,
    PageNumber: 1,
}

ctx := context.TODO()

payments, err := mp.Payment.Get(ctx, opts)
```

Get single payment details

```go
ctx := context.TODO()

payment, err := mp.Payment.Find(ctx, "payment_id")
```

Create payment
```go
params := &mobilepay.PaymentParams{
    Amount:         1050,
    IdempotencyKey: "223dbe4e-2d3b-4484-b870-6c86cff2c07b",
    PaymentPointId: "8f8ed17f-f310-4f40-a7a4-df78185efbdd",
    RedirectUri:    "myapp://redirect",
    Reference:      "payment-1",
    Description:    "test payment #1",
}

ctx := context.TODO()

payment, err := mp.Payment.Create(ctx, params)
```

Capture payment
```go
ctx := context.TODO()

err := mp.Payment.Capture(ctx, "payment_id", 1050)
```
The amount is specified as an integer and is in cents which in danish terms is 'ører'.

Cancel payment
```go
ctx := context.TODO()

err := mp.Payment.Cancel(ctx, "payment_id")
```

List payment refunds
```go
opts := &mobilepay.RefundsListOptions{
    ListOptions: mobilepay.ListOptions{
        PageSize:   10,
        PageNumber: 1,
    },
    PaymentId:      "payment_id",
    PaymentPointId: "payment_point_id",
    CreatedBefore:  "2020-01-02T15:04",
    CreatedAfter:   "2021-01-02T15:04",
}

ctx := context.TODO()

err := mp.Payment.Refunds(ctx, opts)
```

Create payment refund
```go
params := &mobilepay.PaymentRefundParams{
    PaymentId:      "223dbe4e-2d3b-4484-b870-6c86cff2c07b",
    Amount:         550
    IdempotencyKey: "223dbe4e-2d3b-4484-b870-6c86cff2c07b",
    Reference:      "payment-1",
    Description:    "this is a test payment",
}

ctx := context.TODO()

err := mp.Payment.Refunds(ctx, params)
```

### Webhooks

Get single webhook

```go
ctx := context.TODO()

webhook, err := mp.Webhook.Find(ctx, "webhook_id")
```

Get all merchant's webhooks

```go
ctx := context.TODO()

webhooks, err := mp.Webhook.Get(ctx)
```

Create webhook

```go
ctx := context.TODO()

// configure these events to match your needs.
events := []mobilepay.WebhookEvent{
    mobilepay.PaymentReserved.Name(),
    mobilepay.PaymentExpired.Name(),
}

// change the base url accordingly
params := &mobilepay.WebhookCreateParams{
    Events: events,
    Url:    "https://my-api.com/webhooks",
}

webhook, err := mp.Webhook.Create(ctx, params)
```

Update webhook

```go
ctx := context.TODO()

events := []mobilepay.WebhookEvent{
    mobilepay.PaymentPointActivated.Name(),
}

params := &mobilepay.WebhookUpdateParams{
    Events: events,
    Url:    "https://my-api.com/webhooks",
}

webhook, err := mp.Webhook.Update(ctx, "webhook_id", params)
```

Delete webhook

```go
ctx := context.TODO()

err := mp.Webhook.Delete(ctx, "webhook_id")
```

### Verifying webhooks
This library comes with a built in webhook verifier that you can use ensure webhooks was sent by MobilePay.

The example below illustrates how you can verify incoming MobilePay webhooks using the standard library but this verifier could be used with any of the major router libraries out there.

```go
func main() {
    mux := http.NewServeMux()
    
    successHandler := http.HandlerFunc(mobilepaySuccesHandler)
    mux.Handle("/mobilepay/webhooks", mobilepayWebhooks(successHandler, "webhook_url", "webhook_signature_key"))
}

// middleware to verify webhook was sent by MobilePay.
func mobilepayWebhooks(next http.Handler, webhookUrl, webhookSignature string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        verifier, err := mobilepay.NewWebhooksVerifier(r.Header, webhookUrl, webhookSignature)
        if err != nil {
            // handle error
        }
        
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        
        // rewind the body
        r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
        
        _, err = verifier.Write(body)
        if err != nil {
            // handle error
        }
        
        err = verifier.Ensure()
        if err != nil {
            // could not verify mobilepay signature.
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        
        // Signature is OK.
        // Go to next handler in the call stack (mobilepaySuccesHandler)
        next.ServeHTTP(w, r)
    })
}

func mobilepaySuccesHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Webhook OK"))
}
```

# Contributing
You are more than welcome to contribute to this project. Fork and make a Pull Request, or create an Issue if you see any problem.


# MobilePay documentation
- [Payments API Docs](https://mobilepaydev.github.io/MobilePay-Payments-API/docs/introduction)
- [Webhooks API Docs](https://developer.mobilepay.dk/node/5121)

# License
MIT
