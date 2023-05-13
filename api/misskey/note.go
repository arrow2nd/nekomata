package misskey

import "github.com/arrow2nd/nekomata/api/shared"

func (m *Misskey) CreatePost(opts *shared.CreatePostOpts) (*shared.Post, error) {
	return nil, nil
}

func (m Misskey) QuotePost(id string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	return nil, nil
}

func (m Misskey) ReplyPost(id string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	return nil, nil
}

func (m *Misskey) DeletePost(id string) (*shared.Post, error) {
	return nil, nil
}

func (m *Misskey) Reaction(id, reaction string) error {
	return nil
}

func (m *Misskey) UnReaction(id string) error {
	return nil
}

func (m *Misskey) Repost(id string) error {
	return nil
}

func (m *Misskey) UnRepost(id string) error {
	return nil
}

func (m *Misskey) Bookmark(id string) error {
	return nil
}

func (m *Misskey) UnBookmark(id string) error {
	return nil
}
