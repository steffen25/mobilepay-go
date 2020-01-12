package client

import (
	"github.com/steffen25/mobilepay-go"
	"github.com/steffen25/mobilepay-go/appswitch"
)

// MobilePay is the MobilePay client. It contains all the different MobilePay API resources available.
type MobilePay struct {
	// Config is the config used to configure the different MobilePay APIs.
	Config *mobilepay.Config
	// AppSwitch is the client used to invoke AppSwitch API endpoints.
	AppSwitch *appswitch.Client
}

// Init initializes the Mobilepay client with the appropriate config
// as well as giving the ability to override the different backend as needed
func (mp *MobilePay) Init(cfg *mobilepay.Config, backends *mobilepay.Backends) {
	if backends == nil {
		backends = mobilepay.NewBackends(cfg, mobilepay.NewDefaultHTTPClient())
	}
	mp.Config = cfg
	mp.AppSwitch = &appswitch.Client{
		Backend:    backends.AppSwitch,
		MerchantID: cfg.MerchantID,
	}
}

// New initializes the Mobilepay client with the appropriate config
// as well as giving the ability to override the different backend as needed
func New(cfg *mobilepay.Config, backends *mobilepay.Backends) *MobilePay {
	api := MobilePay{Config: cfg}
	api.Init(cfg, backends)

	return &api
}
