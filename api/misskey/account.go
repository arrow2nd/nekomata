package misskey

import (
	"github.com/arrow2nd/nekomata/api"
)

func (m *Misskey) SearchAccounts(query string, limit int) ([]*api.Account, error) {
	return nil, nil
}

func (m *Misskey) GetAccount(id string) (*api.Account, error) {
	return nil, nil
}

func (m *Misskey) GetRelationships(ids []string) ([]*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) GetPosts(id string, limit int) ([]*api.Post, error) {
	return nil, nil
}

func (m *Misskey) Follow(id string) (*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unfollow(id string) (*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Block(id string) (*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unblock(id string) (*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Mute(id string) (*api.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unmute(id string) (*api.Relationship, error) {
	return nil, nil
}
