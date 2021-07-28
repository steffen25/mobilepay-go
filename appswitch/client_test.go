package appswitch

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt"
	"github.com/steffen25/mobilepay-go"
	"github.com/stretchr/testify/assert"
	jose "gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func setup() (c *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/appswitch/api/v1/", http.StripPrefix("/appswitch/api/v1", mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	const privPEM = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`

	const pubPEM = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----`

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
	if err != nil {
		log.Fatalf("could not parse private key %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
	if err != nil {
		log.Fatalf("could not parse public key %v", err)
	}
	// Create a payload signer
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, nil)
	if err != nil {
		log.Fatalf("could not create signer %v", err)
	}

	url, err := url.Parse(server.URL + "/appswitch/api/v1/")
	if err != nil {
		log.Fatalf("could not parse HTTP record server url %v", err)
	}

	cfg := &mobilepay.BackendConfig{
		AppConfig: mobilepay.NewConfig("1234", "1234",
			mobilepay.OptionPrivateKey(privKey),
			mobilepay.OptionPublicKey(pubKey),
			mobilepay.OptionSigner(signer)),
		HTTPClient: mobilepay.NewDefaultHTTPClient(),
		URL:        url.String(),
		TestMode:   false,
	}

	client := &Client{
		Backend:    mobilepay.NewBackendWithConfig(mobilepay.AppSwitchBackend, cfg),
		MerchantID: "1234",
	}

	return client, mux, server.URL, server.Close
}

func TestClient_GetPaymentStatus(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/merchants/1234/orders/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
		  "LatestPaymentStatus": "Captured",
		  "TransactionId": "61872634691623746",
		  "OriginalAmount": 123.45
		}`)
	})
	status, err := client.GetPaymentStatus("1234")
	assert.Nil(t, err)
	assert.NotNil(t, status)
}

func TestClient_GetTransactions(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/merchants/1234/orders/1234/transactions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
		  {
			"TimeStamp": "2016-04-08T07:45:36.533",
			"PaymentStatus": "Captured",
			"TransactionId": "61872634691623746",
			"Amount": 11.5
		  }
		]`)
	})
	transactions, err := client.GetTransactions("1234")
	assert.Nil(t, err)
	assert.NotNil(t, transactions)
}

func TestClient_GetReservations(t *testing.T) {
	client, mux, _, _ := setup()
	from, _ := time.Parse(UrlQueryTimestampLayout, "2020-10-03T20_53")
	to, _ := time.Parse(UrlQueryTimestampLayout, "2020-10-04T20_53")
	endpointParams := &GetReservationsParams{
		From:       from,
		To:         to,
		CustomerID: "",
	}
	mux.HandleFunc("/reservations/merchants/1234/2020-10-03T20_53/2020-10-04T20_53", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
		  {
			"TimeStamp": "2016-04-08T07_45_36Z",
			"OrderId": "DB TESTING 2015060908",
			"TransactionId": "61872634691623746",
			"Amount": 100.25,
			"CaptureType": "Full"
		  },
          {
			"TimeStamp": "2016-04-08T07_45_36Z",
			"OrderId": "DB TESTING 2015060908",
			"TransactionId": "61872634691623799",
			"Amount": 100.25,
			"CaptureType": "Partial"
		  }
		]`)
	})
	transactions, err := client.GetReservations(endpointParams)
	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
}

func TestClient_CancelReservation(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/reservations/merchants/1234/orders/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
    		"TransactionId" : "61872634691623746"
		}`)
	})
	reservation, err := client.CancelReservation("1234")
	assert.Nil(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, "61872634691623746", reservation.TransactionID)
}

func TestClient_Refund(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/merchants/1234/orders/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"TransactionId" : "61872634691623757",
			"OriginalTransactionId" : "61872634691623746",
			"Remainder" : 20.00
		}`)
	})
	refundParams := &RefundParams{
		Amount:  100.00,
		BulkRef: "123456789",
	}
	refundResult, err := client.Refund("1234", refundParams)
	assert.Nil(t, err)
	assert.NotNil(t, refundResult)
}

func TestClient_Capture(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/reservations/merchants/1234/orders/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"TransactionId": "61872634691623746"}`)
	})
	captureParams := &CaptureParams{
		Amount: 100.00,
	}
	transaction, err := client.Capture("1234", captureParams)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, "61872634691623746", transaction.TransactionID)
}

func TestClient_Capture_With_400Error(t *testing.T) {
	client, mux, _, _ := setup()
	mux.HandleFunc("/reservations/merchants/1234/orders/1234", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"Reason":"InvalidAmount"}`)
	})
	captureParams := &CaptureParams{
		Amount: 1000.00,
	}
	transaction, err := client.Capture("1234", captureParams)
	assert.NotNil(t, err)
	assert.IsType(t, &mobilepay.BadRequestError{}, err)
	assert.Empty(t, transaction)
}
