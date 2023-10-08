package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	networkErr := true
	isNotJSON := true
	isError := true
	checkHeaders := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if networkErr {
			// アクセス失敗
			networkErr = false
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if isNotJSON {
			// デコード失敗
			isNotJSON = false
			fmt.Fprintln(w, `<html><head><title>Apps</title></head></html>`)
			return
		} else if isError {
			// エラーレスポンス
			isError = false
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{ "error": "invalid_scope", "error_description": "The requested scope is invalid, unknown, or malformed." }`)
			return
		} else if checkHeaders {
			// ヘッダーが正しい
			checkHeaders = false
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer", "認証情報がある")
			assert.Equal(t, r.Header.Get("Content-Type"), "application/json", "Content-Typeがある")
		} else {
			// メソッドが正しい
			assert.Equal(t, r.Method, http.MethodGet)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"ok": "true"}`)
	}))

	defer ts.Close()

	opts := &requestOpts{
		method: http.MethodPost,
		q:      nil,
		isAuth: false,
	}

	t.Run("リクエストに失敗", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: "http://localhost:9999"}}
		opts.url = endpointAnnouncements.URL(m.opts.Server, nil)

		err := m.request(opts, nil)
		e := &sharedapi.RequestError{}

		assert.ErrorAs(t, err, &e)
	})

	t.Run("アクセス失敗", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: ts.URL}}
		opts.url = endpointAnnouncements.URL(m.opts.Server, nil)

		err := m.request(opts, nil)
		e := &sharedapi.HTTPError{}

		assert.ErrorAs(t, err, &e)
	})

	t.Run("JSONデコードエラー", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: ts.URL}}
		opts.url = endpointAnnouncements.URL(m.opts.Server, nil)

		type a struct {
			hoge string
		}

		res := &a{}
		err := m.request(opts, &res)
		e := &sharedapi.DecodeError{}

		assert.ErrorAs(t, err, &e)
	})

	t.Run("エラーレスポンス", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: ts.URL}}
		opts.url = endpointAnnouncements.URL(m.opts.Server, nil)

		err := m.request(opts, nil)
		e := &errorResponse{}

		assert.ErrorAs(t, err, &e)
	})

	t.Run("指定した内容がヘッダーにあるか", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: ts.URL}}

		opts := &requestOpts{
			method:      http.MethodPost,
			contentType: "application/json",
			url:         endpointAnnouncements.URL(m.opts.Server, nil),
			q:           nil,
			isAuth:      true,
		}

		err := m.request(opts, nil)
		assert.NoError(t, err)
	})

	t.Run("指定したメソッドで送信できているか", func(t *testing.T) {
		m := &Mastodon{opts: &sharedapi.ClientOpts{Server: ts.URL}}

		opts := &requestOpts{
			method: http.MethodGet,
			url:    endpointAnnouncements.URL(m.opts.Server, nil),
			q:      nil,
			isAuth: false,
		}

		err := m.request(opts, nil)
		assert.NoError(t, err)
	})
}
