package mastodon

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"jaytaylor.com/html2text"
)

type mentionAttribute struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Acct     string `json:"acct"`
}

type poll struct {
	ID          string       `json:"id"`
	ExpiresAt   *time.Time   `json:"expires_at"`
	Expired     bool         `json:"expired"`
	Multiple    bool         `json:"multiple"`
	VotesCount  int          `json:"votes_count"`
	VotersCount *int         `json:"voters_count"`
	Voted       bool         `json:"voted"`
	OwnVotes    []int        `json:"own_votes"`
	Options     []pollOption `json:"options"`
}

type pollOption struct {
	Title      string `json:"title"`
	VotesCount int    `json:"votes_count"`
}

type mediaAttachment struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	URL        string    `json:"url"`
	PreviewURL string    `json:"preview_url"`
	TextURL    string    `json:"text_url"`
	Meta       mediaMeta `json:"meta"`
}

type mediaMeta struct {
	Original mediaSize `json:"original"`
	Small    mediaSize `json:"small"`
}

type mediaSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type status struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Sensitive       bool      `json:"sensitive"`
	Visibility      string    `json:"visibility"`
	ReblogsCount    int       `json:"reblogs_count"`
	FavouritesCount int       `json:"favourites_count"`
	Favourited      bool      `json:"favourited"`
	Reblogged       bool      `json:"reblogged"`
	Bookmarked      bool      `json:"bookmarked"`
	Content         string    `json:"content"`
	Reblog          *status   `json:"reblog"`
	Application     struct {
		Name string `json:"name"`
	} `json:"application"`
	Account          account            `json:"account"`
	MediaAttachments []mediaAttachment  `json:"media_attachments"`
	Mentions         []mentionAttribute `json:"mentions"`
	Tags             []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Poll *poll `json:"poll"`
}

// ToPost : shared.Post に変換
func (s *status) ToPost() *shared.Post {
	text, err := html2text.FromString(s.Content)
	if err != nil {
		text = fmt.Sprintf("convert error: %s", err)
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
		post.Reference = s.Reblog.ToPost()
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
	endpoint := statusesEndpoint.URL(m.opts.Server, nil)
	q := m.createPostQuery(opts)

	res := &status{}
	if err := m.request("POST", endpoint, q, true, &res); err != nil {
		return nil, err
	}

	return res.ToPost(), nil
}

func (m *Mastodon) QuotePost(id string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	// NOTE: 引用する機能が公式の手段では存在しないので実装しない
	return nil, errors.New("quote is not available on Mastodon")
}

func (m *Mastodon) ReplyPost(replyToId string, opts *shared.CreatePostOpts) (*shared.Post, error) {
	endpoint := statusesEndpoint.URL(m.opts.Server, nil)

	q := m.createPostQuery(opts)
	q.Add("in_reply_to_id", replyToId)

	res := &status{}
	if err := m.request("POST", endpoint, q, true, &res); err != nil {
		return nil, err
	}

	return res.ToPost(), nil
}

func (m *Mastodon) DeletePost(id string) (*shared.Post, error) {
	u, err := url.JoinPath(statusesEndpoint.URL(m.opts.Server, nil), id)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL for quote: %w", err)
	}

	res := &status{}
	if err := m.request("DELETE", u, nil, true, &res); err != nil {
		return nil, err
	}

	return res.ToPost(), nil
}

func (m *Mastodon) Reaction(id, reaction string) error {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := favouriteEndpoint.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request("POST", endpoint, nil, true, &res); err != nil {
		return err
	}

	if !res.Favourited {
		return fmt.Errorf("failed to favourite (ID: %s)", id)
	}

	return nil
}

func (m *Mastodon) UnReaction(id, reaction string) error {
	p := url.Values{}
	p.Add(":id", id)

	endpoint := unfavouriteEndpoint.URL(m.opts.Server, p)

	res := &status{}
	if err := m.request("POST", endpoint, nil, true, &res); err != nil {
		return err
	}

	if res.Favourited {
		return fmt.Errorf("failed to unfavourite (ID: %s)", id)
	}

	return nil
}
