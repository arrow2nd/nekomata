package mastodon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
