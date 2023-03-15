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

const listenAddr = "localhost:3000"

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
	authURL, sessionID, err := m.createAuthorizeURL(permissions)
	if err != nil {
		return nil, err
	}
	shared.PrintAuthenticateURL(w, authURL)

	// セッションIDを受け取る
	id, err := m.recieveSessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return m.recieveToken(id)
}

func (m *Misskey) createAuthorizeURL(permissions []string) (string, string, error) {
	sessionID, _ := uuid.NewUUID()
	callbackURL, _ := shared.CreateURL(nil, "http://"+listenAddr, "callback")

	q := &url.Values{}
	q.Add("name", m.opts.Name)
	q.Add("callback", callbackURL)
	q.Add("permission", strings.Join(permissions, ","))

	authURL, err := shared.CreateURL(q, m.opts.Server, "miauth", sessionID.String())
	if err != nil {
		return "", "", fmt.Errorf("failed to create URL: %w", err)
	}

	return authURL, sessionID.String(), nil
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
		Addr:    "localhost:3000",
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- serve.ListenAndServe()
	}()

	recievedSessionID := <-sessionID

	// サーバーを閉じる
	if err := serve.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("failed to shut down server: %w", err)
	}

	if err := <-serverErr; err != http.ErrServerClosed {
		return "", fmt.Errorf("failed to listen server: %w", err)
	}

	if recievedSessionID == "" {
		return "", fmt.Errorf("failed to recieve session id")
	}

	return recievedSessionID, nil
}

func (m *Misskey) recieveToken(sessionID string) (*shared.User, error) {
	endpointURL, err := shared.CreateURL(nil, m.opts.Server, "api", "miauth", sessionID, "check")
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	res, err := http.Post(endpointURL, "text/plain", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}

	authRes := &miAuthResponse{}
	decorder := json.NewDecoder(res.Body)
	if err := decorder.Decode(authRes); err != nil {
		return nil, fmt.Errorf("failed to decord json: %w", err)
	}

	if !authRes.OK {
		return nil, fmt.Errorf("failed to get token")
	}

	return &shared.User{
		Token: authRes.Token,
	}, nil
}
