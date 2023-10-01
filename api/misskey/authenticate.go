package misskey

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api"
	"github.com/google/uuid"
)

type miAuthResponse struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}

func (m *Misskey) Authenticate(w io.Writer) (*api.User, error) {
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
	api.PrintAuthenticateURL(w, url)

	// セッションIDを受け取る
	id, err := m.recieveSessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return m.recieveToken(id)
}

func (m *Misskey) createAuthorizeURL(permissions []string) (string, string) {
	q := url.Values{}
	q.Add("name", m.opts.Name)
	q.Add("callback", api.AuthCallbackURL)
	q.Add("permission", strings.Join(permissions, ","))

	sessionID, _ := uuid.NewUUID()
	p := url.Values{}
	p.Add(":session_id", sessionID.String())

	endpoint := endpointMiAuth.URL(m.opts.Server, p)
	return endpoint + "?" + q.Encode(), sessionID.String()
}

func (m *Misskey) recieveSessionID(id string) (string, error) {
	return api.RecieveAuthenticateCode("session", func(sessionID string) bool {
		return sessionID == id
	})
}

func (m *Misskey) recieveToken(sessionID string) (*api.User, error) {
	p := url.Values{}
	p.Add(":session_id", sessionID)

	endpoint := endpointMiAuthCheck.URL(m.opts.Server, p)
	res, err := http.Post(endpoint, "text/plain", nil)
	if err != nil {
		return nil, &api.RequestError{
			URL: endpoint,
			Err: err,
		}
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, api.NewHTTPError(res)
	}

	authRes := &miAuthResponse{}
	decorder := json.NewDecoder(res.Body)
	if err := decorder.Decode(authRes); err != nil {
		return nil, &api.DecodeError{
			URL: endpoint,
			Err: err,
		}
	}

	if !authRes.OK {
		return nil, fmt.Errorf("get token error: invalid authentication URL")
	}

	return &api.User{
		Token: authRes.Token,
	}, nil
}
