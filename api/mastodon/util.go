package mastodon

import (
	"html"
	"regexp"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// statuses2SharedPosts : []*statuseを[]*api.Postに変換
func statuses2SharedPosts(raw []*status) []*sharedapi.Post {
	posts := []*sharedapi.Post{}

	for _, status := range raw {
		posts = append(posts, status.ToShared())
	}

	return posts
}

// html2text : htmlをプレーンテキストに変換
func html2text(h string) string {
	text := regexp.MustCompile(`<p>(.*?)</p>`).ReplaceAllString(h, "$1\n")
	text = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile("<[^>]*>").ReplaceAllString(text, "")
	return strings.TrimSpace(html.UnescapeString(text))
}
