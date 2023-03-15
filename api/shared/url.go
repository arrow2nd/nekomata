package shared

import (
	"fmt"
	"net/url"
)

// CreateURL : URLを作成
func CreateURL(q *url.Values, rawURL string, path ...string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to create URL: %w", err)
	}

	if len(path) != 0 {
		u = u.JoinPath(path...)
	}

	if q != nil {
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
