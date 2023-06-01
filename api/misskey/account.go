package misskey

import (
	"github.com/arrow2nd/nekomata/api/shared"
)

func (m *Misskey) SearchAccounts(query string, limit int) ([]*shared.Account, error) {
	return nil, nil
}

func (m *Misskey) GetAccount(id string) (*shared.Account, error) {
	return nil, nil
}

func (m *Misskey) GetRelationships(ids []string) ([]*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) GetPosts(id string, limit int) ([]*shared.Post, error) {
	return nil, nil
}

func (m *Misskey) Follow(id string) (*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unfollow(id string) (*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Block(id string) (*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unblock(id string) (*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Mute(id string) (*shared.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unmute(id string) (*shared.Relationship, error) {
	return nil, nil
}
