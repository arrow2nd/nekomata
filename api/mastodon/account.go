package mastodon

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"jaytaylor.com/html2text"
)

// account : ユーザー情報
type account struct {
	// ID : ユーザーID
	ID string `json:"id"`
	// Acct : ユーザー名 + ドメイン名からなる文字列 (username@domain)
	Acct string `json:"acct"`
	// DisplayName : 表示名
	DisplayName string `json:"display_name"`
	// Locked : 非公開アカウントか
	Locked bool `json:"locked"`
	// Bot : ボットアカウントか
	Bot bool `json:"bot"`
	// CreatedAt : アカウント作成日
	CreatedAt time.Time `json:"created_at"`
	// Note : BIO
	Note string `json:"note"`
	// FollowersCount : フォロワー数
	FollowersCount int `json:"followers_count"`
	// FollowingCount : フォロイー数
	FollowingCount int `json:"following_count"`
	// StatusesCount : トゥート数
	StatusesCount int `json:"statuses_count"`
	// Fields : カスタムフィールド
	Fields []accountFields `json:"fields"`
}

// accountFields : カスタムフィールド
type accountFields struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ToShared : shared.Account に変換
func (a *account) ToShared() *shared.Account {
	// BIOをプレーンテキストに変換
	bio, err := html2text.FromString(a.Note)
	if err != nil {
		bio = fmt.Sprintf("convert error: %s", err)
	}

	// フィールドをプロフィールに変換
	profiles := []shared.Profile{}
	for _, p := range a.Fields {
		value, err := html2text.FromString(p.Value)
		if err != nil {
			value = fmt.Sprintf("convert error: %s", err)
		}

		profiles = append(profiles, shared.Profile{
			Label: p.Name,
			Value: value,
		})
	}

	return &shared.Account{
		ID:             a.ID,
		Username:       a.Acct,
		DisplayName:    a.DisplayName,
		Private:        a.Locked,
		Bot:            a.Bot,
		Verified:       false,
		BIO:            bio,
		CreatedAt:      a.CreatedAt,
		FollowersCount: a.FollowersCount,
		FollowingCount: a.FollowingCount,
		PostsCount:     a.StatusesCount,
		Profiles:       profiles,
	}
}

// relationship : ユーザーとの関係
type relationship struct {
	ID         string `json:"id"`
	Following  bool   `json:"following"`
	FollowedBy bool   `json:"followed_by"`
	Blocking   bool   `json:"blocking"`
	BlockedBy  bool   `json:"blocked_by"`
	Muting     bool   `json:"muting"`
	Requested  bool   `json:"requested"`
}

// ToShared : shared.Relation に変換
func (r *relationship) ToShared() *shared.Relationship {
	return &shared.Relationship{
		ID:         r.ID,
		Following:  r.Following,
		FollowedBy: r.FollowedBy,
		Blocking:   r.Blocking,
		BlockedBy:  r.BlockedBy,
		Muting:     r.Muting,
		Requested:  r.Requested,
	}
}

func (m *Mastodon) SearchAccounts(query string, limit int) ([]*shared.Account, error) {
	endpoint := endpointAccountsSearch.URL(m.opts.Server, nil)

	q := url.Values{}
	q.Add("q", query)
	q.Add("limit", strconv.Itoa(limit))

	res := []*account{}
	if err := m.request(http.MethodGet, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	accounts := []*shared.Account{}
	for _, account := range res {
		accounts = append(accounts, account.ToShared())
	}

	return accounts, nil
}

func (m *Mastodon) GetAccount(id string) (*shared.Account, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointAccounts.URL(m.opts.Server, p)

	res := &account{}
	if err := m.request(http.MethodGet, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) GetRelationships(ids []string) ([]*shared.Relationship, error) {
	endpoint := endpointRelationships.URL(m.opts.Server, nil)

	q := url.Values{}
	for _, id := range ids {
		q.Add("id[]", id)
	}

	res := []*relationship{}
	if err := m.request(http.MethodGet, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	relationships := []*shared.Relationship{}
	for _, raw := range res {
		relationships = append(relationships, raw.ToShared())
	}

	return relationships, nil
}

func (m *Mastodon) GetPosts(id string, limit int) ([]*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointAccountsStatuses.URL(m.opts.Server, p)

	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))

	res := []*status{}
	if err := m.request(http.MethodGet, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	posts := []*shared.Post{}
	for _, status := range res {
		posts = append(posts, status.ToShared())
	}

	return posts, nil
}

func (m *Mastodon) doAccountAction(id string, e shared.Endpoint) (*shared.Relationship, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := e.URL(m.opts.Server, p)

	res := &relationship{}
	if err := m.request(http.MethodPost, endpoint, nil, true, res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) Follow(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointFollow)
}

func (m *Mastodon) Unfollow(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointUnfollow)
}

func (m *Mastodon) Block(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointBlock)
}

func (m *Mastodon) Unblock(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointUnblock)
}

func (m *Mastodon) Mute(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointMute)
}

func (m *Mastodon) Unmute(id string) (*shared.Relationship, error) {
	return m.doAccountAction(id, endpointUnmute)
}
