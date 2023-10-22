package sharedapi

import "time"

// CreatePostOpts : 投稿の作成に関する設定
type CreatePostOpts struct {
	// Text : 本文
	Text string
	// MediaIDs : 添付するメディアに対応するID
	MediaIDs []string
	// Visibility : 公開範囲
	Visibility string
	// Sensitive : センシティブな内容を含むか
	Sensitive bool
}

// Post : 投稿
type Post struct {
	// ID : 投稿のID
	ID string
	// CreatedAt : 投稿日
	CreatedAt time.Time
	// Visibility : 公開範囲
	Visibility string
	// Sensitive : センシティブな内容を含むか
	Sensitive bool
	// RepostCount : リポスト数
	RepostCount int
	// Reactions : リアクション
	Reactions []Reaction
	// Reacted : リアクション済か
	Reacted bool
	// Reposted : リポスト済か
	Reposted bool
	// Bookmarked : ブックマーク済か
	Bookmarked bool
	// Text : 本文
	Text string
	// Tags : ハッシュタグ
	Tags []Tag
	// Via : 投稿元
	Via string
	// Reference : 引用元
	Reference *Post
	// Author : 投稿者
	Author *Account
	// Poll : アンケート
	// Poll Poll
}

// Reaction : リアクション詳細
type Reaction struct {
	Name  string
	Count int
}

// Tag : ハッシュタグ
type Tag struct {
	Name string
	URL  string
}
