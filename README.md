# Go MobilePay

The unofficial MobilePay Go client library.

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

### Payments
 List payments

```go
opts := mobilepay.ListOptions{
PageSize:   10,
PageNumber: 1,
}

ctx := context.TODO()

data, err := mobilepay.Payment.List(ctx, opts)
```

Get payment

```go
ctx := context.TODO()

payment, err := mobilepay.Payment.Get(ctx, "223dbe4e-2d3b-4484-b870-6c86cff2c07b/")
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

payment, err := mobilepay.Payment.Create(ctx, params)
```

Cancel payment
```go
ctx := context.TODO()

err := mobilepay.Payment.Cancel(ctx, "223dbe4e-2d3b-4484-b870-6c86cff2c07b")
```

Capture payment
```go
ctx := context.TODO()

err := mobilepay.Payment.Capture(ctx, "223dbe4e-2d3b-4484-b870-6c86cff2c07b", 1050)
```

List payment refunds
```go
opts := &mobilepay.RefundsListOptions{
    ListOptions: mobilepay.ListOptions{
        PageSize:   10,
        PageNumber: 1,
    },
    PaymentId:      "223dbe4e-2d3b-4484-b870-6c86cff2c07b",
    PaymentPointId: "8f8ed17f-f310-4f40-a7a4-df78185efbdd",
    CreatedBefore:  "2020-01-02T15:04",
    CreatedAfter:   "2022-01-02T15:04",
}

ctx := context.TODO()

err := mobilepay.Payment.Refunds(ctx, opts)
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

err := mobilepay.Payment.Refunds(ctx, params)
```

# Contributing
You are more than welcome to contribute to this project. Fork and make a Pull Request, or create an Issue if you see any problem.

# License
MIT
