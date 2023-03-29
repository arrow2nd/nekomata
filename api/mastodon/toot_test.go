package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	type result struct {
		s string
		v string
	}

	serverRes := make(chan *result, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverRes <- &result{
			s: r.URL.Query().Get("status"),
			v: r.URL.Query().Get("visibility"),
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `
    {
      "id": "000000",
      "created_at": "2023-01-01T00:00:00.000Z",
      "in_reply_to_id": null,
      "in_reply_to_account_id": null,
      "sensitive": false,
      "spoiler_text": "",
      "visibility": "public",
      "language": "en",
      "uri": "https://mastodon.social/users/User/statuses/000000",
      "url": "https://mastodon.social/@User/000000",
      "replies_count": 5,
      "reblogs_count": 6,
      "favourites_count": 10,
      "favourited": false,
      "reblogged": false,
      "muted": false,
      "bookmarked": false,
      "content": "<p>Hello!</p>",
      "reblog": null,
      "application": { "name": "nekomata for term", "website": null },
      "account": {
        "id": "0",
        "username": "User",
        "acct": "User",
        "display_name": "User",
        "locked": false,
        "bot": false,
        "discoverable": true,
        "group": false,
        "created_at": "2023-01-01T00:00:00.000Z",
        "note": "<p>BIO</p>",
        "url": "https://mastodon.example.com/@User",
        "avatar": "https://files.mastodon.example.com/accounts/avatars/000/000/001/original/pic.jpg",
        "avatar_static": "https://files.mastodon.example.com/accounts/avatars/000/000/001/original/pic.jpg",
        "header": "https://files.mastodon.example.com/accounts/headers/000/000/001/original/pic.png",
        "header_static": "https://files.mastodon.example.com/accounts/headers/000/000/001/original/pic.png",
        "followers_count": 123,
        "following_count": 456,
        "statuses_count": 1000,
        "last_status_at": "2023-01-01T00:00:00.000Z",
        "emojis": [],
        "fields": [
          {
            "name": "A",
            "value": "a",
            "verified_at": null
          },
          {
            "name": "B",
            "value": "b",
            "verified_at": "2023-01-01T00:00:00.000+00:00"
          }
        ]
      },
      "media_attachments": [],
      "mentions": [],
      "tags": [],
      "emojis": [],
      "card": null,
      "poll": null
    }`)
	}))

	defer ts.Close()

	postText := "Hello!"
	postVisibility := "public"
	opts := &shared.CreatePostOpts{Text: postText, Visibility: postVisibility}

	m := &Mastodon{opts: &shared.ClientOpts{Server: ts.URL}}
	res, err := m.CreatePost(opts)

	assert.NoError(t, err)

	t.Run("サーバーにデータを送信できているか", func(t *testing.T) {
		res := <-serverRes
		assert.Equal(t, postText, res.s)
		assert.Equal(t, postVisibility, res.v)
	})

	t.Run("レスポンスをパースできるか", func(t *testing.T) {
		assert.Equal(t, postText, res.Text)
		assert.Equal(t, postVisibility, res.Visibility)
	})
}
