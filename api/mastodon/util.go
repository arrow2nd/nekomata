package mastodon

import "github.com/arrow2nd/nekomata/api/sharedapi"

// statuses2SharedPosts : []*statuseを[]*api.Postに変換
func statuses2SharedPosts(raw []*status) []*sharedapi.Post {
	posts := []*sharedapi.Post{}

	for _, status := range raw {
		posts = append(posts, status.ToShared())
	}

	return posts
}
