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
	    "github.com/steffen25/mobilepay-go"
	    "gopkg.in/dgrijalva/jwt-go.v3"
        jose "gopkg.in/square/go-jose.v2"
)

func main() {
        // Parse a RSA key pair used to generate the authentication signature
        privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
        if err != nil {
            // Handle error
        }
        pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
        if err != nil {
            // Handle error
        }
        // Create a payload signer
        signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, nil)
        if err != nil {
            // Handle error
        }
        
        // Prepare config
        cfg := mobilepay.NewConfig("MERCHANT_ID", "SUBSCRIPTION_KEY", 
            mobilepay.OptionPrivateKey(privKey),
            mobilepay.OptionPublicKey(pubKey),
            mobilepay.OptionSigner(signer),
        )
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
 	    "gopkg.in/dgrijalva/jwt-go.v3"
         jose "gopkg.in/square/go-jose.v2"
 )
 
 func main() {
         // Parse a RSA key pair used to generate the authentication signature
         privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
         if err != nil {
             // Handle error
         }
         pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
         if err != nil {
             // Handle error
         }
         signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, nil)
         if err != nil {
             // Handle error
         }
         
         cfg := mobilepay.NewConfig("MERCHANT_ID", "SUBSCRIPTION_KEY", 
             mobilepay.OptionPrivateKey(privKey),
             mobilepay.OptionPublicKey(pubKey),
             mobilepay.OptionSigner(signer),
         )
        
         httpClient := &http.Client{
            Timeout: 5 * time.Second,
         }
         backends := mobilepay.NewBackends(cfg, httpClient)
         mp := client.New(cfg, backends)
         
         status, err := mp.AppSwitch.GetPaymentStatus("ORDER_ID")
         if err != nil {
            // Handle error
         }
         fmt.Println(status, err)
 }
 ```

## Supported APIs
- [AppSwitch](https://github.com/MobilePayDev/MobilePay-AppSwitch-API)