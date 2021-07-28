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
	const privWithPasswordPEM = `
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,7487BB8910A3741B

iL7m48mbFSIy1Y5xbXWwPTR07ufxu7o+myGUE+AdDeWWISkd5W6Gl44oX/jgXldS
mL/ntUXoZzQz2WKEYLwssAtSTGF+QgSIMvV5faiP+pLYvWgk0oVr42po00CvADFL
eDAJC7LgagYifS1l4EAK4MY8RGCHyJWEN5JAr0fc/Haa3WfWZ009kOWAp8MDuYxB
hQlCKUmnUpXCp5c6jwbjlyinLj8XwzzjZ/rVRsY+t2Z0Vcd5qzR5BV8IJCqbG5Py
z15/EFgMG2N2eYMsiEKgdXeKW2H5XIoWyun/3pBigWaDnTtiWSt9kz2MplqYfIT7
F+0XE3gdDGalAeN3YwFPHCkxxBmcI+s6lQG9INmf2/gkJQ+MOZBVXKmGLv6Qis3l
0eyUz1yZvNzf0zlcUBjiPulLF3peThHMEzhSsATfPomyg5NJ0X7ttd0ybnq+sPe4
qg2OJ8qNhYrqnx7Xlvj61+B2NAZVHvIioma1FzqX8DxQYrnR5S6DJExDqvzNxEz6
5VPQlH2Ig4hTvNzla84WgJ6USc/2SS4ehCReiNvfeNG9sPZKQnr/Ss8KPIYsKGcC
Pz/vEqbWDmJwHb7KixCQKPt1EbD+/uf0YnhskOWM15YiFbYAOZKJ5rcbz2Zu66vg
GAmqcBsHeFR3s/bObEzjxOmMfSr1vzvr4ActNJWVtfNKZNobSehZiMSHL54AXAZW
Yj48pwTbf7b1sbF0FeCuwTFiYxM+yiZVO5ciYOfmo4HUg53PjknKpcKtEFSj02P1
8JRBSb++V0IeMDyZLl12zgURDsvualbJMMBBR8emIpF13h0qdyah431gDhHGBnnC
J5UDGq21/flFjzz0x/Okjwf7mPK5pcmF+uW7AxtHqws6m93yD5+RFmfZ8cb/8CL8
jmsQslj+OIE64ykkRoJWpNBKyQjL3CnPnLmAB6TQKxegR94C7/hP1FvRW+W0AgZy
g2QczKQU3KBQP18Ui1HTbkOUJT0Lsy4FnmJFCB/STPRo6NlJiATKHq/cqHWQUvZd
d4oTMb1opKfs7AI9wiJBuskpGAECdRnVduml3dT4p//3BiP6K9ImWMSJeFpjFAFs
AbBMKyitMs0Fyn9AJRPl23TKVQ3cYeSTxus4wLmx5ECSsHRV6g06nYjBp4GWEqSX
RVclXF3zmy3b1+O5s2chJN6TrypzYSEYXJb1vvQLK0lNXqwxZAFV7Roi6xSG0fSY
EAtdUifLonu43EkrLh55KEwkXdVV8xneUjh+TF8VgJKMnqDFfeHFdmN53YYh3n3F
kpYSmVLRzQmLbH9dY+7kqvnsQm8y76vjug3p4IbEbHp/fNGf+gv7KDng1HyCl9A+
Ow/Hlr0NqCAIhminScbRsZ4SgbRTRgGEYZXvyOtQa/uL6I8t2NR4W7ynispMs0QL
RD61i3++bQXuTi4i8dg3yqIfe9S22NHSzZY/lAHAmmc3r5NrQ1TM1hsSxXawT5CU
anWFjbH6YQ/QplkkAqZMpropWn6ZdNDg/+BUjukDs0HZrbdGy846WxQUvE7G2bAw
IFQ1SymBZBtfnZXhfAXOHoWh017p6HsIkb2xmFrigMj7Jh10VVhdWg==
-----END RSA PRIVATE KEY-----
`

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

	{
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
		assert.Nil(t, err)

		privWithPasswordKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword([]byte(privWithPasswordPEM), "password")
		assert.Nil(t, err)

		cfg, err := NewConfig("test", "test", OptionAppSwitchKeyPairWithPassword([]byte(pubPEM), []byte(privWithPasswordPEM), "password"))
		assert.Nil(t, err)

		expectedConfig := &Config{
			MerchantID:      "test",
			SubscriptionKey: "test",
			AppSwitchSigner: &AppSwitchSigner{
				PrivateKey: privWithPasswordKey,
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
