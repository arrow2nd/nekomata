package mastodon

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/shared"
)

type authenticateResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	CreatedAt   int64  `json:"created_at"`
}

func (m *Mastodon) Authenticate(w io.Writer) (*shared.User, error) {
	permissions := []string{"read", "write", "follow"}

	// 認証URL組み立て
	url := m.createAuthorizeURL(permissions)
	shared.PrintAuthenticateURL(w, url)

	// 認証コードを受け取る
	code, err := m.recieveCode()
	if err != nil {
		return nil, err
	}

	return m.recieveToken(code)
}

func (m *Mastodon) createAuthorizeURL(permissions []string) string {
	q := url.Values{}

	q.Add("response_type", "code")
	q.Add("client_id", m.opts.ID)
	q.Add("redirect_uri", shared.AuthCallbackURL)
	q.Add("scope", strings.Join(permissions, " "))

	endpoint := endpointOauthAuthorize.URL(m.opts.Server, nil)
	return endpoint + "?" + q.Encode()
}

func (m *Mastodon) recieveCode() (string, error) {
	return shared.RecieveAuthenticateCode("code", func(code string) bool {
		return code != ""
	})
}

func (m *Mastodon) recieveToken(code string) (*shared.User, error) {
	q := url.Values{}

	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("client_id", m.opts.ID)
	q.Add("client_secret", m.opts.Secret)
	q.Add("redirect_uri", shared.AuthCallbackURL)

	opts := &requestOpts{
		method: http.MethodPost,
		url:    endpointOauthToken.URL(m.opts.Server, nil),
		q:      q,
		isAuth: false,
	}

	res := authenticateResponse{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return &shared.User{Token: res.AccessToken}, nil
}
