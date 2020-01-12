package appswitch

import "time"

const (
	// CaptureTypeFull is a type where the full amount must also be the amount captured.
	CaptureTypeFull CaptureType = "Full"
	// CaptureTypePartial is a type where the amount to capture can either be the exact or less than the amount reserved.
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

// CanceledReservation maps to a cancelled transaction
type CanceledReservation struct {
	TransactionID string `json:"TransactionId"`
}

// RefundParams holds the parameters used to make a refund.
// Docs https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/refund%20amount%20interface%20description.md#data-types-and-constants
type RefundParams struct {
	Amount  float64 `json:"Amount"`
	BulkRef string  `json:"BulkRef"`
}

// RefundedReservation is the type of a refunded reservation.
type RefundedReservation struct {
	TransactionID         string  `json:"TransactionId"`         // The transaction id of the new refund payment.
	OriginalTransactionID string  `json:"OriginalTransactionId"` // The id of the transaction which is now being (partially) refunded.
	Remainder             float64 `json:"Remainder"`             // The remaining amount. If the transaction is completely refunded then this is 0.00.
}

// CaptureParams holds the parameters used to make a capture.
// See https://github.com/MobilePayDev/MobilePay-AppSwitch-API/blob/master/REST%20APIs/v1/capture%20amount%20interface%20description.md#data-types-and-constants
// Amount:
// Partial: Amount must be > 0.0 and <= reserved amount
// If Amount is less than the reserved amount, then that amount is captured. The remaining amount is then released, i.e. not reserved any more.
// Full: If the reservation type is Full then only the full amount can be captured.
type CaptureParams struct {
	Amount  float64 `json:"Amount"`
	BulkRef string  `json:"BulkRef"` // A reference for bulking payments on the merchants account statement.
}

// CapturedReservation is the type of a captured reservation.
type CapturedReservation struct {
	TransactionID string `json:"TransactionId"` // This is the transaction ID of the payment transaction. It is the same transaction id returned when reserving the amount.
}
