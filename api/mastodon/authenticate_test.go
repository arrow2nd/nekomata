package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthorizeURL(t *testing.T) {
	m := &Mastodon{
		client: &sharedapi.ClientCredential{
			ID: "hoge",
		},
		user: &sharedapi.UserCredential{
			Server: "https://example.com",
		},
	}

	u := m.createAuthorizeURL("aaaa bbbb")

	endpoint := endpointOauthAuthorize.URL(m.user.Server, nil)
	want := endpoint + "?client_id=hoge&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&response_type=code&scope=aaaa+bbbb"
	assert.Equal(t, want, u)
}

func TestRecieveToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{ "access_token": "USER_TOKEN", "token_type": "Bearer", "scope": "read write follow push", "created_at": 0 }`)
	}))

	defer ts.Close()

	m := &Mastodon{
		client: &sharedapi.ClientCredential{
			Name:   "test",
			ID:     "id",
			Secret: "secret",
		},
		user: &sharedapi.UserCredential{
			Server: ts.URL,
		},
	}

	token, err := m.recieveToken("CODE")
	assert.NoError(t, err)
	assert.Equal(t, "USER_TOKEN", token)
}
