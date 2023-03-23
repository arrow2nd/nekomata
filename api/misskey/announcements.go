package misskey

import (
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
)

type announcementsOpts struct {
	WithUnreads bool `json:"withUnreads"`
}

type announcementsResponse struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	Text      string  `json:"text"`
	Title     string  `json:"title"`
}

func (m *Misskey) GetAnnouncements() ([]*shared.Announcement, error) {
	req := &announcementsOpts{
		WithUnreads: false,
	}

	res := []*announcementsResponse{}
	if err := m.post(announcementsEndpoint, req, &res); err != nil {
		return nil, err
	}

	results := []*shared.Announcement{}
	for _, r := range res {
		publishedAt, _ := time.Parse(time.RFC3339, r.CreatedAt)

		a := &shared.Announcement{
			ID:          r.ID,
			PublishedAt: publishedAt,
			UpdatedAt:   nil,
			Title:       r.Title,
			Text:        r.Text,
		}

		if r.UpdatedAt != nil {
			u, _ := time.Parse(time.RFC3339, *r.UpdatedAt)
			a.UpdatedAt = &u
		}

		results = append(results, a)
	}

	return results, nil
}
