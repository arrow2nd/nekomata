package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	networkErr := true
	isNotJSON := true
	isError := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a := r.Header.Get("Authorization"); a != "" && strings.HasPrefix(a, "Bearer") {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{ "s": "authorization" }`)
			return
		} else if networkErr {
			networkErr = false
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if isNotJSON {
			isNotJSON = false
			fmt.Fprintln(w, `<html><head><title>Apps</title></head></html>`)
			return
		} else if isError {
			isError = false
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{ "error": "invalid_scope", "error_description": "The requested scope is invalid, unknown, or malformed." }`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{ "s": "%s" }`, r.Method)
	}))

	defer ts.Close()

	t.Run("リクエストに失敗", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: "http://localhost:9999"}}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("POST", u, nil, false, nil)
		e := &shared.RequestError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("アクセス失敗", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("POST", u, nil, false, nil)
		e := &shared.HTTPError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("JSONデコードエラー", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("POST", u, nil, false, nil)
		e := &shared.DecodeError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("エラーレスポンス", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("POST", u, nil, false, nil)
		e := &errorResponse{}
		assert.ErrorAs(t, err, &e)
	})

	type r struct {
		S string `json:"s"`
	}

	t.Run("認証情報がヘッダーにあるか", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		res := &r{}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("POST", u, nil, true, res)
		assert.NoError(t, err)
		assert.Equal(t, "authorization", res.S)
	})

	t.Run("指定したメソッドで送信できているか", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		res := &r{}
		u := announcementsEndpoint.URL(m.opts.Server, nil)
		err := m.request("GET", u, nil, false, res)
		assert.NoError(t, err)
		assert.Equal(t, "GET", res.S)
	})
}