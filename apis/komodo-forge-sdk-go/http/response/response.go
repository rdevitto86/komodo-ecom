package httpresponse

import (
	"encoding/json"
	"fmt"
	httpErr "komodo-forge-sdk-go/http/errors"
	"net/http"
)

type APIResponse struct {
	Status  	int
	Body    	[]byte // raw response body
	Headers 	http.Header
	RequestID string
	Error 		*httpErr.ErrorCode
}

// Unmarshals the response body into the target struct
func Bind(res *http.Response, target any) (*APIResponse, error) {
	if res == nil {
		return nil, fmt.Errorf("failed to bind response - response is nil")
	}

	body, err := json.Marshal(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to bind response - failed to marshal response body: %v", err)
	}

	return &APIResponse{
		Status:   	res.StatusCode,
		Body:     	body,
		Headers:  	res.Header,
		RequestID: 	res.Header.Get("X-Request-ID"),
		Error:    	nil,
	}, nil
}

// ResponseWriter wraps http.ResponseWriter to capture status code and bytes written.
type ResponseWriter struct {
	http.ResponseWriter
	Status       int
	BytesWritten int
	WroteHeader  bool
}

func (wtr *ResponseWriter) WriteHeader(code int) {
	if !wtr.WroteHeader {
		wtr.Status = code
		wtr.WroteHeader = true
	}
	wtr.ResponseWriter.WriteHeader(code)
}

func (wtr *ResponseWriter) Write(b []byte) (int, error) {
	if !wtr.WroteHeader { wtr.WriteHeader(http.StatusOK) }
	num, err := wtr.ResponseWriter.Write(b)
	wtr.BytesWritten += num
	return num, err
}

func (wtr *ResponseWriter) Unwrap() http.ResponseWriter {
	return wtr.ResponseWriter
}

func IsSuccess(status int) bool { return status >= 200 && status < 300 }
func IsError(status int) bool { return status >= 400 && status < 600 }
func IsRedirect(status int) bool { return status >= 300 && status < 400 }
func IsInformational(status int) bool { return status >= 100 && status < 200 }
