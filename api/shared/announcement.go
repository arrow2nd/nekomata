package shared

import "time"

type AnnouncementOpts struct {
	Limit int
}

// Announcement : サービスからのお知らせ
type Announcement struct {
	ID          string
	PublishedAt time.Time
	UpdatedAt   *time.Time
	Title       string
	Text        string
}
