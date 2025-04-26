package apierrors

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Response represents a JSON:API error response body.
type Response struct {
	RequestID string   `json:"requestID,omitempty"`
	Errors    []*Error `json:"errors,omitempty"`
}

// Add appends additional errors to an API error response
func (resp *Response) Add(errs ...error) {
	for _, err := range errs {
		switch v := err.(type) {
		case *Response:
			if resp.RequestID == "" {
				resp.RequestID = v.RequestID
			}
			resp.Errors = append(resp.Errors, v.Errors...)
		default:
			resp.Errors = append(resp.Errors, FromError(err))
		}
	}
}

// HasErrors returns if there are any errors in the response
func (resp *Response) HasErrors() bool {
	return len(resp.Errors) > 0
}

// StatusCode returns the (highest) status code from the slice of errors in a response
func (resp *Response) StatusCode() int {
	var highestStatusCode int
	for _, e := range resp.Errors {
		if e.Status > highestStatusCode {
			highestStatusCode = e.Status
		}
	}

	return highestStatusCode
}

// MarshalJSON handles marshaling the error into JSON, per the JSON-marshaling interface
func (resp *Response) MarshalJSON() ([]byte, error) {
	marshalResponse := *resp
	return json.Marshal(marshalResponse)
}

// Error converts all errors contained within a response into a comma-separated string
func (resp *Response) Error() string {
	errStrings := make([]string, 0, len(resp.Errors))
	for i := range resp.Errors {
		errStrings = append(errStrings, resp.Errors[i].Error())
	}

	return strings.Join(errStrings, ",")
}

// Write will write and errors status and JSON marshalled bytes to the http response
func (resp *Response) Write(w http.ResponseWriter) error {
	w.WriteHeader(resp.StatusCode())
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// NewResponse creates an API error response from a slice of errors
func NewResponse(errs ...error) *Response {
	// Prepare the response
	resp := &Response{}

	// Add all provided errors
	resp.Add(errs...)

	// Return the final response
	return resp
}
