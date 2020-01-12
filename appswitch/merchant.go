package appswitch

import (
	"strings"
	"time"
)

const (
	// PaymentTypeReserved represents a type for a payment that is reserved
	PaymentTypeReserved PaymentStatusType = "Reserved"
	// PaymentTypeCancelled represents a type for a payment that is cancelled by a merchant
	PaymentTypeCancelled PaymentStatusType = "Cancelled"
	// PaymentTypeCaptured represents a type for a payment that is captured by a merchant
	PaymentTypeCaptured PaymentStatusType = "Captured"
	// PaymentTypeTotalRefund represents a type for a payment that have been fully refunded by a merchant
	PaymentTypeTotalRefund PaymentStatusType = "TotalRefund"
	// PaymentTypePartialRefund represents a type for a payment that have been partially refunded by a merchant
	PaymentTypePartialRefund PaymentStatusType = "PartialRefund"
	// PaymentTypeRejected represents a type for a payment where the reservation, capture, refund or cancellation is rejected
	PaymentTypeRejected PaymentStatusType = "Rejected"
)

// TimestampLayout is a date format used by the AppSwitch API. It returns timestamps in ISO8601 (UTC time) format.
const TimestampLayout = "2006-01-02T15:04:05.000"

// PaymentStatusType is a type to match the different payment status types
type PaymentStatusType string

// PaymentStatus represents the result from the /api/v1/merchants/{merchantId}/orders/{orderId} endpoint
type PaymentStatus struct {
	LatestPaymentStatus string  `json:"LatestPaymentStatus"`
	TransactionID       string  `json:"TransactionId"`
	OriginalAmount      float64 `json:"OriginalAmount"`
}

// PaymentTransaction is the type of a payment status
type PaymentTransaction struct {
	TimeStamp     MobilePayTimestamp `json:"TimeStamp"`
	PaymentStatus string             `json:"PaymentStatus"`
	TransactionID string             `json:"TransactionId"`
	Amount        float64            `json:"Amount"`
}

// MobilePayTimestamp is a time.Time wrapper in order to parse timestamps returned by the AppSwitch API.
// Since Go expects a RFC3339 format we use a custom time type to parse the format returned from MobilePay
type MobilePayTimestamp struct {
	time.Time
}

// UnmarshalJSON handles deserialization of a AppSwitch API timestamp.
func (mpTime *MobilePayTimestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(TimestampLayout, s)
	if err != nil {
		return err
	}
	mpTime.Time = t
	return nil
}
