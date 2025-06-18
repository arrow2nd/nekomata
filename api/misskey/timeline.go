package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

// TimelineRequest : タイムライン取得のリクエスト
type TimelineRequest struct {
	I           string  `json:"i"`
	WithFiles   *bool   `json:"withFiles,omitempty"`
	WithRenotes *bool   `json:"withRenotes,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
	SinceId     *string `json:"sinceId,omitempty"`
	UntilId     *string `json:"untilId,omitempty"`
	ExcludeNsfw *bool   `json:"excludeNsfw,omitempty"`
}

func (m *Misskey) GetGlobalTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	req := TimelineRequest{
		I:     m.user.Token,
		Limit: &limit,
	}

	if sinceID != "" {
		req.SinceId = &sinceID
	}

	var notes []MisskeyNote
	if err := m.post(endpointNotesGlobalTimeline, req, &notes); err != nil {
		return nil, err
	}

	posts := make([]*sharedapi.Post, len(notes))
	for i, note := range notes {
		posts[i] = convertToPost(&note)
	}

	return posts, nil
}

func (m *Misskey) GetLocalTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	req := TimelineRequest{
		I:     m.user.Token,
		Limit: &limit,
	}

	if sinceID != "" {
		req.SinceId = &sinceID
	}

	var notes []MisskeyNote
	if err := m.post(endpointNotesLocalTimeline, req, &notes); err != nil {
		return nil, err
	}

	posts := make([]*sharedapi.Post, len(notes))
	for i, note := range notes {
		posts[i] = convertToPost(&note)
	}

	return posts, nil
}

// HomeTimelineRequest : ホームタイムライン取得のリクエスト
type HomeTimelineRequest struct {
	I                        string  `json:"i"`
	Limit                    *int    `json:"limit,omitempty"`
	SinceId                  *string `json:"sinceId,omitempty"`
	UntilId                  *string `json:"untilId,omitempty"`
	IncludeMyRenotes         *bool   `json:"includeMyRenotes,omitempty"`
	IncludeRenotedMyNotes    *bool   `json:"includeRenotedMyNotes,omitempty"`
	IncludeLocalRenotes      *bool   `json:"includeLocalRenotes,omitempty"`
	WithFiles                *bool   `json:"withFiles,omitempty"`
	WithRenotes              *bool   `json:"withRenotes,omitempty"`
}

func (m *Misskey) GetHomeTimeline(sinceID string, limit int) ([]*sharedapi.Post, error) {
	req := HomeTimelineRequest{
		I:     m.user.Token,
		Limit: &limit,
	}

	if sinceID != "" {
		req.SinceId = &sinceID
	}

	var notes []MisskeyNote
	if err := m.post(endpointNotesTimeline, req, &notes); err != nil {
		return nil, err
	}

	posts := make([]*sharedapi.Post, len(notes))
	for i, note := range notes {
		posts[i] = convertToPost(&note)
	}

	return posts, nil
}

// UserListTimelineRequest : ユーザーリストタイムライン取得のリクエスト
type UserListTimelineRequest struct {
	I                        string  `json:"i"`
	ListId                   string  `json:"listId"`
	Limit                    *int    `json:"limit,omitempty"`
	SinceId                  *string `json:"sinceId,omitempty"`
	UntilId                  *string `json:"untilId,omitempty"`
	IncludeMyRenotes         *bool   `json:"includeMyRenotes,omitempty"`
	IncludeRenotedMyNotes    *bool   `json:"includeRenotedMyNotes,omitempty"`
	IncludeLocalRenotes      *bool   `json:"includeLocalRenotes,omitempty"`
	WithFiles                *bool   `json:"withFiles,omitempty"`
	WithRenotes              *bool   `json:"withRenotes,omitempty"`
}

func (m *Misskey) GetListTimeline(listID, sinceID string, limit int) ([]*sharedapi.Post, error) {
	req := UserListTimelineRequest{
		I:      m.user.Token,
		ListId: listID,
		Limit:  &limit,
	}

	if sinceID != "" {
		req.SinceId = &sinceID
	}

	var notes []MisskeyNote
	if err := m.post(endpointNotesUserListTimeline, req, &notes); err != nil {
		return nil, err
	}

	posts := make([]*sharedapi.Post, len(notes))
	for i, note := range notes {
		posts[i] = convertToPost(&note)
	}

	return posts, nil
}
