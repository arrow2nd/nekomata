package mastodon

import "fmt"

type errorResponse struct {
	Type        string `json:"error"`
	Description string `json:"error_description"`
}

func (h *errorResponse) Error() string {
	return fmt.Sprintf("http error: %s (type: %s)", h.Description, h.Type)
}
