package misskey

import (
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
)

type announcementsOpts struct {
	Limit       int  `json:"limit"`
	WithUnreads bool `json:"withUnreads"`
}

type announcementsResponse struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	Text      string  `json:"text"`
	Title     string  `json:"title"`
	ImageURL  string  `json:"imageUrl"`
}

func (m *Misskey) GetAnnouncements() ([]*shared.Announcement, error) {
	res := []*announcementsResponse{}
	req := &announcementsOpts{
		Limit:       10,
		WithUnreads: false,
	}

	if err := m.post("announcements", req, &res); err != nil {
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
