package misskey

import (
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
)

type announcementsOpts struct {
	WithUnreads bool `json:"withUnreads"`
}

type announcementsResponse struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Text      string     `json:"text"`
	Title     string     `json:"title"`
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
		results = append(results, &shared.Announcement{
			ID:          r.ID,
			PublishedAt: r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
			Title:       r.Title,
			Text:        r.Text,
		})
	}

	return results, nil
}
