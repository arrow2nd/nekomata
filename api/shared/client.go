package shared

import "io"

type Client interface {
	// Authenticate : アプリケーション認証を行なってアクセストークンを取得
	Authenticate(w io.Writer) (*User, error)
	// GetAnnouncements : サーバーからのお知らせを取得
	GetAnnouncements() ([]*Announcement, error)
	// CreatePost : 投稿を作成
	CreatePost(opts *CreatePostOpts) (*Post, error)
	// QuotePost : 投稿を引用
	QuotePost(text string, opts *CreatePostOpts) (*Post, error)
	// ReplyPost : 投稿に返信
	ReplyPost(text string, opts *CreatePostOpts) (*Post, error)
	// DeletePost : 投稿を削除
	DeletePost(id string) (*Post, error)
	// Reaction : 投稿にリアクション
	Reaction(id string, reactionName string) error
	// UnReaction : リアクションを削除
	UnReaction(id string) error
	// Repost : 投稿をリポスト
	Repost(id string) error
	// UnRepost : リポストを削除
	UnRepost(id string) error
	// Bookbart : 投稿をブックマーク
	Bookmark(id string) error
	// UnBookmark : ブックマークを解除
	UnBookmark(id string) error
	// Follow : ユーザーをフォロー
	Follow(id string) error
	// UnFollow : ユーザーのフォローを解除
	UnFollow(id string) error
}

// ClientOpts : クライアントの設定
type ClientOpts struct {
	// Server : 接続先のサーバーのURL
	Server string
	// Name : クライアント名
	Name string
	// ID : クライアント ID
	ID string
	// Secret : クライアントシークレット
	Secret string
	// UserToken : ユーザーのアクセストークン
	UserToken string
}
