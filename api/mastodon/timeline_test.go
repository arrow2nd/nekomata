package mastodon_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/mastodon"
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

var clientCred = &sharedapi.ClientCredential{
	Name:   "hoge",
	ID:     "fuga",
	Secret: "piyo",
}

func createTestServer(t *testing.T, local string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0", r.FormValue("since_id"))
		assert.Equal(t, "5", r.FormValue("limit"))
		assert.Equal(t, local, r.FormValue("local"))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": "1"}]`)
	}))
}

func TestGetGlobalTimeLine(t *testing.T) {
	ts := createTestServer(t, "false")
	defer ts.Close()

	m, err := mastodon.New(clientCred, &sharedapi.UserCredential{Server: ts.URL})
	assert.NoError(t, err)

	posts, err := m.GetGlobalTimeline("0", 5)
	assert.NoError(t, err)

	assert.Len(t, posts, 1)
}

func TestGetLocalTimeLine(t *testing.T) {
	ts := createTestServer(t, "true")
	defer ts.Close()

	m, err := mastodon.New(clientCred, &sharedapi.UserCredential{Server: ts.URL})
	assert.NoError(t, err)

	posts, err := m.GetLocalTimeline("0", 5)
	assert.NoError(t, err)

	assert.Len(t, posts, 1)
}

func TestGetHomeTimeLine(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0", r.FormValue("since_id"))
		assert.Equal(t, "5", r.FormValue("limit"))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": "1"}]`)
	}))

	defer ts.Close()

	m, err := mastodon.New(clientCred, &sharedapi.UserCredential{Server: ts.URL})
	assert.NoError(t, err)

	posts, err := m.GetHomeTimeline("0", 5)
	assert.NoError(t, err)

	assert.Len(t, posts, 1)
}

func TestGetListTimeLine(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), "/12345")
		assert.Equal(t, "0", r.FormValue("since_id"))
		assert.Equal(t, "5", r.FormValue("limit"))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": "1"}]`)
	}))

	defer ts.Close()

	m, err := mastodon.New(clientCred, &sharedapi.UserCredential{Server: ts.URL})
	assert.NoError(t, err)

	posts, err := m.GetListTimeline("12345", "0", 5)
	assert.NoError(t, err)

	assert.Len(t, posts, 1)
}
