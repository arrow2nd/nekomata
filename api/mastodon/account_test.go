package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

const mockAccount = `
{
  "id": "1",
  "username": "hoge",
  "acct": "hoge",
  "display_name": "hoge",
  "locked": false,
  "bot": false,
  "discoverable": true,
  "group": false,
  "created_at": "2020-08-16T00:00:00.000Z",
  "note": "<p>おもち</p>",
  "url": "https://example.com/",
  "avatar": "https://example.com/",
  "avatar_static": "https://example.com/",
  "header": "https://example.com/",
  "header_static": "https://example.com/",
  "followers_count": 24,
  "following_count": 22,
  "statuses_count": 473,
  "last_status_at": "2023-05-29",
  "noindex": false,
  "emojis": [],
  "fields": [
    {
      "name": "first",
      "value": "<p>1st</p>",
      "verified_at": null
    },
    {
      "name": "second",
      "value": "<a href=\"https://example.com/\">hello!</a>",
      "verified_at": null
    }
  ]
}
`

var wantAccount = sharedapi.Account{
	ID:             "1",
	Username:       "hoge",
	DisplayName:    "hoge",
	Private:        false,
	Bot:            false,
	Verified:       false,
	BIO:            "おもち",
	CreatedAt:      time.Date(2020, time.August, 16, 0, 0, 0, 0, time.UTC),
	FollowersCount: 24,
	FollowingCount: 22,
	PostsCount:     473,
	Profiles: []sharedapi.Profile{
		{Label: "first", Value: "1st"},
		{Label: "second", Value: "hello! ( https://example.com/ )"},
	},
}

const mockRelationship = `
{
  "id": "0",
  "following": true,
  "showing_reblogs": false,
  "notifying": false,
  "followed_by": false,
  "blocking": false,
  "blocked_by": false,
  "muting": false,
  "muting_notifications": false,
  "requested": false,
  "domain_blocking": false,
  "endorsed": false
}
`

var wantRelationship = sharedapi.Relationship{
	ID:         "0",
	Following:  true,
	FollowedBy: false,
	Blocking:   false,
	BlockedBy:  false,
	Muting:     false,
	Requested:  false,
}

func createMockServer(t *testing.T, id, res string) *httptest.Server {
	isError := false

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isError {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{ "error": "Record not found" }`)
			return
		}

		assert.Contains(t, r.URL.String(), "/"+id, "パスパラメータにユーザーIDが含まれているか")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, res)
		isError = true
	}))
}

func TestAccountToShared(t *testing.T) {
	note := "こんにちわ"
	a := &account{
		ID:             "id",
		Acct:           "acct@example.com",
		DisplayName:    "おもち",
		Locked:         true,
		Bot:            true,
		CreatedAt:      time.Now(),
		Note:           "<p>" + note + "</p>",
		FollowersCount: 1,
		FollowingCount: 2,
		StatusesCount:  3,
		Fields: []accountFields{
			{Name: "好きなもの", Value: "おこめ"},
		},
	}

	got := a.ToShared()
	assert.ObjectsAreEqualValues(a, got)
}

func TestRelationshipToShared(t *testing.T) {
	r := &relationship{
		ID:         "id",
		Following:  true,
		FollowedBy: false,
		Blocking:   true,
		BlockedBy:  false,
		Muting:     true,
		Requested:  false,
	}

	got := r.ToShared()
	assert.ObjectsExportedFieldsAreEqual(r, got)
}

func TestSearchAccounts(t *testing.T) {
	isError := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isError {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, `{ "error": "Remote data could not be fetched" }`)
			return
		}

		assert.Contains(t, r.URL.String(), "?limit=1&q=hoge", "クエリパラメータが正しいか")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "[%s]", mockAccount)

		isError = true
	}))

	defer ts.Close()

	t.Run("取得できる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.SearchAccounts("hoge", 1)
		assert.NoError(t, err)

		assert.Equal(t, wantAccount, *r[0])
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.SearchAccounts("hoge", 1)
		assert.Error(t, err)
	})
}

func TestGetAccount(t *testing.T) {
	id := "1"
	ts := createMockServer(t, id, mockAccount)
	defer ts.Close()

	t.Run("取得できる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.GetAccount(id)
		assert.NoError(t, err)

		assert.Equal(t, wantAccount, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.GetAccount(id)
		assert.Error(t, err)
	})
}

func TestGetRelationships(t *testing.T) {
	isError := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isError {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{ "error": "The access token is invalid" }`)
			return
		}

		assert.Contains(t, r.URL.String(), "id%5B%5D=1234&id%5B%5D=5678", "クエリパラメータが正しいか")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "[%s, %s]", mockRelationship, mockRelationship)

		isError = true
	}))

	defer ts.Close()

	t.Run("取得できる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.GetRelationships([]string{"1234", "5678"})
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r[0])
		assert.Equal(t, wantRelationship, *r[1])
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.GetRelationships([]string{""})
		assert.Error(t, err)
	})
}

func TestGetPosts(t *testing.T) {
	id := "1"
	ts := createMockServer(t, id, `[{ "id": "0" }, { "id": "1" }]`)
	defer ts.Close()

	t.Run("取得できる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		res, err := m.GetPosts(id, 0)
		assert.NoError(t, err)

		assert.Equal(t, "0", res[0].ID)
		assert.Equal(t, "1", res[1].ID)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.GetPosts(id, 0)
		assert.Error(t, err)
	})
}

func TestFollow(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("フォローできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Follow(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Follow(id)
		assert.Error(t, err)
	})
}

func TestUnfollow(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("アンフォローできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Unfollow(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Unfollow(id)
		assert.Error(t, err)
	})
}

func TestBlock(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("ブロックできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Block(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Block(id)
		assert.Error(t, err)
	})
}

func TestUnblock(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("アンブロックできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Unblock(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Unblock(id)
		assert.Error(t, err)
	})
}

func TestMute(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("ミュートできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Mute(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Mute(id)
		assert.Error(t, err)
	})
}

func TestUnmute(t *testing.T) {
	id := "012345"
	ts := createMockServer(t, id, mockRelationship)
	defer ts.Close()

	t.Run("アンミュートできる", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		r, err := m.Unmute(id)
		assert.NoError(t, err)

		assert.Equal(t, wantRelationship, *r)
	})

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(nil, &sharedapi.UserCredential{Server: ts.URL})

		_, err := m.Unmute(id)
		assert.Error(t, err)
	})
}
