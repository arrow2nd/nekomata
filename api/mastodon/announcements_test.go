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

const mockAnnouncements = `
[
  {
    "id": "0",
    "content": "<p>text_1</p>",
    "starts_at": null,
    "ends_at": null,
    "all_day": false,
    "published_at": "2023-01-01T00:00:00.000Z",
    "updated_at": "2023-01-02T00:00:00.000Z",
    "read": true,
    "mentions": [],
    "statuses": [],
    "tags": [],
    "emojis": [],
    "reactions": []
  },
  {
    "id": "1",
    "content": "<p>text_2</p>",
    "starts_at": null,
    "ends_at": null,
    "all_day": false,
    "published_at": "2022-01-01T00:40:00.000Z",
    "updated_at": "2022-01-02T00:00:00.000Z",
    "read": true,
    "mentions": [],
    "statuses": [],
    "tags": [],
    "emojis": [],
    "reactions": []
  }
]`

func TestAnnouncementToShared(t *testing.T) {
	content := "hoge"

	a := &announcement{
		ID:          "id",
		Content:     "<p>" + content + "<p>",
		PublishedAt: time.Now(),
		UpdatedAt:   time.Now().Add(time.Hour),
	}

	got := a.ToShared()
	assert.Equal(t, a.ID, got.ID, "IDが一致")
	assert.Equal(t, content, got.Text, "本文が一致")
	assert.Equal(t, a.PublishedAt, got.PublishedAt, "公開日が一致")
	assert.Equal(t, a.UpdatedAt, *got.UpdatedAt, "更新日が一致")
}

func TestGetAnnouncements(t *testing.T) {
	isError := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isError {
			isError = false
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{ "error": "The access token is invalid" }`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockAnnouncements)
	}))

	defer ts.Close()

	t.Run("エラーが返る", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})

		_, err := m.GetAnnouncements()
		assert.Error(t, err)
	})

	t.Run("取得できる", func(t *testing.T) {
		m := New(&shared.ClientOpts{Server: ts.URL})

		res, err := m.GetAnnouncements()
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		assert.Equal(t, "0", res[0].ID)
		assert.Equal(t, int64(1672531200), res[0].PublishedAt.Unix())
		assert.Equal(t, int64(1672617600), res[0].UpdatedAt.Unix())
		assert.Equal(t, "text_1", res[0].Text)

		assert.Equal(t, "1", res[1].ID)
		assert.Equal(t, int64(1640997600), res[1].PublishedAt.Unix())
		assert.Equal(t, int64(1641081600), res[1].UpdatedAt.Unix())
		assert.Equal(t, "text_2", res[1].Text)
	})
}
