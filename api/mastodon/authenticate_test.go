package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthorizeURL(t *testing.T) {
	m := &Mastodon{opts: &shared.ClientOpts{Server: "https://example.com", ID: "hoge"}}
	u := m.createAuthorizeURL([]string{"aaaa", "bbbb"})

	want := oauthAuthorizeEndpoint.URL(m.opts.Server) + "?client_id=hoge&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&response_type=code&scope=aaaa+bbbb"
	assert.Equal(t, want, u)
}

func TestRecieveCode(t *testing.T) {
	type result struct {
		code string
		err  error
	}

	run := func(r chan *result) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: "https://example.com"}}
		recieveCode, err := m.recieveCode()
		r <- &result{code: recieveCode, err: err}
	}

	postCallback := func(code string) (*http.Response, error) {
		q := url.Values{}
		q.Add("code", code)
		return http.Post(shared.AuthCallbackURL+"?"+q.Encode(), "", nil)
	}

	t.Run("受け取れるか", func(t *testing.T) {
		result := make(chan *result, 1)
		go run(result)

		wantCode := "CODE"
		res, err := postCallback(wantCode)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.NoError(t, r.err)
		assert.Equal(t, wantCode, r.code)
	})

	t.Run("空の場合エラーを返すか", func(t *testing.T) {
		result := make(chan *result, 1)
		go run(result)

		res, err := postCallback("")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.ErrorContains(t, r.err, "failed to recieve code")
	})
}

func TestRecieveToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{ "access_token": "USER_TOKEN", "token_type": "Bearer", "scope": "read write follow push", "created_at": 0 }`)
	}))

	defer ts.Close()

	m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
	res, err := m.recieveToken("CODE")
	assert.NoError(t, err)
	assert.Equal(t, "USER_TOKEN", res.Token)
}
