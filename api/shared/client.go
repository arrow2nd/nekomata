package shared

import "io"

type Client interface {
	// Authenticate : アプリケーション認証を行なってアクセストークンを取得
	Authenticate(io.Writer) (*User, error)
	// GetAnnouncements : サーバーからのお知らせを取得
	GetAnnouncements() ([]*Announcement, error)
	// CreatePost : 投稿を作成
	CreatePost(*CreatePostOpts) (*Post, error)
	// QuotePost : 投稿を引用
	QuotePost(string, *CreatePostOpts) (*Post, error)
	// ReplyPost : 投稿に返信
	ReplyPost(string, *CreatePostOpts) (*Post, error)
	// DeletePost : 投稿を削除
	DeletePost(string) (*Post, error)
	// Reaction : 投稿にリアクション
	Reaction(string, string) error
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
