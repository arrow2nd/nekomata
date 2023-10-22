package mastodon

import (
	"net/http"
	"net/url"
	"time"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// announcement : アナウンス
type announcement struct {
	// ID : アナウンス ID
	ID string `json:"id"`
	// Content : 内容
	Content string `json:"content"`
	// PublishedAt : 配信日
	PublishedAt time.Time `json:"published_at"`
	// UpdatedAt : 更新日
	UpdatedAt time.Time `json:"updated_at"`
}

// ToShared : api.Announcement に変換
func (a *announcement) ToShared() *sharedapi.Announcement {
	return &sharedapi.Announcement{
		ID:          a.ID,
		PublishedAt: a.PublishedAt,
		UpdatedAt:   &a.UpdatedAt,
		Title:       "",
		Text:        html2text(a.Content),
	}
}

func (m *Mastodon) GetAnnouncements() ([]*sharedapi.Announcement, error) {
	q := url.Values{}
	q.Add("with_dismissed", "false")

	opts := &requestOpts{
		method: http.MethodGet,
		url:    endpointAnnouncements.URL(m.user.Server, nil),
		q:      q,
		isAuth: true,
	}

	res := []*announcement{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	results := []*sharedapi.Announcement{}
	for _, r := range res {
		results = append(results, r.ToShared())
	}

	return results, nil
}
