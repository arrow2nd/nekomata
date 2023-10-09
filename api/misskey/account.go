package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

func (m *Misskey) SearchAccounts(query string, limit int) ([]*sharedapi.Account, error) {
	return nil, nil
}

func (m *Misskey) GetAccount(id string) (*sharedapi.Account, error) {
	return nil, nil
}

func (m *Misskey) GetLoginAccount() (*sharedapi.Account, error) {
	return nil, nil
}

func (m *Misskey) GetRelationships(ids []string) ([]*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) GetPosts(id string, limit int) ([]*sharedapi.Post, error) {
	return nil, nil
}

func (m *Misskey) Follow(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unfollow(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Block(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unblock(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Mute(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}

func (m *Misskey) Unmute(id string) (*sharedapi.Relationship, error) {
	return nil, nil
}
