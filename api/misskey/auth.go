package misskey

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/google/uuid"
)

func (m *Misskey) Authenticate(w io.Writer) (*shared.AuthResponse, error) {
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

	miauth := &miAuth{
		Name:        m.config.AppName,
		Host:        m.config.Host,
		Permissions: permissions,
	}

	return miauth.Run(w)
}

const listenAddr = "localhost:3000"

type miAuthResponse struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}

type miAuth struct {
	Name        string
	Host        string
	Permissions []string

	sessionID string
	res       chan string
}

func (m *miAuth) Run(w io.Writer) (*shared.AuthResponse, error) {
	m.res = make(chan string, 1)

	// 認証URLを組み立て
	authURL := m.createAuthURL(m.Permissions)
	shared.PrintAuthURL(w, authURL)

	// サーバーを建ててリダイレクトを待機
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", m.handleCallback)

	serve := http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		if err := serve.ListenAndServe(); err != http.ErrServerClosed {
			// TODO: ここのエラーも返したい
			log.Fatalf("server error: %v", err)
		}
	}()

	sessionID := <-m.res

	// サーバーを閉じる
	if err := serve.Shutdown(context.Background()); err != nil {
		return nil, err
	}

	return m.recieveToken(sessionID)
}

func (m *miAuth) createAuthURL(permissions []string) string {
	ID, _ := uuid.NewUUID()
	m.sessionID = ID.String()
	u := shared.CreateURL(m.Host, "miauth", m.sessionID)

	q := url.Values{}
	q.Add("name", m.Name)
	q.Add("callback", "http://"+listenAddr+"/callback")
	q.Add("permission", strings.Join(permissions, ","))

	u.RawQuery = q.Encode()
	return u.String()
}

func (m *miAuth) handleCallback(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")

	// 受け取ったセッションIDが正しいかチェック
	if m.sessionID != sessionID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m.res <- sessionID
	w.Write([]byte("Authentication complete! You may close this page."))
}

func (m *miAuth) recieveToken(sessionID string) (*shared.AuthResponse, error) {
	endpointURL := shared.CreateURL(m.Host, "api", "miauth", sessionID, "check")
	res, err := http.Post(endpointURL.String(), "text/plain", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: status %s", res.Status)
	}

	decorder := json.NewDecoder(res.Body)
	authRes := &miAuthResponse{}
	if err := decorder.Decode(authRes); err != nil {
		return nil, err
	}

	if !authRes.OK {
		return nil, fmt.Errorf("failed to get token")
	}

	return &shared.AuthResponse{
		Token: authRes.Token,
	}, nil
}
