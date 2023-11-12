package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

func (m *Misskey) CreatePost(opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	return nil, nil
}

func (m Misskey) QuotePost(id string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	return nil, nil
}

func (m Misskey) ReplyPost(id string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) DeletePost(id string) error {
	return nil
}

func (m *Misskey) Reaction(id, reaction string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Unreaction(id string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Repost(id string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Unrepost(id string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Bookmark(id string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Unbookmark(id string) (*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) GetVisibilityList() []string {
	return []string{}
}
