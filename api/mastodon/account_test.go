package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

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

var wantRelationship = shared.Relationship{
	ID:         "0",
	Following:  true,
	FollowedBy: false,
	Blocking:   false,
	BlockedBy:  false,
	Muting:     false,
	Requested:  false,
}

func createMockServer(t *testing.T, id string) *httptest.Server {
	isError := false

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), id, "URLにユーザーIDが含まれているか")

		if isError {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{ "error": "Record not found" }`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockRelationship)
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
	assert.Equal(t, a.ID, got.ID)
	assert.Equal(t, a.Acct, got.Username)
	assert.Equal(t, a.DisplayName, got.DisplayName)
	assert.Equal(t, a.Locked, got.Private)
	assert.Equal(t, a.Bot, got.Bot)
	assert.Equal(t, a.CreatedAt, got.CreatedAt)
	assert.Equal(t, note, got.BIO)
	assert.Equal(t, a.FollowersCount, got.FollowersCount)
	assert.Equal(t, a.FollowingCount, got.FollowingCount)
	assert.Equal(t, a.Fields[0].Name, got.Profiles[0].Label)
	assert.Equal(t, a.Fields[0].Value, got.Profiles[0].Value)
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
	assert.Equal(t, r.ID, got.ID)
	assert.True(t, got.Following)
	assert.False(t, got.FollowedBy)
	assert.True(t, got.Blocking)
	assert.False(t, got.BlockedBy)
	assert.True(t, got.Muting)
	assert.False(t, got.Requested)
}

func TestFollow(t *testing.T) {
	id := "012345"

	ts := createMockServer(t, id)
	defer ts.Close()

	t.Run("成功", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		r, err := m.Follow(id)
		assert.Equal(t, wantRelationship, *r)
		assert.NoError(t, err)
	})

	t.Run("失敗", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		_, err := m.Follow(id)
		assert.Error(t, err)
	})
}

func TestUnfollow(t *testing.T) {
	id := "012345"

	ts := createMockServer(t, id)
	defer ts.Close()

	t.Run("成功", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		r, err := m.Unfollow(id)
		assert.Equal(t, wantRelationship, *r)
		assert.NoError(t, err)
	})

	t.Run("失敗", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		_, err := m.Unfollow(id)
		assert.Error(t, err)
	})
}

func TestBlock(t *testing.T) {
	id := "012345"

	ts := createMockServer(t, id)
	defer ts.Close()

	t.Run("成功", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		r, err := m.Block(id)
		assert.Equal(t, wantRelationship, *r)
		assert.NoError(t, err)
	})

	t.Run("失敗", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		_, err := m.Block(id)
		assert.Error(t, err)
	})
}

func TestUnblock(t *testing.T) {
	id := "012345"

	ts := createMockServer(t, id)
	defer ts.Close()

	t.Run("成功", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		r, err := m.Unblock(id)
		assert.Equal(t, wantRelationship, *r)
		assert.NoError(t, err)
	})

	t.Run("失敗", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})
		_, err := m.Unblock(id)
		assert.Error(t, err)
	})
}
