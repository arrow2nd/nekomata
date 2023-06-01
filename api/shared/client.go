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
	Reaction(id string, reactionName string) (*Post, error)
	// Unreaction : リアクションを削除
	Unreaction(id string) (*Post, error)
	// Repost : 投稿をリポスト
	Repost(id string) (*Post, error)
	// Unrepost : リポストを削除
	Unrepost(id string) (*Post, error)
	// Bookmark : 投稿をブックマーク
	Bookmark(id string) (*Post, error)
	// Unbookmark : ブックマークを解除
	Unbookmark(id string) (*Post, error)
	// SearchAccounts : アカウントを検索
	SearchAccounts(query string, limit int) ([]*Account, error)
	// GetAccount : アカウント情報を取得
	GetAccount(id string) (*Account, error)
	// GetRelationships : ユーザーとの関係を取得
	GetRelationships(ids []string) ([]*Relationship, error)
	// GetPosts : アカウントの投稿を取得
	GetPosts(id string) ([]*Post, error)
	// Follow : ユーザーをフォロー
	Follow(id string) (*Relationship, error)
	// Unfollow : ユーザーのフォローを解除
	Unfollow(id string) (*Relationship, error)
	// Block : ユーザーをブロック
	Block(id string) (*Relationship, error)
	// Unblock : ユーザーのブロックを解除
	Unblock(id string) (*Relationship, error)
	// Mute : ユーザーをミュート
	Mute(id string) (*Relationship, error)
	// Unmute : ユーザーをミュート
	Unmute(id string) (*Relationship, error)
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
