package mastodon

import (
	"io"
	"net/http"
	"net/url"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

var scope = "read write follow"

type authenticateResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	CreatedAt   int64  `json:"created_at"`
}

func (m *Mastodon) Authenticate(w io.Writer) (string, error) {
	// 認証URL組み立て
	url := m.createAuthorizeURL(scope)
	sharedapi.PrintAuthenticateURL(w, url)

	// 認証コードを受け取る
	code, err := m.recieveCode()
	if err != nil {
		return "", err
	}

	// アクセストークンを受け取る
	return m.recieveToken(code)
}

func (m *Mastodon) createAuthorizeURL(permissions string) string {
	q := url.Values{}

	q.Add("response_type", "code")
	q.Add("client_id", m.client.ID)
	q.Add("redirect_uri", sharedapi.AuthCallbackURL)
	q.Add("scope", permissions)

	endpoint := endpointOauthAuthorize.URL(m.user.Server, nil)
	return endpoint + "?" + q.Encode()
}

func (m *Mastodon) recieveCode() (string, error) {
	return sharedapi.RecieveAuthenticateCode("code", func(code string) bool {
		return code != ""
	})
}

func (m *Mastodon) recieveToken(code string) (string, error) {
	q := url.Values{}

	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", sharedapi.AuthCallbackURL)
	q.Add("client_id", m.client.ID)
	q.Add("client_secret", m.client.Secret)

	opts := &requestOpts{
		method: http.MethodPost,
		url:    endpointOauthToken.URL(m.user.Server, nil),
		q:      q,
		isAuth: false,
	}

	res := authenticateResponse{}
	if err := m.request(opts, &res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

type applicationResponse struct {
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (m *Mastodon) RegisterNewApplication() (string, string, error) {
	q := url.Values{}

	q.Add("client_name", m.client.Name)
	q.Add("redirect_uris", sharedapi.AuthCallbackURL)
	q.Add("scopes", scope)

	opts := &requestOpts{
		method: http.MethodPost,
		url:    endpointApps.URL(m.user.Server, nil),
		q:      q,
		isAuth: false,
	}

	res := applicationResponse{}
	if err := m.request(opts, &res); err != nil {
		return "", "", err
	}

	return res.ClientID, res.ClientSecret, nil
}
