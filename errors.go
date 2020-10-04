package mobilepay

import (
	"encoding/json"
	"fmt"
	"time"
)

// AuthError represents a type for handling auth errors e.g. missing subscription key
type AuthError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// BadRequestError represents a type for handling 4xx errors
type BadRequestError struct {
	Reason string `json:"Reason"`
}

// RateLimitError represents a type for handling 429 too many requests
type RateLimitError struct {
	// TODO figure out if this has any headers like retry after e.g
	RetryAfter time.Duration
}

// ServerError represents a type for handling server errors
type ServerError struct {
	CorrelationID string `json:"CorrelationId"`
	Errortype     string `json:"Errortype"`
	Message       string `json:"Message"`
}

// ResponseDecodingError represents a type for handling json decoding errors
type ResponseDecodingError struct {
	Body    []byte
	Message string
	Status  int
}

// ResponseError represents a type of the least known error message returned from the api
type ResponseError struct {
	Status  int
	Message string
}

// A helper method to check if the error is empty.
// This means if the error struct represents its initial values
func (e *AuthError) Empty() bool {
	return e.StatusCode == 0 && e.Message == ""
}

// A helper method to check if the error is empty.
// This means if the error struct represents its initial values
func (e *BadRequestError) Empty() bool {
	return e.Reason == ""
}

// A helper method to check if the error is empty.
// This means if the error struct represents its initial values
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
	return "MobilePay rate limit exceeded"
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

func HandleAuthError(body []byte, statusCode int) error {
	authError := &AuthError{}
	err := json.Unmarshal(body, &authError)
	if err != nil {
		return NewResponseDecodingError(body, err, statusCode)
	}

	if authError.Empty() {
		responseErr := &ResponseError{
			Status:  statusCode,
			Message: string(body),
		}
		return responseErr
	}

	return authError
}

func HandleRequestError(body []byte, statusCode int) error {
	badRequestError := &BadRequestError{}
	err := json.Unmarshal(body, &badRequestError)
	if err != nil {
		return NewResponseDecodingError(body, err, statusCode)
	}

	if badRequestError.Empty() {
		responseErr := &ResponseError{
			Status:  statusCode,
			Message: string(body),
		}
		return responseErr
	}

	return badRequestError
}

func HandleServerError(body []byte, statusCode int) error {
	// somehow a status 500 can also be a BadRequestError
	// {"Reason":"BackendError"} status: 500
	serverErr := &ServerError{}
	err := json.Unmarshal(body, &serverErr)
	if err != nil {
		return NewResponseDecodingError(body, err, statusCode)
	}

	if serverErr.Empty() {
		// try bad request
		badRequestError := &BadRequestError{}
		err = json.Unmarshal(body, &badRequestError)
		if err != nil {
			return NewResponseDecodingError(body, err, statusCode)
		}

		if badRequestError.Empty() {
			responseErr := &ResponseError{
				Status:  statusCode,
				Message: string(body),
			}
			return responseErr
		}

		return badRequestError
	}

	return serverErr
}

func NewResponseDecodingError(body []byte, err error, statusCode int) error {
	return ResponseDecodingError{
		Body:    body,
		Message: err.Error(),
		Status:  statusCode,
	}
}
