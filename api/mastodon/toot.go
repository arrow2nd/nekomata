package mastodon

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

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
	MediaAttachments []media `json:"media_attachments"`
	// Mentions : メンション先
	Mentions []mentionAttribute `json:"mentions"`
	// Tags : タグ
	Tags []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
	// Poll : アンケート
	Poll *poll `json:"poll"`
}

// mentionAttribute : メンション先のユーザー
type mentionAttribute struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
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

// ToShared : api.Post に変換
// TODO: テスト書く
func (s *status) ToShared() *sharedapi.Post {
	post := &sharedapi.Post{
		ID:          s.ID,
		CreatedAt:   s.CreatedAt,
		Visibility:  s.Visibility,
		Sensitive:   s.Sensitive,
		RepostCount: s.ReblogsCount,
		Reactions: []sharedapi.Reaction{{
			Name:    "Like",
			Count:   s.FavouritesCount,
			Reacted: s.Favourited,
		}},
		Reposted:   s.Reblogged,
		Bookmarked: s.Bookmarked,
		Text:       html2text(s.Content),
		Tags:       []sharedapi.Tag{},
		Mentions:   []sharedapi.Mention{},
		Author:     s.Account.ToShared(),
	}

	if s.Tags != nil && len(s.Tags) > 0 {
		for _, tag := range s.Tags {
			post.Tags = append(post.Tags, sharedapi.Tag{
				Name: tag.Name,
				URL:  tag.URL,
			})
		}
	}

	if s.Mentions != nil && len(s.Mentions) > 0 {
		for _, mention := range s.Mentions {
			post.Mentions = append(post.Mentions, sharedapi.Mention{
				ID:          mention.ID,
				DisplayName: mention.Username,
				Username:    mention.Acct,
			})
		}
	}

	if app := s.Application; app != nil {
		post.Via = app.Name
	}

	if s.Reblog != nil {
		post.Reference = s.Reblog.ToShared()
	}

	return post
}

func (m *Mastodon) createPostQuery(opts *sharedapi.CreatePostOpts) url.Values {
	q := url.Values{}
	q.Add("status", opts.Text)
	q.Add("visibility", opts.Visibility)

	if opts.Sensitive {
		q.Add("sensitive", "true")
	}

	if len(opts.MediaIDs) > 0 {
		q.Add("media_ids[]", strings.Join(opts.MediaIDs, ","))
	}

	return q
}

func (m *Mastodon) CreatePost(opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	q := m.createPostQuery(opts)

	requestOpts := &requestOpts{
		method: http.MethodPost,
		url:    endpointStatuses.URL(m.user.Server, nil),
		q:      q,
		isAuth: true,
	}

	res := status{}
	if err := m.request(requestOpts, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) QuotePost(id string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	// NOTE: 引用する機能が公式の手段では存在しないので実装しない
	return nil, errors.New("quote is not available on Mastodon")
}

func (m *Mastodon) ReplyPost(replyToId string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	q := m.createPostQuery(opts)
	q.Add("in_reply_to_id", replyToId)

	requestOpts := &requestOpts{
		method: http.MethodPost,
		url:    endpointStatuses.URL(m.user.Server, nil),
		q:      q,
		isAuth: true,
	}

	res := status{}
	if err := m.request(requestOpts, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) DeletePost(id string) (*sharedapi.Post, error) {
	u, err := url.JoinPath(endpointStatuses.URL(m.user.Server, nil), id)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL for quote: %w", err)
	}

	opts := &requestOpts{
		method: http.MethodDelete,
		url:    u,
		q:      nil,
		isAuth: true,
	}

	res := status{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) doTootAction(id string, e sharedapi.Endpoint) (*sharedapi.Post, error) {
	p := url.Values{}
	p.Add(":id", id)

	opts := &requestOpts{
		method: http.MethodPost,
		url:    e.URL(m.user.Server, p),
		q:      nil,
		isAuth: true,
	}

	res := status{}
	if err := m.request(opts, &res); err != nil {
		return nil, err
	}

	return res.ToShared(), nil
}

func (m *Mastodon) Reaction(id, reaction string) (*sharedapi.Post, error) {
	return m.doTootAction(id, endpointFavourite)
}

func (m *Mastodon) Unreaction(id string) (*sharedapi.Post, error) {
	return m.doTootAction(id, endpointUnfavourite)
}

func (m *Mastodon) Repost(id string) (*sharedapi.Post, error) {
	post, err := m.doTootAction(id, endpointReblog)
	if err != nil {
		return nil, err
	}
	return post.Reference, nil
}

func (m *Mastodon) Unrepost(id string) (*sharedapi.Post, error) {
	return m.doTootAction(id, endpointUnreblog)
}

func (m *Mastodon) Bookmark(id string) (*sharedapi.Post, error) {
	return m.doTootAction(id, endpointBookmark)
}

func (m *Mastodon) Unbookmark(id string) (*sharedapi.Post, error) {
	return m.doTootAction(id, endpointUnbookmark)
}
