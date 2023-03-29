package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestGetAnnouncements(t *testing.T) {
	isNotHTMLContent := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isNotHTMLContent {
			isNotHTMLContent = false
			fmt.Fprintln(w, `[ { "id": "0", "content": "This is plain text", "published_at": "2023-01-01T00:00:00.000Z", "updated_at": "2023-01-02T00:00:00.000Z" }`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `
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
    ]`)
	}))

	defer ts.Close()

	t.Run("Contentパースエラー", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
		_, err := m.GetAnnouncements()
		assert.Error(t, err)
	})

	t.Run("内容を取得できるか", func(t *testing.T) {
		m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
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
