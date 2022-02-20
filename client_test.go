package mobilepay

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/url"

	"net/http"
	"testing"
	"time"
)

func TestClient_Timeout(t *testing.T) {
	var config = &Config{}

	client := New("test", "test", config)

	assert.Equal(t, 10*time.Second, client.client.Timeout)
}

func TestClient_Auth_Keys(t *testing.T) {
	var config = &Config{}

	client := New("client_id", "api_key", config)

	assert.Len(t, client.headers, 2)
	assert.Contains(t, client.headers, "x-ibm-client-id")
	assert.Contains(t, client.headers, "Authorization")
	assert.Equal(t, "client_id", client.headers["x-ibm-client-id"])
	assert.Equal(t, "Bearer api_key", client.headers["Authorization"])
}

func TestClient_CheckResponse(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/capture_payment_409_amount_too_large.json")
	if err != nil {
		t.Fatal(err)
	}

	gock.New(TestBaseUrl).
		Post("/v1/payments/25df9ee7-5608-4b7a-98d0-df649861075b/capture").
		Reply(409).
		JSON(testdata)

	url := fmt.Sprintf("%s/%s", TestBaseUrl, "/v1/payments/25df9ee7-5608-4b7a-98d0-df649861075b/capture")
	res, err := http.Post(url, "application/json", nil)
	assert.Nil(t, err)

	err = CheckResponse(res)
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	mpError, ok := err.(*ErrorResponse)
	assert.True(t, ok)
	assert.Equal(t, 409, mpError.StatusCode)
	assert.Equal(t, "amount_too_large", mpError.Conflict.Code)
}

func TestClient_CheckResponse_Unknown_Error_Format(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New(TestBaseUrl).
		Post("/v1/payments/25df9ee7-5608-4b7a-98d0-df649861075b/capture").
		Reply(409).
		JSON("Unknown error")

	url := fmt.Sprintf("%s/%s", TestBaseUrl, "/v1/payments/25df9ee7-5608-4b7a-98d0-df649861075b/capture")
	res, err := http.Post(url, "application/json", nil)
	assert.Nil(t, err)

	err = CheckResponse(res)
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	mpError, ok := err.(*ErrorResponse)
	assert.True(t, ok)
	assert.Equal(t, 409, mpError.StatusCode)
	assert.Equal(t, "Unknown error", mpError.Message)
}

func TestClient_ErrorResponse(t *testing.T) {
	error := ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusConflict,
			Request: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "https",
					Host:   "bla.com",
				},
			},
		},
		Message: "Unknown error",
		Conflict: ConflictError{
			Code:          "amount_too_large",
			Message:       "Cannot capture a larger amount than is reserved.",
			CorrelationID: "d503b7ed-b5d0-4751-b3ac-52ecd7cd3a4a",
			Origin:        "MPY",
		},
		StatusCode: http.StatusConflict,
	}

	expected := fmt.Sprintf("%s %s: %d %v", "POST", "https://bla.com", 409, "Unknown error")

	assert.Equal(t, expected, error.Error())
}
