package appswitch

import (
	"fmt"
	"github.com/steffen25/mobilepay-go"
	"net/http"
	"time"
)

// UrlQueryTimestampLayout represents a date format used to encode dates as query parameters for the AppSwitch API
const UrlQueryTimestampLayout = "2006-01-02T15_04"

// Client is used to invoke APIs related to AppSwitch.
type Client struct {
	Backend    mobilepay.Backend
	MerchantID string
}

// GetPaymentStatus retrieves the status of a payment
func (c Client) GetPaymentStatus(orderID string) (*PaymentStatus, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s", c.MerchantID, orderID)
	status := &PaymentStatus{}
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, nil, nil, status)

	return status, err
}

// GetTransactions retrieves the transactions related to a specific payment
func (c Client) GetTransactions(orderID string) ([]*PaymentTransaction, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s/transactions", c.MerchantID, orderID)
	var statuses []*PaymentTransaction
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, nil, nil, &statuses)

	return statuses, err
}

// GetReservations returns a list of reservations created by a merchant within a given time frame
func (c Client) GetReservations(params *GetReservationsParams) ([]*Reservation, error) {
	url := fmt.Sprintf("/reservations/merchants/%s/%s/%s", c.MerchantID, dateFormatter(params.From), dateFormatter(params.To))
	var reservations []*Reservation
	err := c.Backend.Call(http.MethodGet, url, c.MerchantID, params, nil, &reservations)

	return reservations, err
}

// CancelReservation cancels a previously submitted reservation made by a merchant
func (c Client) CancelReservation(orderID string) (*CanceledReservation, error) {
	url := fmt.Sprintf("/reservations/merchants/%s/orders/%s", c.MerchantID, orderID)
	canceledReservation := &CanceledReservation{}
	err := c.Backend.Call(http.MethodDelete, url, c.MerchantID, nil, nil, canceledReservation)

	return canceledReservation, err
}

// Refund refunds the transaction amount, either the entire amount or just a part of the amount.
// It is possible to refund transactions within a year after capture.
func (c Client) Refund(orderID string, params *RefundParams) (*RefundedReservation, error) {
	url := fmt.Sprintf("/merchants/%s/orders/%s", c.MerchantID, orderID)
	refundedReservation := &RefundedReservation{}
	err := c.Backend.Call(http.MethodPut, url, c.MerchantID, nil, params, refundedReservation)

	return refundedReservation, err
}

// Capture captures the transaction, i.e. carries out the actual payment.
// It is important to know which capture type to use, since it must match the reservation type.
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
