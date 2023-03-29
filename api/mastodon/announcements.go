package mastodon

import (
	"net/url"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"jaytaylor.com/html2text"
)

type announcement struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Mastodon) GetAnnouncements() ([]*shared.Announcement, error) {
	q := url.Values{}
	q.Add("with_dismissed", "false")

	res := []*announcement{}
	url := announcementsEndpoint.URL(m.opts.Server)
	if err := m.request("GET", url, q, true, &res); err != nil {
		return nil, err
	}

	results := []*shared.Announcement{}
	for _, r := range res {
		// Content は HTML 文字列なので普通の文字列に変換する
		text, err := html2text.FromString(r.Content, html2text.Options{PrettyTables: true})
		if err != nil {
			return nil, err
		}

		results = append(results, &shared.Announcement{
			ID:          r.ID,
			PublishedAt: r.PublishedAt,
			UpdatedAt:   &r.UpdatedAt,
			Title:       "",
			Text:        text,
		})
	}

	return results, nil
}
