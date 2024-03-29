package mobilepay

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const paymentsBasePath = "v1/payments"

type ListOptions struct {
	PageSize   int `url:"pageSize"`
	PageNumber int `url:"pageNumber"`
}

type RefundsListOptions struct {
	ListOptions
	PaymentId      string `url:"paymentId"`
	PaymentPointId string `url:"paymentPointId,omitempty"`
	CreatedBefore  string `url:"createdBefore,omitempty"`
	CreatedAfter   string `url:"createdAfter,omitempty"`
}

type PaymentService interface {
	Get(context.Context, ListOptions) (*PaymentsRoot, error)
	Find(context.Context, string) (*Payment, error)
	Create(context.Context, *PaymentParams) (*CreatePaymentResponse, error)

	Cancel(ctx context.Context, paymentId string) error
	Capture(ctx context.Context, paymentId string, amount int) error
}

type PaymentServiceOp struct {
	Refund RefundService
	client *Client
}

var _ PaymentService = &PaymentServiceOp{}

type Payment struct {
	PaymentId               string `json:"paymentId,omitempty"`
	Amount                  int    `json:"amount,omitempty"`
	Description             string `json:"description,omitempty"`
	PaymentPointId          string `json:"paymentPointId,omitempty"`
	Reference               string `json:"reference,omitempty"`
	MobilePayAppRedirectUri string `json:"mobilePayAppRedirectUri,omitempty"`
	State                   string `json:"state,omitempty"`
	InitiatedOn             string `json:"initiatedOn,omitempty"`
	LastUpdatedOn           string `json:"lastUpdatedOn,omitempty"`
	MerchantId              string `json:"merchantId,omitempty"`
	IsoCurrencyCode         string `json:"isoCurrencyCode,omitempty"`
	PaymentPointName        string `json:"paymentPointName,omitempty"`
}

type PaymentsRoot struct {
	Payments       []Payment `json:"payments"`
	PageSize       int       `json:"pageSize"`
	NextPageNumber int       `json:"nextPageNumber"`
}

type CreatePaymentResponse struct {
	PaymentId               string `json:"paymentId"`
	MobilePayAppRedirectUri string `json:"mobilePayAppRedirectUri"`
}

type PaymentParams struct {
	Amount         int    `json:"amount"`
	IdempotencyKey string `json:"idempotencyKey"`
	PaymentPointId string `json:"paymentPointId"`
	RedirectUri    string `json:"redirectUri"`
	Reference      string `json:"reference"`
	Description    string `json:"description"`
}

func (ps PaymentServiceOp) Get(ctx context.Context, opts ListOptions) (*PaymentsRoot, error) {
	path := paymentsBasePath

	path, err := addOptions(path, opts)
	if err != nil {
		log.Println("err", err)
		return nil, err
	}

	req, err := ps.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(PaymentsRoot)
	_, err = ps.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (ps *PaymentServiceOp) Find(ctx context.Context, paymentId string) (*Payment, error) {
	if paymentId == "" {
		ps.client.Logger.Errorf("paymentParams cannot be empty")

		return nil, newArgError("paymentId", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", paymentsBasePath, paymentId)

	req, err := ps.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(Payment)
	_, err = ps.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (ps *PaymentServiceOp) Create(ctx context.Context, paymentParams *PaymentParams) (*CreatePaymentResponse, error) {
	if paymentParams == nil {
		ps.client.Logger.Errorf("paymentParams cannot be nil %v", paymentParams)

		return nil, newArgError("paymentParams", "cannot be nil")
	}

	path := paymentsBasePath

	req, err := ps.client.NewRequest(ctx, http.MethodPost, path, paymentParams)
	if err != nil {
		return nil, err
	}

	root := new(CreatePaymentResponse)
	_, err = ps.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (ps *PaymentServiceOp) Cancel(ctx context.Context, paymentId string) error {
	if paymentId == "" {
		ps.client.Logger.Errorf("paymentParams cannot be empty")

		return newArgError("paymentId", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/cancel", paymentsBasePath, paymentId)

	req, err := ps.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	_, err = ps.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PaymentServiceOp) Capture(ctx context.Context, paymentId string, amount int) error {
	if paymentId == "" {
		ps.client.Logger.Errorf("paymentId cannot be empty", paymentId)

		return newArgError("paymentId", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/capture", paymentsBasePath, paymentId)

	type captureRequest struct {
		Amount int `json:"amount"`
	}

	requestData := &captureRequest{Amount: amount}

	req, err := ps.client.NewRequest(ctx, http.MethodPost, path, requestData)
	if err != nil {
		return err
	}

	_, err = ps.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
