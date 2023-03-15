package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
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
	q := &url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", m.opts.ID)
	q.Add("redirect_uri", shared.AuthCallbackURL)
	q.Add("scope", strings.Join(permissions, " "))

	return oauthAuthorizeEndpoint.URL(m.opts.Server) + "?" + q.Encode()
}

func (m *Mastodon) recieveCode() (string, error) {
	code := make(chan string, 1)
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		recieved := r.URL.Query().Get("code")
		if recieved != "" {
			code <- recieved
			w.Write([]byte("Authentication complete! You may close this page."))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		code <- ""
	})

	// サーバーを建ててリダイレクトを待機
	serve := http.Server{
		Addr:    shared.AuthCallbackAddr,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- serve.ListenAndServe()
	}()

	recievedCode := <-code

	if err := serve.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("server shutdown error: %w", err)
	}

	if err := <-serverErr; err != http.ErrServerClosed {
		return "", fmt.Errorf("listen server error: %w", err)
	}

	if recievedCode == "" {
		return "", fmt.Errorf("failed to recieve authenticate code")
	}

	return recievedCode, nil
}

func (m *Mastodon) recieveToken(code string) (*shared.User, error) {
	q := url.Values{}
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("client_id", m.opts.ID)
	q.Add("client_secret", m.opts.Secret)
	q.Add("redirect_uri", shared.AuthCallbackURL)

	endpoint := oauthTokenEndpoint.URL(m.opts.Server)
	res, err := http.PostForm(endpoint, q)
	if err != nil {
		return nil, &shared.RequestError{
			Endpoint: oauthTokenEndpoint,
			Err:      err,
		}
	}

	defer res.Body.Close()
	decorder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		e := &errorResponse{}
		if err := decorder.Decode(e); err != nil {
			return nil, shared.NewHTTPError(res)
		}

		return nil, e
	}

	authRes := &authenticateResponse{}
	if err := decorder.Decode(authRes); err != nil {
		return nil, &shared.DecodeError{
			Endpoint: oauthTokenEndpoint,
			Err:      err,
		}
	}

	return &shared.User{
		Token: authRes.AccessToken,
	}, nil
}
