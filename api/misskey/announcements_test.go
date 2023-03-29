package misskey_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arrow2nd/nekomata/api"
	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

const mockAnnouncements = `
[
  {
    "id": "hogehoge",
    "createdAt": "2023-01-01T00:00:00.000Z",
    "updatedAt": "2023-01-02T00:00:00.000Z",
    "text": "text_1",
    "title": "title_1",
    "imageUrl": null
  },
  {
    "id": "fugafuga",
    "createdAt": "2022-01-01T00:40:00.000Z",
    "updatedAt": "2022-01-02T00:00:00.000Z",
    "text": "text_2",
    "title": "title_2",
    "imageUrl": null
  }
]`

const mockAnnouncementsWithNull = `
[{
  "id": "hogehoge",
  "createdAt": "2023-01-01T00:00:00.000Z",
  "updatedAt": null,
  "text_1": "text",
  "title_1": "title",
  "imageUrl": null
}]`

func TestGetAnnouncements(t *testing.T) {
	nullUpdatedAt := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if nullUpdatedAt {
			nullUpdatedAt = false
			fmt.Fprintln(w, mockAnnouncementsWithNull)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockAnnouncements)
	}))

	defer ts.Close()

	t.Run("nullのフィールドに対応できるか", func(t *testing.T) {
		m, _ := api.NewClient(os.Stdout, api.ServiceMisskey, &shared.ClientOpts{Server: ts.URL})
		res, err := m.GetAnnouncements()
		assert.NoError(t, err)
		assert.Nil(t, res[0].UpdatedAt)
	})

	t.Run("内容を取得できるか", func(t *testing.T) {
		m, _ := api.NewClient(os.Stdout, api.ServiceMisskey, &shared.ClientOpts{Server: ts.URL})
		res, err := m.GetAnnouncements()
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		assert.Equal(t, "hogehoge", res[0].ID)
		assert.Equal(t, int64(1672531200), res[0].PublishedAt.Unix())
		assert.Equal(t, int64(1672617600), res[0].UpdatedAt.Unix())
		assert.Equal(t, "title_1", res[0].Title)
		assert.Equal(t, "text_1", res[0].Text)

		assert.Equal(t, "fugafuga", res[1].ID)
		assert.Equal(t, int64(1640997600), res[1].PublishedAt.Unix())
		assert.Equal(t, int64(1641081600), res[1].UpdatedAt.Unix())
		assert.Equal(t, "title_2", res[1].Title)
		assert.Equal(t, "text_2", res[1].Text)
	})
}
