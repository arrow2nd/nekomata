package mastodon

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/arrow2nd/nekomata/api"
)

func (m *Mastodon) getGlobalTimeline(sinceID string, limit int, local bool) ([]*api.Post, error) {
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

	opts := &requestOpts{
		method: http.MethodGet,
		url:    endpointTimelinePublic.URL(m.opts.Server, nil),
		q:      q,
		isAuth: true,
	}

	res := []*status{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return statuses2SharedPosts(res), nil
}

func (m *Mastodon) GetGlobalTimeline(sinceID string, limit int) ([]*api.Post, error) {
	return m.getGlobalTimeline(sinceID, limit, false)
}

func (m *Mastodon) GetLocalTimeline(sinceID string, limit int) ([]*api.Post, error) {
	return m.getGlobalTimeline(sinceID, limit, true)
}

func (m *Mastodon) GetHomeTimeline(sinceID string, limit int) ([]*api.Post, error) {
	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))

	if sinceID != "" {
		q.Add("since_id", sinceID)
	}

	opts := &requestOpts{
		method: http.MethodGet,
		url:    endpointTimelineHome.URL(m.opts.Server, nil),
		q:      q,
		isAuth: true,
	}

	res := []*status{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return statuses2SharedPosts(res), nil
}

func (m *Mastodon) GetListTimeline(listID, sinceID string, limit int) ([]*api.Post, error) {
	p := url.Values{}
	p.Add(":list_id", listID)

	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))

	if sinceID != "" {
		q.Add("since_id", sinceID)
	}

	opts := &requestOpts{
		method: http.MethodGet,
		url:    endpointTimelineList.URL(m.opts.Server, p),
		q:      q,
		isAuth: true,
	}

	res := []*status{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return statuses2SharedPosts(res), nil
}
