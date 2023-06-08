package mastodon

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/arrow2nd/nekomata/api/shared"
)

func (m *Mastodon) getGlobalTimeline(sinceID string, limit int, local bool) ([]*shared.Post, error) {
	endpoint := endpointTimelinePublic.URL(m.opts.Server, nil)

	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))

	if local {
		q.Add("local", "true")
	} else {
		q.Add("local", "false")
	}

	if sinceID != "" {
		q.Add("since_id", sinceID)
	}

	res := []*status{}
	if err := m.request(http.MethodGet, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	return statuses2SharedPosts(res), nil
}

func (m *Mastodon) GetGlobalTimeline(sinceID string, limit int) ([]*shared.Post, error) {
	return m.getGlobalTimeline(sinceID, limit, false)
}

func (m *Mastodon) GetLocalTimeline(sinceID string, limit int) ([]*shared.Post, error) {
	return m.getGlobalTimeline(sinceID, limit, true)
}

func (m *Mastodon) GetHomeTimeline(sinceID string, limit int) ([]*shared.Post, error) {
	return nil, nil
}

func (m *Mastodon) GetListTimeline(sinceID string, limit int) ([]*shared.Post, error) {
	return nil, nil
}
