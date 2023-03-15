package shared

import (
	"fmt"
	"net/http"
)

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
	return fmt.Sprintf("http error (%s) code: %d status: %s", h.URL, h.StatusCode, h.Status)
}

type DecodeError struct {
	Name string
	Err  error
}

func (d *DecodeError) Error() string {
	return fmt.Sprintf("%s decode error: %v", d.Name, d.Err)
}
