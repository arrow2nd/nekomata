package misskey

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthorizeURL(t *testing.T) {
	m := &Misskey{
		opts: &shared.ClientOpts{Name: "test_app", Server: "https://example.com"},
	}
	permissions := []string{"aaaa", "bbbb"}

	t.Run("正常", func(t *testing.T) {
		u, sessionID, err := m.createAuthorizeURL(permissions)
		r := regexp.MustCompile("https://example.com/miauth/.+callback=http%3A%2F%2Flocalhost%3A3000%2Fcallback&name=test_app&permission=aaaa%2Cbbbb")

		assert.NoError(t, err)
		assert.NotEqual(t, "", sessionID, "セッションIDがあるか")
		assert.Regexp(t, r, u, "正しい形式で生成されているか")
	})

	t.Run("URLの組み立てに失敗", func(t *testing.T) {
		m.opts.Server = ":"
		_, _, err := m.createAuthorizeURL(permissions)
		assert.ErrorContains(t, err, "failed to create URL")
	})
}

func TestRecieveSessionID(t *testing.T) {
	type result struct {
		id  string
		err error
	}

	run := func(r chan *result, id string) {
		m := &Misskey{
			opts: &shared.ClientOpts{Name: "test_app", Server: "https://example.com"},
		}
		recieveID, err := m.recieveSessionID(id)
		r <- &result{id: recieveID, err: err}
	}

	postCallback := func(id string) (*http.Response, error) {
		q := &url.Values{}
		q.Set("session", id)
		url, _ := shared.CreateURL(q, "http://"+listenAddr, "callback")
		return http.Post(url, "", nil)
	}

	t.Run("セッションIDが受け取れるか", func(t *testing.T) {
		id, _ := uuid.NewUUID()
		wantSessionID := id.String()

		result := make(chan *result, 1)
		go run(result, wantSessionID)

		res, err := postCallback(wantSessionID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.NoError(t, r.err)
		assert.Equal(t, wantSessionID, r.id, "受け取ったセッションIDが生成したものと一致するか")
	})

	t.Run("セッションIDが一致しない場合エラーを返すか", func(t *testing.T) {
		id, _ := uuid.NewUUID()
		wantSessionID := id.String()

		result := make(chan *result, 1)
		go run(result, wantSessionID)

		res, err := postCallback("hogehoge")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.ErrorContains(t, r.err, "failed to recieve session id")
	})
}

func TestRecieveToken(t *testing.T) {
	isNotJSON := true
	isExpired := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/miauth/SESSION_ID/check" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if isNotJSON {
			isNotJSON = false
			fmt.Fprintln(w, `<html><head><title>Apps</title></head></html>`)
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

	t.Run("URLの組み立てに失敗", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: ":"}}
		_, err := m.recieveToken("SESSION_ID")
		assert.ErrorContains(t, err, "failed to create URL")
	})

	t.Run("リクエストに失敗", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: "http://localhost:9999"}}
		_, err := m.recieveToken("SESSION_ID")
		assert.ErrorContains(t, err, "failed to request")
	})

	t.Run("アクセス失敗", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: ts.URL}}
		_, err := m.recieveToken("hoge")
		assert.ErrorContains(t, err, "http error")
	})

	t.Run("JSONデコードエラー", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: ts.URL}}
		_, err := m.recieveToken("SESSION_ID")
		assert.ErrorContains(t, err, "failed to decord json")
	})

	t.Run("URL期限切れ", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: ts.URL}}
		_, err := m.recieveToken("SESSION_ID")
		assert.ErrorContains(t, err, "failed to get token")
	})

	t.Run("アクセストークンが取得できるか", func(t *testing.T) {
		m := &Misskey{opts: &shared.ClientOpts{Server: ts.URL}}
		res, err := m.recieveToken("SESSION_ID")
		assert.NoError(t, err)
		assert.Equal(t, "USER_TOKEN", res.Token)
	})
}