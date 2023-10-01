package misskey

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/arrow2nd/nekomata/api"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthorizeURL(t *testing.T) {
	m := &Misskey{opts: &api.ClientOpts{Name: "test_app", Server: "https://example.com"}}

	u, sessionID := m.createAuthorizeURL([]string{"aaaa", "bbbb"})
	assert.NotEqual(t, "", sessionID, "セッションIDがあるか")

	pathParams := url.Values{}
	pathParams.Add(":session_id", sessionID)
	endpoint := endpointMiAuth.URL(m.opts.Server, pathParams)

	want := endpoint + "?callback=http%3A%2F%2Flocalhost%3A3000%2Fcallback&name=test_app&permission=aaaa%2Cbbbb"
	assert.Equal(t, want, u, "正しい形式で生成されているか")
}

func TestRecieveToken(t *testing.T) {
	isExpired := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/miauth/SESSION_ID/check" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if isExpired {
			isExpired = false
			fmt.Fprintln(w, `{"ok": false}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"ok": true, "token": "USER_TOKEN"}`)
	}))

	defer ts.Close()

	t.Run("URL期限切れ", func(t *testing.T) {
		m := &Misskey{opts: &api.ClientOpts{Server: ts.URL}}
		_, err := m.recieveToken("SESSION_ID")
		assert.ErrorContains(t, err, "get token error")
	})

	t.Run("アクセストークンが取得できるか", func(t *testing.T) {
		m := &Misskey{opts: &api.ClientOpts{Server: ts.URL}}
		res, err := m.recieveToken("SESSION_ID")
		assert.NoError(t, err)
		assert.Equal(t, "USER_TOKEN", res.Token)
	})
}
