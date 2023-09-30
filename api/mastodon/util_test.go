package mastodon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatuses2SharedPosts(t *testing.T) {
	raw := []*status{
		{
			ID:      "0",
			Content: "hoge",
		},
		{
			ID:      "1",
			Content: "fuga",
		},
	}

	posts := statuses2SharedPosts(raw)

	assert.Len(t, raw, 2)
	assert.Equal(t, "0", posts[0].ID)
	assert.Equal(t, "hoge", posts[0].Text)
	assert.Equal(t, "1", posts[1].ID)
	assert.Equal(t, "fuga", posts[1].Text)
}
