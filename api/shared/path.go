package shared

import "net/url"

// CreateJoinURL : パスを結合したURLを作成
func CreateURL(scheme, host string, v ...string) *url.URL {
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
	}

	return u.JoinPath(v...)
}
