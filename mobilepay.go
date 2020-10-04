package mobilepay

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	AppSwitchBackend                BackendType = "appswitch"
	SubscriptionsBackend            BackendType = "subscriptions"
	AppSwitchAPIURL                 string      = "https://api.mobeco.dk/appswitch/api/v1"
	DefaultTimeout                              = 5 * time.Second
	AppSwitchAuthSignatureHeaderKey string      = "AuthenticationSignature"
	AppSwitchSubscriptionHeaderKey  string      = "Ocp-Apim-Subscription-Key"
	TestModeHeaderKey               string      = "Test-mode"
)

type Config struct {
	MerchantID      string
	SubscriptionKey string
	//IBMClientID     string
	//IBMClientSecret string
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
	AppSwitchSigner jose.Signer
}

type BackendType string

type Backend interface {
	Call(method, path, key string, queryParams, body, resource interface{}) error
}

type BackendConfig struct {
	AppConfig  *Config
	HTTPClient *http.Client
	URL        string
	TestMode   bool
}

type Backends struct {
	AppSwitch Backend
}

type BackendImplementation struct {
	AppConfig  *Config
	Type       BackendType
	URL        string
	HTTPClient *http.Client
	TestMode   bool
}

type responseParser func(*http.Response) error

func newJSONParser(resource interface{}) responseParser {
	return func(res *http.Response) error {
		return json.NewDecoder(res.Body).Decode(resource)
	}
}

// Option defines an option to be applied on a config
// Since the config is used across different backends some backends might need a specific element
type Option func(*Config)

// OptionPrivateKey sets the private key of a config
func OptionPrivateKey(privKey *rsa.PrivateKey) func(*Config) {
	return func(c *Config) {
		c.PrivateKey = privKey
	}
}

// OptionPublicKey sets the public key of a config
func OptionPublicKey(pubKey *rsa.PublicKey) func(*Config) {
	return func(c *Config) {
		c.PublicKey = pubKey
	}
}

// OptionPublicKey sets the public key of a config
func OptionSigner(signer jose.Signer) func(*Config) {
	return func(c *Config) {
		c.AppSwitchSigner = signer
	}
}

func NewConfig(merchantId, subscriptionKey string, options ...Option) *Config {
	cfg := &Config{
		MerchantID:      merchantId,
		SubscriptionKey: subscriptionKey,
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg
}

func NewBackends(cfg *Config, httpClient *http.Client) *Backends {
	appSwitchConfig := &BackendConfig{AppConfig: cfg, HTTPClient: httpClient}

	return &Backends{
		AppSwitch: NewBackendWithConfig(AppSwitchBackend, appSwitchConfig),
	}
}

func (b *BackendImplementation) Call(method, path, key string, queryParams, body, resource interface{}) error {
	_path := path
	if queryParams != nil {
		path, err := addOptions(path, queryParams)
		if err != nil {
			return err
		}
		_path = path
	}

	var _bodyBytes []byte
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		_bodyBytes = bodyBytes
	}

	bodyBuffer := bytes.NewBuffer(_bodyBytes)

	req, err := b.NewRequest(method, _path, bodyBuffer, "application/json")
	if err != nil {
		return err
	}

	if b.TestMode {
		err := logRequest(req)
		if err != nil {
			log.Printf("could not log request %v", err)
		}
	}

	if err := b.Do(req, newJSONParser(resource)); err != nil {
		return err
	}

	return nil
}

func (b *BackendImplementation) NewRequest(method, path string, body *bytes.Buffer, contentType string) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = b.URL + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)

	if b.TestMode {
		req.Header.Set(TestModeHeaderKey, "true")
	}

	bodyString := body.String()

	switch b.Type {
	case AppSwitchBackend:
		req.Header.Set(AppSwitchSubscriptionHeaderKey, b.AppConfig.SubscriptionKey)
		err = setMobilePayAuthHeader(req, bodyString, b.AppConfig.AppSwitchSigner, b.AppConfig.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (b *BackendImplementation) Do(req *http.Request, parser responseParser) error {
	res, err := b.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if b.TestMode {
		err := logResponse(res)
		if err != nil {
			log.Printf("could not log response %v", err)
		}
	}

	err = CheckResponseError(res)
	if err != nil {
		return err
	}

	return parser(res)
}

func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}

func NewBackendWithConfig(backendType BackendType, config *BackendConfig) Backend {
	if config.HTTPClient == nil {
		config.HTTPClient = NewDefaultHTTPClient()
	}

	switch backendType {
	case AppSwitchBackend:
		if config.URL == "" {
			config.URL = AppSwitchAPIURL
		}
		return newBackendImplementation(backendType, config)
	}

	return nil
}

// TODO: There must be a better way to do all this error unmarshalling
func CheckResponseError(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return HandleAuthError(bodyBytes, res.StatusCode)
	}

	if res.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(res.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitError{time.Duration(retry) * time.Second}
	}

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return HandleRequestError(bodyBytes, res.StatusCode)
	}

	return HandleServerError(bodyBytes, res.StatusCode)
}

// Used for debugging purposes
func logRequest(req *http.Request) error {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	log.Println(string(requestDump))

	return nil
}

// Used for debugging purposes
func logResponse(res *http.Response) error {
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return err
	}
	log.Println(string(requestDump))

	return nil
}

// addOptions will add query parameters to the string s using the go-querystring library
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// setMobilePayAuthHeader will generate the authentication signature for the HTTP request.
// It will add the header AuthenticationSignature to the request based on the payload(URL and body)
// See https://github.com/MobilePayDev/MobilePay-Merchant-API-Security/blob/master/Merchant-request.md
func setMobilePayAuthHeader(req *http.Request, body string, signer jose.Signer, pubKey *rsa.PublicKey) error {
	// concat url and body
	payload := req.URL.String() + body
	// hash payload
	payloadSha1, err := sha1Hash([]byte(payload))
	if err != nil {
		return err
	}
	// base64 encode hashed payload
	payloadBase64 := base64.StdEncoding.EncodeToString(payloadSha1)
	// generate JSON Web Signature for this payload
	signature, err := generateJWS(signer, pubKey, []byte(payloadBase64))
	if err != nil {
		return err
	}
	signatureCompact, err := signature.CompactSerialize()
	if err != nil {
		return err
	}
	req.Header.Set(AppSwitchAuthSignatureHeaderKey, signatureCompact)

	return nil
}

// sha1Hash generates a SHA-1 hash based on the input and return the digest as a byte array
func sha1Hash(data []byte) ([]byte, error) {
	hasher := sha1.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}

	digest := hasher.Sum(nil)

	return digest, nil
}

func generateJWS(signer jose.Signer, pubKey *rsa.PublicKey, payload []byte) (*jose.JSONWebSignature, error) {
	// Sign a sample payload. Calling the signer returns a protected JWS object,
	// which can then be serialized for output afterwards. An error would
	// indicate a problem in an underlying cryptographic primitive.
	object, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()

	// Parse the serialized, protected JWS object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseSigned(serialized)
	if err != nil {
		return nil, err
	}

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	_, err = object.Verify(pubKey)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func newBackendImplementation(backendType BackendType, config *BackendConfig) Backend {
	return &BackendImplementation{
		AppConfig:  config.AppConfig,
		Type:       backendType,
		URL:        config.URL,
		HTTPClient: config.HTTPClient,
		TestMode:   config.TestMode,
	}
}
