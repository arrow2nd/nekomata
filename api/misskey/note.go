package misskey

import "github.com/arrow2nd/nekomata/api"

func (m *Misskey) CreatePost(opts *api.CreatePostOpts) (*api.Post, error) {
	return nil, nil
}

func (m Misskey) QuotePost(id string, opts *api.CreatePostOpts) (*api.Post, error) {
	return nil, nil
}

func (m Misskey) ReplyPost(id string, opts *api.CreatePostOpts) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) DeletePost(id string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Reaction(id, reaction string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Unreaction(id string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Repost(id string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Unrepost(id string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Bookmark(id string) (*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Unbookmark(id string) (*api.Post, error) {
	return nil, nil
}
