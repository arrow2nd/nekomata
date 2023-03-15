package shared_test

import (
	"net/url"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestCreateURL(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		q := &url.Values{}
		q.Add("id", "hoge")
		u, err := shared.CreateURL(q, "https://example.com", "test")
		want := "https://example.com/test?id=hoge"
		assert.NoError(t, err)
		assert.Equal(t, want, u)
	})

	t.Run("異常", func(t *testing.T) {
		_, err := shared.CreateURL(nil, ":", "test")
		assert.ErrorContains(t, err, "failed to create URL")
	})

	t.Run("パスを省略した場合", func(t *testing.T) {
		u, err := shared.CreateURL(nil, "https://example.com")
		want := "https://example.com"
		assert.NoError(t, err)
		assert.Equal(t, want, u)
	})

	t.Run("クエリパラメータを省略した場合", func(t *testing.T) {
		u, err := shared.CreateURL(nil, "https://example.com", "test")
		want := "https://example.com/test"
		assert.NoError(t, err)
		assert.Equal(t, want, u)
	})
}
