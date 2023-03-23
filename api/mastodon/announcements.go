package mastodon

import "github.com/arrow2nd/nekomata/api/shared"

type announcementsOpts struct {
	WithDismissed bool `json:"with_dismissed"`
}

type announcementsResponse struct {
	ID          string  `json:"id"`
	Content     string  `json:"content"`
	StartsAt    *string `json:"starts_at"`
	EndsAt      *string `json:"ends_at"`
	AllDay      bool    `json:"all_day"`
	PublishedAt string  `json:"published_at"`
	UpdatedAt   string  `json:"updated_at"`
	Read        bool    `json:"read"`
}

func (m *Mastodon) GetAnnouncements() ([]*shared.Announcement, error) {
	// res := []*announcementsResponse{}
	// req := &announcementsOpts{
	// 	WithDismissed: false,
	// }
	return nil, nil
}
