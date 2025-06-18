package misskey

import (
	"time"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// MisskeyNote : Misskeyのノート
type MisskeyNote struct {
	ID         string        `json:"id"`
	CreatedAt  time.Time     `json:"createdAt"`
	Text       *string       `json:"text"`
	CW         *string       `json:"cw"`
	User       MisskeyUser   `json:"user"`
	UserID     string        `json:"userId"`
	Reply      *MisskeyNote  `json:"reply"`
	Renote     *MisskeyNote  `json:"renote"`
	Files      []interface{} `json:"files"`
	Tags       []string      `json:"tags"`
	Visibility string        `json:"visibility"`
	Reactions  map[string]int `json:"reactions"`
	RenoteCount int          `json:"renoteCount"`
	RepliesCount int         `json:"repliesCount"`
	MyReaction *string       `json:"myReaction"`
}

// NotesCreateRequest : ノート作成のリクエスト
type NotesCreateRequest struct {
	I          string    `json:"i"`
	Visibility *string   `json:"visibility,omitempty"`
	Text       *string   `json:"text,omitempty"`
	CW         *string   `json:"cw,omitempty"`
	ViaMobile  *bool     `json:"viaMobile,omitempty"`
	LocalOnly  *bool     `json:"localOnly,omitempty"`
	FileIds    []string  `json:"fileIds,omitempty"`
	ReplyId    *string   `json:"replyId,omitempty"`
	RenoteId   *string   `json:"renoteId,omitempty"`
}

func (m *Misskey) CreatePost(opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	req := NotesCreateRequest{
		I:    m.user.Token,
		Text: &opts.Text,
	}

	if opts.Visibility != "" {
		req.Visibility = &opts.Visibility
	}

	if len(opts.MediaIDs) > 0 {
		req.FileIds = opts.MediaIDs
	}

	if opts.Sensitive {
		cw := "センシティブな内容を含む可能性があります"
		req.CW = &cw
	}

	var note MisskeyNote
	if err := m.post(endpointNotesCreate, req, &note); err != nil {
		return nil, err
	}

	return convertToPost(&note), nil
}

func (m *Misskey) QuotePost(id string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	req := NotesCreateRequest{
		I:        m.user.Token,
		Text:     &opts.Text,
		RenoteId: &id,
	}

	if opts.Visibility != "" {
		req.Visibility = &opts.Visibility
	}

	if len(opts.MediaIDs) > 0 {
		req.FileIds = opts.MediaIDs
	}

	if opts.Sensitive {
		cw := "センシティブな内容を含む可能性があります"
		req.CW = &cw
	}

	var note MisskeyNote
	if err := m.post(endpointNotesCreate, req, &note); err != nil {
		return nil, err
	}

	return convertToPost(&note), nil
}

func (m *Misskey) ReplyPost(id string, opts *sharedapi.CreatePostOpts) (*sharedapi.Post, error) {
	req := NotesCreateRequest{
		I:       m.user.Token,
		Text:    &opts.Text,
		ReplyId: &id,
	}

	if opts.Visibility != "" {
		req.Visibility = &opts.Visibility
	}

	if len(opts.MediaIDs) > 0 {
		req.FileIds = opts.MediaIDs
	}

	if opts.Sensitive {
		cw := "センシティブな内容を含む可能性があります"
		req.CW = &cw
	}

	var note MisskeyNote
	if err := m.post(endpointNotesCreate, req, &note); err != nil {
		return nil, err
	}

	return convertToPost(&note), nil
}

// NotesDeleteRequest : ノート削除のリクエスト
type NotesDeleteRequest struct {
	I      string `json:"i"`
	NoteId string `json:"noteId"`
}

func (m *Misskey) DeletePost(id string) error {
	req := NotesDeleteRequest{
		I:      m.user.Token,
		NoteId: id,
	}

	return m.post(endpointNotesDelete, req, nil)
}

// NotesReactionsCreateRequest : リアクション作成のリクエスト
type NotesReactionsCreateRequest struct {
	I        string `json:"i"`
	NoteId   string `json:"noteId"`
	Reaction string `json:"reaction"`
}

func (m *Misskey) Reaction(id, reaction string) (*sharedapi.Post, error) {
	req := NotesReactionsCreateRequest{
		I:        m.user.Token,
		NoteId:   id,
		Reaction: reaction,
	}

	if err := m.post(endpointNotesReactionsCreate, req, nil); err != nil {
		return nil, err
	}

	// リアクション後のノート情報を取得する必要があるが、今は空を返す
	return &sharedapi.Post{ID: id}, nil
}

// NotesReactionsDeleteRequest : リアクション削除のリクエスト
type NotesReactionsDeleteRequest struct {
	I      string `json:"i"`
	NoteId string `json:"noteId"`
}

func (m *Misskey) Unreaction(id string) (*sharedapi.Post, error) {
	req := NotesReactionsDeleteRequest{
		I:      m.user.Token,
		NoteId: id,
	}

	if err := m.post(endpointNotesReactionsDelete, req, nil); err != nil {
		return nil, err
	}

	return &sharedapi.Post{ID: id}, nil
}

func (m *Misskey) Repost(id string) (*sharedapi.Post, error) {
	// MisskeyではリノートはCreatePostでテキストなしでrenoteIdを指定する
	req := NotesCreateRequest{
		I:        m.user.Token,
		RenoteId: &id,
	}

	var note MisskeyNote
	if err := m.post(endpointNotesCreate, req, &note); err != nil {
		return nil, err
	}

	return convertToPost(&note), nil
}

// NotesUnrenoteRequest : リノート取り消しのリクエスト
type NotesUnrenoteRequest struct {
	I      string `json:"i"`
	NoteId string `json:"noteId"`
}

func (m *Misskey) Unrepost(id string) (*sharedapi.Post, error) {
	req := NotesUnrenoteRequest{
		I:      m.user.Token,
		NoteId: id,
	}

	if err := m.post(endpointNotesUnrenote, req, nil); err != nil {
		return nil, err
	}

	return &sharedapi.Post{ID: id}, nil
}

// NotesFavoritesCreateRequest : お気に入り作成のリクエスト
type NotesFavoritesCreateRequest struct {
	I      string `json:"i"`
	NoteId string `json:"noteId"`
}

func (m *Misskey) Bookmark(id string) (*sharedapi.Post, error) {
	req := NotesFavoritesCreateRequest{
		I:      m.user.Token,
		NoteId: id,
	}

	if err := m.post(endpointNotesFavoritesCreate, req, nil); err != nil {
		return nil, err
	}

	return &sharedapi.Post{ID: id}, nil
}

// NotesFavoritesDeleteRequest : お気に入り削除のリクエスト
type NotesFavoritesDeleteRequest struct {
	I      string `json:"i"`
	NoteId string `json:"noteId"`
}

func (m *Misskey) Unbookmark(id string) (*sharedapi.Post, error) {
	req := NotesFavoritesDeleteRequest{
		I:      m.user.Token,
		NoteId: id,
	}

	if err := m.post(endpointNotesFavoritesDelete, req, nil); err != nil {
		return nil, err
	}

	return &sharedapi.Post{ID: id}, nil
}

func (m *Misskey) GetVisibilityList() []string {
	return []string{"public", "home", "followers", "specified"}
}

// convertToPost : MisskeyNoteをsharedapi.Postに変換
func convertToPost(note *MisskeyNote) *sharedapi.Post {
	text := ""
	if note.Text != nil {
		text = *note.Text
	}

	// リアクションを変換
	reactions := make([]sharedapi.Reaction, 0, len(note.Reactions))
	for name, count := range note.Reactions {
		reacted := note.MyReaction != nil && *note.MyReaction == name
		reactions = append(reactions, sharedapi.Reaction{
			Name:    name,
			Count:   count,
			Reacted: reacted,
		})
	}

	// タグを変換
	tags := make([]sharedapi.Tag, len(note.Tags))
	for i, tag := range note.Tags {
		tags[i] = sharedapi.Tag{
			Name: tag,
			URL:  "", // MisskeyのタグURLは別途取得が必要
		}
	}

	// 引用元の変換
	var reference *sharedapi.Post
	if note.Renote != nil {
		reference = convertToPost(note.Renote)
	}

	return &sharedapi.Post{
		ID:          note.ID,
		CreatedAt:   note.CreatedAt,
		Visibility:  note.Visibility,
		Sensitive:   note.CW != nil,
		RepostCount: note.RenoteCount,
		Reactions:   reactions,
		Reposted:    false, // リポスト状態の取得は別途必要
		Bookmarked:  false, // ブックマーク状態の取得は別途必要
		Text:        text,
		Mentions:    []sharedapi.Mention{}, // メンションのパースは別途必要
		Tags:        tags,
		Via:         "", // Via情報は別途取得が必要
		Reference:   reference,
		Author:      convertToAccount(&note.User),
	}
}
