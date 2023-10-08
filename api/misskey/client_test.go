package misskey

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	networkErr := true
	isNotJSON := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if networkErr {
			networkErr = false
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if isNotJSON {
			isNotJSON = false
			fmt.Fprintln(w, `<html><head><title>Apps</title></head></html>`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{ "ok": true }`)
	}))

	defer ts.Close()

	t.Run("リクエストに失敗", func(t *testing.T) {
		m := &Misskey{opts: &sharedapi.ClientOpts{Server: "http://localhost:9999"}}
		err := m.post(endpointAnnouncements, &announcementsOpts{}, &announcementsResponse{})
		e := &sharedapi.RequestError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("アクセス失敗", func(t *testing.T) {
		m := &Misskey{opts: &sharedapi.ClientOpts{Server: ts.URL}}
		err := m.post(endpointAnnouncements, &announcementsOpts{}, &announcementsResponse{})
		e := &sharedapi.HTTPError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("JSONデコードエラー", func(t *testing.T) {
		m := &Misskey{opts: &sharedapi.ClientOpts{Server: ts.URL}}
		err := m.post(endpointAnnouncements, &announcementsOpts{}, &announcementsResponse{})
		e := &sharedapi.DecodeError{}
		assert.ErrorAs(t, err, &e)
	})

	t.Run("正常", func(t *testing.T) {
		type r struct {
			OK bool `json:"ok"`
		}

		m := &Misskey{opts: &sharedapi.ClientOpts{Server: ts.URL}}

		res := &r{}
		err := m.post(endpointAnnouncements, &announcementsOpts{}, &res)
		assert.NoError(t, err)

		assert.True(t, res.OK)
	})
}
