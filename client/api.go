package client

import (
	"github.com/steffen25/mobilepay-go"
	"github.com/steffen25/mobilepay-go/appswitch"
)

type MobilePay struct {
	Config    *mobilepay.Config
	AppSwitch *appswitch.Client
}

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

func New(cfg *mobilepay.Config, backends *mobilepay.Backends) *MobilePay {
	api := MobilePay{Config: cfg}
	api.Init(cfg, backends)

	return &api
}
