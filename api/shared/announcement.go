package shared

import "time"

// Announcement : サービスからのお知らせ
type Announcement struct {
	ID          string
	PublishedAt time.Time
	UpdatedAt   *time.Time
	Title       string
	Text        string
}
