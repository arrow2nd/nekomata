package mastodon

import (
	"github.com/arrow2nd/nekomata/api/shared"
)

// statuses2SharedPosts : []*statuseを[]*shared.Postに変換
func statuses2SharedPosts(raw []*status) []*shared.Post {
	posts := []*shared.Post{}

	for _, status := range raw {
		posts = append(posts, status.ToShared())
	}

	return posts
}
