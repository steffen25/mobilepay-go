# mobilepay-go
Go library for accessing various MobilePay APIs.

[![Build Status](https://github.com/steffen25/mobilepay-go/workflows/golangci-lint/badge.svg)](https://github.com/steffen25/mobilepay-go/actions?query=workflow%3Agolangci-lint)
[![Test suite Status](https://github.com/steffen25/mobilepay-go/workflows/test-suite/badge.svg)](https://github.com/steffen25/mobilepay-go/actions?query=workflow%3Agolangci-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/steffen25/mobilepay-go)](https://goreportcard.com/report/github.com/steffen25/mobilepay-go)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/steffen25/mobilepay-go)
[![codecov](https://codecov.io/gh/steffen25/mobilepay-go/branch/master/graph/badge.svg)](https://codecov.io/gh/steffen25/mobilepay-go)

This library is still work in progress.

## Installing
    $ go get -u github.com/steffen25/mobilepay-go
    
## AppSwitch Examples

### Client configuration

  ```go
 import (
        "net/http"
        "time"
        "fmt"

 	    "github.com/steffen25/mobilepay-go"
        "github.com/steffen25/mobilepay-go/client"
 )

const pubPEM = `
-----BEGIN PUBLIC KEY-----
....
-----END PUBLIC KEY-----
`

const privPEM = `
-----BEGIN RSA PRIVATE KEY-----
....
-----END RSA PRIVATE KEY-----
`
 
 func main() {
         // Prepare config
         cfg, err := mobilepay.NewConfig("MERCHANT_ID", "SUBSCRIPTION_KEY",
             mobilepay.OptionAppSwitchKeyPair([]byte(pubPEM), []byte(privPEM)),
         )
         if err != nil {
             // handle error
         }
        
         httpClient := &http.Client{
            Timeout: 5 * time.Second,
         }
         backends := mobilepay.NewBackends(cfg, httpClient)
         mp := client.New(cfg, backends)
        // see examples for how to use the mp client below.
 }
 ```
 
 ### Get payment status
 
  ```go
 import (
        "net/http"
        "time"
        "fmt"

 	    "github.com/steffen25/mobilepay-go"
        "github.com/steffen25/mobilepay-go/client"
 )
 
 func main() {
         // client configuration skipped - see above.
         mp := client.New(cfg, backends)
         // get payment status
         status, err := mp.AppSwitch.GetPaymentStatus("ORDER_ID")
         if err != nil {
            // Handle error
         }
         fmt.Println(status, err)
 }
 ```

## Supported APIs
- [AppSwitch](https://github.com/MobilePayDev/MobilePay-AppSwitch-API)