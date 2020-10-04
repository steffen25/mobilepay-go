package mobilepay

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/square/go-jose.v2"
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
	{
		cfg := NewConfig("test", "test")
		expectedConfig := &Config{
			MerchantID:      "test",
			SubscriptionKey: "test",
			PrivateKey:      nil,
			PublicKey:       nil,
			AppSwitchSigner: nil,
		}
		assert.Equal(t, expectedConfig, cfg)
	}

	{
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		pubKey, _ := privateKey.Public().(*rsa.PublicKey)
		signer, _ := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, nil)
		cfg := NewConfig("test", "test", OptionPrivateKey(privateKey), OptionPublicKey(pubKey), OptionSigner(signer))
		expectedConfig := &Config{
			MerchantID:      "test",
			SubscriptionKey: "test",
			PrivateKey:      privateKey,
			PublicKey:       pubKey,
			AppSwitchSigner: signer,
		}
		assert.Equal(t, expectedConfig, cfg)
	}
}

func TestNewBackends(t *testing.T) {
	cfg := NewConfig("test", "test")
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
