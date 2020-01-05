package appswitch

import (
	"fmt"
	"github.com/steffen25/mobilepay-go"
	"net/http"
	"time"
)

const UrlQueryTimestampLayout = "2006-01-02T15_04"

type Client struct {
	Backend    mobilepay.Backend
	MerchantID string
}

// /api/v1/merchants/{merchantId}/orders/{orderId}
func (c Client) GetPaymentStatus(orderID string) (*PaymentStatus, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s", c.MerchantID, orderID)
	status := &PaymentStatus{}
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, nil, nil, status)

	return status, err
}

// /api/v1/merchants/{merchantId}/orders/{orderId}/transaction
func (c Client) GetTransactions(orderID string) ([]*PaymentTransaction, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s/transactions", c.MerchantID, orderID)
	var statuses []*PaymentTransaction
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, nil, nil, &statuses)

	return statuses, err
}

// /api/v1/reservations/merchants/{merchantId}/{datetimeFrom}/{datetimeTo}?customerId={customerId}
func (c Client) GetReservations(params *GetReservationsParams) ([]*Reservation, error) {
	url := fmt.Sprintf("/reservations/merchants/%s/%s/%s", c.MerchantID, dateFormatter(params.From), dateFormatter(params.To))
	var reservations []*Reservation
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, params, nil, &reservations)

	return reservations, err
}

// /api/v1/reservations/merchants/{merchantId}/orders/{orderId}
func (c Client) CancelReservation(orderID string) (*CanceledReservation, error) {
	url := fmt.Sprintf("/reservations/merchants/%s/orders/%s", c.MerchantID, orderID)
	canceledReservation := &CanceledReservation{}
	err := c.Backend.Call(http.MethodDelete, url, c.MerchantID, nil, nil, canceledReservation)

	return canceledReservation, err
}

// /api/v1/merchants/{merchantId}/orders/{orderId}
func (c Client) Refund(orderID string, params *RefundParams) (*RefundedReservation, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s", c.MerchantID, orderID)
	refundedReservation := &RefundedReservation{}
	err := c.Backend.Call(http.MethodPut, url, c.MerchantID, nil, params, refundedReservation)

	return refundedReservation, err
}

// /api/v1/reservations/merchants/{merchantId}/orders/{orderId}
func (c Client) Capture(orderID string, params *CaptureParams) (*CapturedReservation, error) {
	url := fmt.Sprintf("/reservations/merchants/%s/orders/%s", c.MerchantID, orderID)
	capturedReservation := &CapturedReservation{}
	err := c.Backend.Call(http.MethodPut, url, c.MerchantID, nil, params, capturedReservation)

	return capturedReservation, err
}

// dateFormatter is a small utility function that formats a date into a valid date query string used by AppSwitch
func dateFormatter(date time.Time) string {
	return date.Format(UrlQueryTimestampLayout)
}
