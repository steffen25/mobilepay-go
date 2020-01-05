package appswitch

import (
	"strings"
	"time"
)

const (
	PaymentTypeReserved      PaymentStatusType = "Reserved"
	PaymentTypeCancelled     PaymentStatusType = "Cancelled"
	PaymentTypeCaptured      PaymentStatusType = "Captured"
	PaymentTypeTotalRefund   PaymentStatusType = "TotalRefund"
	PaymentTypePartialRefund PaymentStatusType = "PartialRefund"
	PaymentTypeRejected      PaymentStatusType = "Rejected"
)

// AppSwitch API returns timestamps in ISO8601 (UTC time) format.
const TimestampLayout = "2006-01-02T15:04:05.000"

// PaymentStatusType is a type to match the different payment status types
type PaymentStatusType string

// PaymentStatus represents the result from the /api/v1/merchants/{merchantId}/orders/{orderId} endpoint
type PaymentStatus struct {
	LatestPaymentStatus string  `json:"LatestPaymentStatus"`
	TransactionID       string  `json:"TransactionId"`
	OriginalAmount      float64 `json:"OriginalAmount"`
}

// PaymentStatus represents the result from the /api/v1/merchants/{merchantId}/orders/{orderId}/transactions endpoint
type PaymentTransaction struct {
	TimeStamp     MobilePayTimestamp `json:"TimeStamp"`
	PaymentStatus string             `json:"PaymentStatus"`
	TransactionID string             `json:"TransactionId"`
	Amount        float64            `json:"Amount"`
}

// Since Go expects a RFC3339 format we use a custom time type to parse the format returned from MobilePay
type MobilePayTimestamp struct {
	time.Time
}

func (mpTime *MobilePayTimestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(TimestampLayout, s)
	if err != nil {
		return err
	}
	mpTime.Time = t
	return nil
}
