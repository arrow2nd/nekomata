package mastodon

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// account : ユーザー情報
type account struct {
	// ID : ユーザーID
	ID string `json:"id"`
	// Username : ユーザー名
	Username string `json:"username"`
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

func (m *Mastodon) Follow(id string) error {
	p := url.Values{}
	p.Add(":id", id)

	res := &account{}
	endpoint := endpointFollow.URL(m.opts.Server, p)

	err := m.request(http.MethodPost, endpoint, nil, true, res)
	if err != nil {
		return fmt.Errorf("failed to follow (ID: %s): %w", id, err)
	}

	return nil
}

func (m *Mastodon) UnFollow(id string) error {
	p := url.Values{}
	p.Add(":id", id)

	res := &account{}
	endpoint := endpointUnfollow.URL(m.opts.Server, p)

	err := m.request(http.MethodPost, endpoint, nil, true, res)
	if err != nil {
		return fmt.Errorf("failed to unfollow (ID: %s): %w", id, err)
	}

	return nil
}
