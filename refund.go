package mobilepay

import (
	"context"
	"net/http"
)

const refundsBasePath = "v1/refunds"

type RefundService interface {
	List(ctx context.Context, opt *RefundsListOptions) (*RefundsRoot, error)
	Create(ctx context.Context, createRequest *RefundParams) (*Refund, error)
}

type RefundParams struct {
	IdempotencyKey string `json:"idempotencyKey"`
	PaymentId      string `json:"paymentId"`
	Amount         int    `json:"amount"`
	Reference      string `json:"reference"`
	Description    string `json:"description"`
}

type Refund struct {
	RefundId        string `json:"refundId"`
	PaymentId       string `json:"paymentId"`
	Amount          int    `json:"amount"`
	RemainingAmount int    `json:"remainingAmount,omitempty"`
	Description     string `json:"description"`
	Reference       string `json:"reference"`
	CreatedOn       string `json:"createdOn"`
}

type RefundsRoot struct {
	Refunds        []Refund `json:"refunds"`
	PageSize       int      `json:"pageSize"`
	NextPageNumber int      `json:"nextPageNumber"`
}

type RefundServiceOp struct {
	client *Client
}

var _ RefundService = &RefundServiceOp{}

func (rs *RefundServiceOp) Create(ctx context.Context, refundParams *RefundParams) (*Refund, error) {
	if refundParams == nil {
		rs.client.Logger.Errorf("refundParams cannot be nil")

		return nil, NewArgError("refundParams", "cannot be nil")
	}

	path := refundsBasePath

	req, err := rs.client.NewRequest(ctx, http.MethodPost, path, refundParams)
	if err != nil {
		return nil, err
	}

	root := new(Refund)
	_, err = rs.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}

func (rs RefundServiceOp) List(ctx context.Context, opts *RefundsListOptions) (*RefundsRoot, error) {
	path := refundsBasePath

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}

	req, err := rs.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(RefundsRoot)
	_, err = rs.client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root, err
}
