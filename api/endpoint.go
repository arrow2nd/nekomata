package api

import (
	"net/url"
	"strings"
)

type Endpoint string

func (e Endpoint) URL(host string, pathParams url.Values) string {
	endpoint := string(e)

	// パスパラメータを設定
	if pathParams != nil {
		for k, v := range pathParams {
			endpoint = strings.Replace(endpoint, k, v[0], 1)
		}
	}

	return host + endpoint
}
