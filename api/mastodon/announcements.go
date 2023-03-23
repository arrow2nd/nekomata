package mastodon

import (
	"net/url"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"jaytaylor.com/html2text"
)

type announcementsResponse struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	PublishedAt string `json:"published_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (m *Mastodon) GetAnnouncements() ([]*shared.Announcement, error) {
	q := url.Values{}
	q.Add("with_dismissed", "false")

	res := []*announcementsResponse{}
	if err := m.request("GET", announcementsEndpoint, q, true, &res); err != nil {
		return nil, err
	}

	results := []*shared.Announcement{}
	for _, r := range res {
		publishedAt, _ := time.Parse(time.RFC3339, r.PublishedAt)
		updatedAt, _ := time.Parse(time.RFC3339, r.UpdatedAt)

		// Content は HTML 文字列なので普通の文字列に変換する
		text, err := html2text.FromString(r.Content, html2text.Options{PrettyTables: true})
		if err != nil {
			return nil, err
		}

		results = append(results, &shared.Announcement{
			ID:          r.ID,
			PublishedAt: publishedAt,
			UpdatedAt:   &updatedAt,
			Title:       "",
			Text:        text,
		})
	}

	return results, nil
}
