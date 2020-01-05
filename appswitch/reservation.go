package appswitch

import "time"

const (
	CaptureTypeFull    CaptureType = "Full"
	CaptureTypePartial CaptureType = "Partial"
)

// CaptureType is a type to match the different capture types
type CaptureType string

// Reservation represents the result from the /api/v1/reservations/merchants/{merchantId}/{datetimeFrom}/{datetimeTo}?customerId={customerId} endpoint
// Docs https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/get%20reservations%20interface%20description.md#data-types-and-constants
type Reservation struct {
	TimeStamp     string  `json:"TimeStamp"`
	OrderID       string  `json:"OrderId"`
	TransactionID string  `json:"TransactionId"`
	Amount        float64 `json:"Amount"`
	CaptureType   string  `json:"CaptureType"`
}

// GetReservationsParams is a type used to query the get reservations endpoint
// Docs https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/get%20reservations%20interface%20description.md#data-types-and-constants
type GetReservationsParams struct {
	From       time.Time `url:"-"`
	To         time.Time `url:"-"`
	CustomerID string    `url:"customerId,omitempty"`
}

type CanceledReservation struct {
	TransactionID string `json:"TransactionId"`
}

// Docs https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/refund%20amount%20interface%20description.md#data-types-and-constants
type RefundParams struct {
	Amount  float64 `json:"Amount"`
	BulkRef string  `json:"BulkRef"`
}

type RefundedReservation struct {
	TransactionID         string  `json:"TransactionId"`
	OriginalTransactionID string  `json:"OriginalTransactionId"`
	Remainder             float64 `json:"Remainder"`
}

// See https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/capture%20amount%20interface%20description.md#data-types-and-constants
type CaptureParams struct {
	Amount  float64 `json:"Amount"`
	BulkRef string  `json:"BulkRef"`
}

type CapturedReservation struct {
	TransactionID string `json:"TransactionId"`
}
