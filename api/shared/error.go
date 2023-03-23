package shared

import (
	"fmt"
	"net/http"
)

type RequestError struct {
	Endpoint Endpoint
	Err      error
}

func (r RequestError) Error() string {
	return fmt.Sprintf("request error (%s): %v", r.Endpoint, r.Err)
}

type HTTPError struct {
	Status     string
	StatusCode int
	URL        string
}

func NewHTTPError(res *http.Response) *HTTPError {
	return &HTTPError{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		URL:        res.Request.URL.String(),
	}
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("http error (%s): code %d, status %s", h.URL, h.StatusCode, h.Status)
}

type DecodeError struct {
	Endpoint Endpoint
	Err      error
}

func (d *DecodeError) Error() string {
	return fmt.Sprintf("decode error (%s): %v", d.Endpoint, d.Err)
}
