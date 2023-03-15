package misskey

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/google/uuid"
)

type miAuthResponse struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}

func (m *Misskey) Authenticate(w io.Writer) (*shared.User, error) {
	permissions := []string{
		"read:account",
		"read:blocks",
		"write:blocks",
		"write:drive",
		"read:favorites",
		"read:following",
		"write:favorites",
		"write:following",
		"write:mutes",
		"write:notes",
		"read:notifications",
		"write:notifications",
		"write:reactions",
		"write:votes",
	}

	// 認証URL組み立て
	url, sessionID := m.createAuthorizeURL(permissions)
	shared.PrintAuthenticateURL(w, url)

	// セッションIDを受け取る
	id, err := m.recieveSessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return m.recieveToken(id)
}

func (m *Misskey) createAuthorizeURL(permissions []string) (string, string) {
	q := &url.Values{}
	q.Add("name", m.opts.Name)
	q.Add("callback", shared.AuthCallbackURL)
	q.Add("permission", strings.Join(permissions, ","))

	sessionID, _ := uuid.NewUUID()
	u := miAuthEndpoint.URL(m.opts.Server)
	u = strings.Replace(u, ":session_id", sessionID.String(), 1)

	return u + "?" + q.Encode(), sessionID.String()
}

func (m *Misskey) recieveSessionID(id string) (string, error) {
	mux := http.NewServeMux()

	sessionID := make(chan string, 1)
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		recieved := r.URL.Query().Get("session")
		if recieved == id {
			sessionID <- recieved
			w.Write([]byte("Authentication complete! You may close this page."))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		sessionID <- ""
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

	recievedSessionID := <-sessionID

	// サーバーを閉じる
	if err := serve.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("shutdown server error: %w", err)
	}

	if err := <-serverErr; err != http.ErrServerClosed {
		return "", fmt.Errorf("listen server error: %w", err)
	}

	if recievedSessionID == "" {
		return "", fmt.Errorf("failed to recieve session id")
	}

	return recievedSessionID, nil
}

func (m *Misskey) recieveToken(sessionID string) (*shared.User, error) {
	url := miAuthCheckEndpoint.URL(m.opts.Server)
	url = strings.Replace(url, ":session_id", sessionID, 1)

	res, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return nil, &shared.RequestError{
			Endpoint: miAuthCheckEndpoint,
			Err:      err,
		}
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, shared.NewHTTPError(res)
	}

	authRes := &miAuthResponse{}
	decorder := json.NewDecoder(res.Body)
	if err := decorder.Decode(authRes); err != nil {
		return nil, &shared.DecodeError{
			Endpoint: miAuthCheckEndpoint,
			Err:      err,
		}
	}

	if !authRes.OK {
		return nil, fmt.Errorf("get token error: invalid authentication URL")
	}

	return &shared.User{
		Token: authRes.Token,
	}, nil
}
