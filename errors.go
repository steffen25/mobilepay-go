package mobilepay

import (
	"fmt"
	"time"
)

type AuthError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type BadRequestError struct {
	Reason string `json:"Reason"`
}

type RateLimitError struct {
	// TODO figure out if this has any headers like retry after e.g
	RetryAfter time.Duration
}

type ServerError struct {
	CorrelationID string `json:"CorrelationId"`
	Errortype     string `json:"Errortype"`
	Message       string `json:"Message"`
}

type ResponseDecodingError struct {
	Body    []byte
	Message string
	Status  int
}

type ResponseError struct {
	Status  int
	Message string
}

func (e *AuthError) Empty() bool {
	return e.StatusCode == 0 && e.Message == ""
}

func (e *BadRequestError) Empty() bool {
	return e.Reason == ""
}

func (e *ServerError) Empty() bool {
	return e.CorrelationID == "" && e.Message == "" && e.Errortype == ""
}

func (e AuthError) Error() string {
	return e.Message
}

func (e BadRequestError) Error() string {
	return e.Reason
}

func (e *RateLimitError) Error() string {
	return fmt.Sprint("mobilepay rate limit exceeded")
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("mobilepay server error detected. Message: %s Type: %s CorrelationId %s", e.Message, e.Errortype, e.CorrelationID)
}

func (e ResponseDecodingError) Error() string {
	return e.Message
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("Unknown error. Status %d Body: %s", e.Status, e.Message)
}
