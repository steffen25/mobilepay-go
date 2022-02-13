package mobilepay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

// https://mobilepaydev.github.io/MobilePay-Payments-API/docs/payments-refunds/create-payments
const (
	libraryVersion       = "1.0.0"
	defaultBaseURL       = "https://api.sandbox.mobilepay.dk/"
	TestBaseUrl          = "https://api.sandbox.mobilepay.dk"
	userAgent            = "mobilepay-go/" + libraryVersion
	mediaType            = "application/json"
	ibmClientIdHeaderKey = "x-ibm-client-id"
	DefaultTimeout       = 10 * time.Second
)

type Client struct {
	// HTTP client used to communicate with the MobilePay App Payment API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Services used for communicating with the API

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string

	Logger LeveledLoggerInterface

	// MobilePay API services used for communicating with the API.
	Payment *PaymentServiceOp // we are using a struct over an interface to support multiple interfaces implemented by the struct properties.
	Webhook WebhookService
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ClientOpt are options for New.
type ClientOpt func(*Client) error

type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`
}

type responseParser func(*http.Response) error

func newJSONParser(resource interface{}) responseParser {
	return func(res *http.Response) error {
		return json.NewDecoder(res.Body).Decode(resource)
	}
}

type Config struct {
	HTTPClient *http.Client
	Logger     LeveledLoggerInterface
}

func NewClient(IbmClientId, apiKey string, config *Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = newDefaultHTTPClient()
	}

	if config.Logger == nil {
		config.Logger = DefaultLeveledLogger
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    config.HTTPClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Logger:    config.Logger,
	}

	// we wrap the refund service inside the payment service to follow a more RESTful approach
	refundService := &RefundServiceOp{client: c}

	c.Payment = &PaymentServiceOp{client: c, Refund: refundService}
	c.Webhook = &WebhookServiceOp{client: c}

	c.headers = make(map[string]string)

	c.headers[ibmClientIdHeaderKey] = IbmClientId
	c.headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)

	return c
}

// New returns a new MobilePay App Payment API client instance.
func New(config *Config, opts ...ClientOpt) (*Client, error) {
	c := NewClient("", "", config)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func SetIbmClientId(clientId string) ClientOpt {
	return func(c *Client) error {
		c.headers[ibmClientIdHeaderKey] = clientId
		return nil
	}
}

func SetApiKey(apiKey string) ClientOpt {
	return func(c *Client) error {
		c.headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
		return nil
	}
}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	var req *http.Request

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	c.Logger.Infof("Requesting %v %v%v", req.Method, req.URL.Host, req.URL.Path)

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)

	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()

	return origURL.String(), nil
}

func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	//requestDump, err := httputil.DumpResponse(resp, true)
	//if err != nil {
	//	return nil, err
	//}
	//log.Println(string(requestDump))

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, copyErr := io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
			if copyErr != nil {
				err = copyErr
			}
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = newJSONParser(v)(resp)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
