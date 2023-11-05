package sharedapi

import "io"

type Client interface {
	// IsStreamingSupported : ストリーミングをサポートしているか
	IsStreamingSupported() bool

	// RegisterNewApplication : サーバーにアプリケーションを登録
	RegisterNewApplication() (string, string, error)
	// Authenticate : アプリケーション認証を行なってアクセストークンを取得
	Authenticate(w io.Writer) (string, error)

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

	// UploadMedia : メディアをアップロード (画像のみ対応)
	UploadMedia(filename string, src io.Reader) (string, error)

	// Reaction : 投稿にリアクション
	Reaction(id string, reactionName string) (*Post, error)
	// RemoveReaction : リアクションを削除
	RemoveReaction(id string) (*Post, error)
	// Repost : 投稿をリポスト
	Repost(id string) (*Post, error)
	// RemoveRepost : リポストを削除
	RemoveRepost(id string) (*Post, error)
	// Bookmark : 投稿をブックマーク
	Bookmark(id string) (*Post, error)
	// RemoveBookmark : ブックマークを解除
	RemoveBookmark(id string) (*Post, error)

	// SearchAccounts : アカウントを検索
	SearchAccounts(query string, limit int) ([]*Account, error)
	// GetAccount : アカウント情報を取得
	GetAccount(id string) (*Account, error)
	// GetLoginAccount : ログイン中のアカウント情報を取得
	GetLoginAccount() (*Account, error)
	// GetRelationships : ユーザーとの関係を取得
	GetRelationships(ids []string) ([]*Relationship, error)
	// GetPosts : アカウントの投稿を取得
	GetPosts(id string, limit int) ([]*Post, error)

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

	// GetGlobalTimeline : グローバルタイムラインを取得
	GetGlobalTimeline(sinceID string, limit int) ([]*Post, error)
	// GetLocalTimeline : ローカルタイムラインを取得
	GetLocalTimeline(sinceID string, limit int) ([]*Post, error)
	// GetHomeTimeline : ホームタイムラインを取得
	GetHomeTimeline(sinceID string, limit int) ([]*Post, error)
	// GetListTimeline : リストタイムラインを取得
	GetListTimeline(listID, sinceID string, limit int) ([]*Post, error)

	// StreamingGlobalTimeline : グローバルタイムラインをストリーミング
	StreamingGlobalTimeline(opts *StreamingTimelineOpts) error
	// StreamingLocalTimeline : ローカルタイムラインをストリーミング
	StreamingLocalTimeline(opts *StreamingTimelineOpts) error
	// StreamingHomeTimeline : ホームタイムラインをストリーミング
	StreamingHomeTimeline(opts *StreamingTimelineOpts) error
	// StreamingListTimeline : リストタイムラインをストリーミング
	StreamingListTimeline(opts *StreamingListTimelineOpts) error
}

// ClientCredential : クライアントの資格情報
type ClientCredential struct {
	// Service : サービス種別
	Service string
	Name    string
	ID      string
	Secret  string
}

// IsUncertified : 資格情報が無いかどうか
func (c *ClientCredential) IsUncertified() bool {
	return c.ID == "" || c.Secret == ""
}

// UserCredential : ユーザーの資格情報
type UserCredential struct {
	// Server : 接続先サーバーのURL
	Server string
	// Token : トークン
	Token string
}
