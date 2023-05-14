package mastodon

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"jaytaylor.com/html2text"
)

// mentionAttribute : メンション先のユーザー
type mentionAttribute struct {
	ID   string `json:"id"`
	Acct string `json:"acct"`
}

// poll : アンケート
type poll struct {
	// ID : アンケートID
	ID string `json:"id"`
	// ExpiresAt : 終了日時
	ExpiresAt *time.Time `json:"expires_at"`
	// Expired : 終了しているか
	Expired bool `json:"expired"`
	// Multiple : 複数投票可能か
	Multiple bool `json:"multiple"`
	// VotesCount : 総投票数
	VotesCount int `json:"votes_count"`
	// VotersCount : 総投票ユーザー数
	VotersCount *int `json:"voters_count"`
	// Voted : 投票済か
	Voted bool `json:"voted"`
	// OwnVotes : 自分が投票した項目
	OwnVotes []int `json:"own_votes"`
	// Options : アンケート項目
	Options []pollOption `json:"options"`
}

// pollOption : アンケート項目
type pollOption struct {
	// Title : タイトル
	Title string `json:"title"`
	// VotersCount : 投票数
	VotesCount int `json:"votes_count"`
}

// mediaAttachment : 添付メディア
type mediaAttachment struct {
	// ID : メディアID
	ID string `json:"id"`
	// Type : 種類 (unknown, image, gifv, video, audio)
	Type string `json:"type"`
	// URL : オリジナルのメディアを指すURL
	URL string `json:"url"`
	// PreviewURL : スケールダウンされたメディアを指すURL
	PreviewURL string `json:"preview_url"`
	// Meta : メタ情報
	Meta mediaMeta `json:"meta"`
}

// mediaMeta : メディアのメタ情報
type mediaMeta struct {
	// Original : オリジナルサイズ
	Original mediaSize `json:"original"`
	// Small : 縮小サイズ
	Small mediaSize `json:"small"`
}

// mediaSize : メディアサイズ
type mediaSize struct {
	// Width : 幅
	Width int `json:"width"`
	// Height : 高さ
	Height int `json:"height"`
}

// status : 投稿
type status struct {
	// ID : 投稿ID
	ID string `json:"id"`
	// CreatedAt : 投稿日時
	CreatedAt time.Time `json:"created_at"`
	// Sensitive : センシティブな内容を含むか
	Sensitive bool `json:"sensitive"`
	// Visibility : 公開範囲 (public, unlisted, private, direct)
	Visibility string `json:"visibility"`
	// ReblogsCount : ブースト数
	ReblogsCount int `json:"reblogs_count"`
	// FavouritesCount : いいね数
	FavouritesCount int `json:"favourites_count"`
	// Favourited : いいね済か
	Favourited bool `json:"favourited"`
	// Reblogged : ブースト済か
	Reblogged bool `json:"reblogged"`
	// Bookmarked : ブックマーク済か
	Bookmarked bool `json:"bookmarked"`
	// Content : 本文 (HTML)
	Content string `json:"content"`
	// SpoilerText : 要約文 (これより下の文章は展開するまで折りたたまれる)
	SpoilerText string `json:"spoiler_text"`
	// Reblog : 引用元
	Reblog *status `json:"reblog"`
	// Application : 投稿元のクライアント
	Application *struct {
		Name string `json:"name"`
	} `json:"application"`
	// Account : 投稿ユーザー
	Account account `json:"account"`
	// MediaAttachments : 添付メディア
	MediaAttachments []mediaAttachment `json:"media_attachments"`
	// Mentions : メンション先
	Mentions []mentionAttribute `json:"mentions"`
	// Tags : タグ
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
	// Poll : アンケート
	Poll *poll `json:"poll"`
}

// ToShared : shared.Post に変換
// TODO: テスト書く
func (s *status) ToShared() *shared.Post {
	text, err := html2text.FromString(s.Content)
	if err != nil {
		text = fmt.Sprintf("convert error: %s", err.Error())
	}

	post := &shared.Post{
		ID:          s.ID,
		CreatedAt:   s.CreatedAt,
		Visibility:  s.Visibility,
		Sensitive:   s.Sensitive,
		RepostCount: s.ReblogsCount,
		Reactions:   []shared.Reaction{{Name: "Fav", Count: s.FavouritesCount}},
		Reacted:     s.Favourited,
		Reposted:    s.Reblogged,
		Bookmarked:  s.Bookmarked,
		Text:        text,
		Via:         s.Application.Name,
	}

	if s.Reblog != nil {
		post.Reference = s.Reblog.ToShared()
	}

	return post
}

func (m *Mastodon) createPostQuery(opts *shared.CreatePostOpts) url.Values {
	q := url.Values{}
	q.Add("status", opts.Text)
	q.Add("visibility", opts.Visibility)

	if opts.Sensitive {
		q.Add("sensitive", "true")
	}

	// TODO: 未検証。もしかするとカンマ区切りじゃないかも
	if len(opts.MediaIDs) > 0 {
		q.Add("media_ids", strings.Join(opts.MediaIDs, ","))
	}

	return q
}

func (m *Mastodon) CreatePost(opts *shared.CreatePostOpts) (*shared.Post, error) {
	endpoint := endpointStatuses.URL(m.opts.Server, nil)
	q := m.createPostQuery(opts)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) QuotePost(id string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	// NOTE: 引用する機能が公式の手段では存在しないので実装しない
	return nil, errors.New("quote is not available on Mastodon")
}

func (m *Mastodon) ReplyPost(replyToId string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	endpoint := endpointStatuses.URL(m.opts.Server, nil)

	q := m.createPostQuery(opts)
	q.Add("in_reply_to_id", replyToId)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, q, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) DeletePost(id string) (*shared.Post, error) {
	u, err := url.JoinPath(endpointStatuses.URL(m.opts.Server, nil), id)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL for quote: %w", err)
	}

	res := &status{}
	if err := m.request(http.MethodDelete, u, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) Reaction(id, reaction string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointFavourite.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) UnReaction(id string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointUnfavourite.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) Repost(id string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointReblog.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) UnRepost(id string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointUnreblog.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) Bookmark(id string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointBookmark.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) UnBookmark(id string) (*shared.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := endpointUnbookmark.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request(http.MethodPost, endpoint, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}
