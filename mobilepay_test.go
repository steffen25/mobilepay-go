package mobilepay

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewBackendWithConfig(t *testing.T) {
	{
		backend := NewBackendWithConfig(AppSwitchBackend, &BackendConfig{
			URL: "https://api.mobilepay.dk/v1",
		}).(*BackendImplementation)
		assert.Equal(t, "https://api.mobilepay.dk/v1", backend.URL)
	}

	{
		backend := NewBackendWithConfig(AppSwitchBackend, &BackendConfig{
			URL: "", // NewBackendWithConfig will default to the api url for the specificed backend type
		}).(*BackendImplementation)
		assert.Equal(t, "https://api.mobeco.dk/appswitch/api/v1", backend.URL)
	}

	{
		backend := NewBackendWithConfig("bad backend", &BackendConfig{
			URL: "",
		})
		assert.Nil(t, backend)
	}

}

func TestNewConfig(t *testing.T) {
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

	{
		cfg, err := NewConfig("test", "test")
		assert.Nil(t, err)
		expectedConfig := &Config{
			MerchantID:      "test",
			SubscriptionKey: "test",
			AppSwitchSigner: nil,
		}
		assert.Equal(t, expectedConfig, cfg)
	}

	{
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
		assert.Nil(t, err)

		privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
		assert.Nil(t, err)

		cfg, err := NewConfig("test", "test", OptionAppSwitchKeyPair([]byte(pubPEM), []byte(privPEM)))
		assert.Nil(t, err)

		expectedConfig := &Config{
			MerchantID:      "test",
			SubscriptionKey: "test",
			AppSwitchSigner: &AppSwitchSigner{
				PrivateKey: privKey,
				PublicKey:  pubKey,
				Signer:     cfg.AppSwitchSigner.Signer, // since the signer contains a slice we cannot compare its equality if we create one outside.
			},
		}
		assert.Equal(t, expectedConfig, cfg)
	}
}

func TestNewBackends(t *testing.T) {
	cfg, err := NewConfig("test", "test")
	assert.Nil(t, err)
	httpClient := &http.Client{}
	backends := NewBackends(cfg, httpClient)
	assert.Equal(t, httpClient, backends.AppSwitch.(*BackendImplementation).HTTPClient)
}

func TestAddOptions(t *testing.T) {
	{
		var nilPnt *string
		_, err := addOptions(":", nilPnt)
		assert.NoError(t, err)
	}

	{
		_, err := addOptions(":", nil)
		assert.Error(t, err)
	}

	{
		type test interface{}
		var tester test
		_, err := addOptions("/", &tester)
		assert.Error(t, err)
	}

	{
		type test struct {
			CustomerID string `url:"customerId,omitempty"`
		}
		params := &test{CustomerID: "+4588888888"}
		path, err := addOptions("/appswitch/api/v1/reservations/merchants/APPDK0000000000/2020-01-04T19_58/2020-01-04T20_58", params)
		assert.NoError(t, err)
		assert.Equal(t, "/appswitch/api/v1/reservations/merchants/APPDK0000000000/2020-01-04T19_58/2020-01-04T20_58?customerId=%2B4588888888", path)
	}
}

type mockReadCloser struct {
	mock.Mock
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCheckResponseError(t *testing.T) {
	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"data": "data"}`)),
		}

		err := CheckResponseError(res)
		assert.NoError(t, err)
	}

	{
		mockReadCloser := mockReadCloser{}
		mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusBadRequest,
			Body:       &mockReadCloser,
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		mockReadCloser.AssertExpectations(t)
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"statusCode": 401, "message": "Access denied due to missing subscription key. Make sure to include subscription key when making requests to an API."}`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Access denied due to missing subscription key. Make sure to include subscription key when making requests to an API.")
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.IsType(t, ResponseDecodingError{}, err)
		assert.EqualError(t, err, "invalid character '<' looking for beginning of value")
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"Reason": "MerchantNotFound"}`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.EqualError(t, err, "MerchantNotFound")
	}

	{
		headers := make(http.Header)
		headers.Add("Retry-After", "10")

		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusTooManyRequests,
			Body:       ioutil.NopCloser(strings.NewReader(``)),
			Header:     headers,
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.EqualError(t, err, "MobilePay rate limit exceeded")
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(`{"Reason": "BackendError"}`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.IsType(t, &BadRequestError{}, err)
		assert.EqualError(t, err, "BackendError")
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(`{"CorrelationId": "a658ab24-70ab-4d05-a792-e5995f237c10", "Errortype": "ServerError", "Message": "Attempted to divide by zero."}`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.IsType(t, &ServerError{}, err)
		assert.EqualError(t, err, "mobilepay server error detected. Message: Attempted to divide by zero. Type: ServerError CorrelationId a658ab24-70ab-4d05-a792-e5995f237c10")
	}

	{
		res := &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(`{"test": "test"}`)),
		}

		err := CheckResponseError(res)
		assert.Error(t, err)
		assert.IsType(t, &ResponseError{}, err)
		assert.EqualError(t, err, "Unknown error. Status 500 Body: {\"test\": \"test\"}")
	}
}
