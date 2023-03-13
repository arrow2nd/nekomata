package shared

import "net/url"

// CreateJoinURL : パスを結合したURLを作成
func CreateURL(host string, v ...string) *url.URL {
	u := &url.URL{}

	u.Scheme = "https"
	u.Host = host
	u = u.JoinPath(v...)

	return u
}
