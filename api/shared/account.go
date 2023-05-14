package shared

import "time"

// Account : アカウント
type Account struct {
	// ID : ユーザーID
	ID string
	// Username : ユーザー名 + ドメイン名
	Username string
	// DisplayName : 表示名
	DisplayName string
	// Private : 非公開アカウントか
	Private bool
	// Bot : ボットアカウントか
	Bot bool
	// Verified : 認証アカウントか
	Verified bool
	// BIO : 自己紹介文
	BIO string
	// CreatedAt : アカウント作成日
	CreatedAt time.Time
	// FollowersCount : フォロワー数
	FollowersCount int
	// FollowingCount : フォロイー数
	FollowingCount int
	// PostsCount : 投稿数
	PostsCount int
	// Profiles : プロフィール一覧
	Profiles []Profile
}

// Profile : プロフィール
type Profile struct {
	Label string
	Value string
}
