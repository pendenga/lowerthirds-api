package apierrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/xid"
)

// ErrorSource represents the source parameter of errors.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// Error defines the structure for a JSON-formatted API error
type Error struct {
	ID         string                  `json:"id,omitempty"`
	Status     int                     `json:"status,omitempty"`
	Code       string                  `json:"code,omitempty"`
	Title      string                  `json:"title,omitempty"`
	Detail     string                  `json:"detail,omitempty"`
	Source     *ErrorSource            `json:"source,omitempty"`
	Meta       *map[string]interface{} `json:"meta,omitempty"`
	innerError error
}

// Error handles printing the error as a string, per the error interface.
func (err *Error) Error() string {
	return fmt.Sprintf("%s - %s (%s) - %s\n\t%s", err.ID, err.Code, http.StatusText(err.Status), err.Title, err.Detail)
}

// MarshalJSON handles marshaling the error into JSON, per the JSON-marshaling interface
func (err *Error) MarshalJSON() ([]byte, error) {
	marshalErr := *err
	return json.Marshal(marshalErr)
}

// Unwrap enables the unwrapping of errors, which can be used with the `Printf` "%w" verb
func (err *Error) Unwrap() error {
	return err.innerError
}

// WithSource sets a source pointer and parameter of an error
func (err *Error) WithSource(pointer string, parameter string) *Error {
	if err.Source == nil {
		err.Source = &ErrorSource{}
	}

	err.Source.Pointer = pointer
	err.Source.Parameter = parameter

	return err
}

// Write will write and errors status and JSON marshalled bytes to the http response
func (err *Error) Write(w http.ResponseWriter) error {
	w.WriteHeader(err.Status)
	b, e := json.Marshal(err)
	if e != nil {
		return e
	}
	_, e = w.Write(b)
	return e
}

// New creates a new `apierrors.Error`
func New(statusCode int, code string, title string, detail string, args ...interface{}) *Error {
	// Use fmt.Sprintf only if there are arguments
	if len(args) > 0 {
		detail = fmt.Sprintf(detail, args...)
	}

	return &Error{
		ID:     xid.New().String(),
		Status: statusCode,
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}

// FromError creates a new `apierrors.Error` from a standard error.
func FromError(originalErr error) *Error {
	var apierror *Error
	switch {
	case errors.As(originalErr, &apierror):
		return apierror
	case errors.Is(originalErr, context.Canceled):
		return New(499, "Client Closed Request",
			"The client has closed the request before the server could send a response.", originalErr.Error())
	default:
		err := New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error", originalErr.Error())
		err.innerError = originalErr
		return err
	}
}
