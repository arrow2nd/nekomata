package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

func (m *Misskey) GetGlobalTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) GetLocalTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) GetHomeTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) GetListTimeline(listID, sinceID string, limit int) ([]*sharedapi.Post, error) {
	return nil, nil
}
