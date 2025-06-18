package misskey

import (
	"time"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// MisskeyUser : Misskeyのユーザー情報
type MisskeyUser struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Host         *string   `json:"host"`
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	Location     *string   `json:"location"`
	Birthday     *string   `json:"birthday"`
	Lang         *string   `json:"lang"`
	AvatarURL    *string   `json:"avatarUrl"`
	BannerURL    *string   `json:"bannerUrl"`
	IsBot        bool      `json:"isBot"`
	IsCat        bool      `json:"isCat"`
	IsLocked     bool      `json:"isLocked"`
	IsExplorable bool      `json:"isExplorable"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LastFetchedAt *time.Time `json:"lastFetchedAt"`
	BannerColor  *string   `json:"bannerColor"`
	IsAdmin      bool      `json:"isAdmin"`
	IsModerator  bool      `json:"isModerator"`
	IsSilenced   bool      `json:"isSilenced"`
	IsSuspended  bool      `json:"isSuspended"`
	IsVerified   bool      `json:"isVerified"`
	Fields       []MisskeyField `json:"fields"`
	FollowersCount int     `json:"followersCount"`
	FollowingCount int     `json:"followingCount"`
	NotesCount     int     `json:"notesCount"`
}

// MisskeyField : Misskeyのプロフィールフィールド
type MisskeyField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MisskeyRelation : Misskeyでのユーザー関係
type MisskeyRelation struct {
	ID                 string `json:"id"`
	IsFollowing        bool   `json:"isFollowing"`
	HasPendingFollowRequestFromYou bool `json:"hasPendingFollowRequestFromYou"`
	HasPendingFollowRequestToYou   bool `json:"hasPendingFollowRequestToYou"`
	IsFollowed         bool   `json:"isFollowed"`
	IsBlocking         bool   `json:"isBlocking"`
	IsBlocked          bool   `json:"isBlocked"`
	IsMuted            bool   `json:"isMuted"`
	IsRenoteMuted      bool   `json:"isRenoteMuted"`
}

// SearchUsersRequest : ユーザー検索のリクエスト
type SearchUsersRequest struct {
	I      string  `json:"i"`
	Query  string  `json:"query"`
	Offset *int    `json:"offset,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Origin *string `json:"origin,omitempty"`
	Detail *bool   `json:"detail,omitempty"`
}

func (m *Misskey) SearchAccounts(query string, limit int) ([]*sharedapi.Account, error) {
	req := SearchUsersRequest{
		I:     m.user.Token,
		Query: query,
		Limit: &limit,
	}

	var users []MisskeyUser
	if err := m.post(endpointUsersSearch, req, &users); err != nil {
		return nil, err
	}

	accounts := make([]*sharedapi.Account, len(users))
	for i, user := range users {
		accounts[i] = convertToAccount(&user)
	}

	return accounts, nil
}

// UsersShowRequest : ユーザー詳細取得のリクエスト
type UsersShowRequest struct {
	I      string  `json:"i"`
	UserID *string `json:"userId,omitempty"`
	Username *string `json:"username,omitempty"`
	Host   *string `json:"host,omitempty"`
}

func (m *Misskey) GetAccount(id string) (*sharedapi.Account, error) {
	req := UsersShowRequest{
		I:      m.user.Token,
		UserID: &id,
	}

	var user MisskeyUser
	if err := m.post(endpointUsersShow, req, &user); err != nil {
		return nil, err
	}

	return convertToAccount(&user), nil
}

// IRequest : 自分の情報取得のリクエスト
type IRequest struct {
	I string `json:"i"`
}

func (m *Misskey) GetLoginAccount() (*sharedapi.Account, error) {
	req := IRequest{
		I: m.user.Token,
	}

	var user MisskeyUser
	if err := m.post(endpointI, req, &user); err != nil {
		return nil, err
	}

	return convertToAccount(&user), nil
}

// UsersRelationRequest : ユーザー関係取得のリクエスト
type UsersRelationRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

func (m *Misskey) GetRelationships(ids []string) ([]*sharedapi.Relationship, error) {
	relationships := make([]*sharedapi.Relationship, len(ids))

	for i, id := range ids {
		req := UsersRelationRequest{
			I:      m.user.Token,
			UserID: id,
		}

		var relation MisskeyRelation
		if err := m.post(endpointUsersRelation, req, &relation); err != nil {
			return nil, err
		}

		relationships[i] = convertToRelationship(&relation)
	}

	return relationships, nil
}

// UsersNotesRequest : ユーザーのノート取得のリクエスト
type UsersNotesRequest struct {
	I              string  `json:"i"`
	UserID         string  `json:"userId"`
	IncludeReplies *bool   `json:"includeReplies,omitempty"`
	Limit          *int    `json:"limit,omitempty"`
	SinceID        *string `json:"sinceId,omitempty"`
	UntilID        *string `json:"untilId,omitempty"`
	WithFiles      *bool   `json:"withFiles,omitempty"`
	ExcludeNsfw    *bool   `json:"excludeNsfw,omitempty"`
}

func (m *Misskey) GetPosts(id string, limit int) ([]*sharedapi.Post, error) {
	req := UsersNotesRequest{
		I:      m.user.Token,
		UserID: id,
		Limit:  &limit,
	}

	var notes []MisskeyNote
	if err := m.post(endpointUsersNotes, req, &notes); err != nil {
		return nil, err
	}

	posts := make([]*sharedapi.Post, len(notes))
	for i, note := range notes {
		posts[i] = convertToPost(&note)
	}

	return posts, nil
}

// FollowingCreateRequest : フォロー作成のリクエスト
type FollowingCreateRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
	WithReplies *bool `json:"withReplies,omitempty"`
}

func (m *Misskey) Follow(id string) (*sharedapi.Relationship, error) {
	req := FollowingCreateRequest{
		I:      m.user.Token,
		UserID: id,
	}

	var user MisskeyUser
	if err := m.post(endpointFollowingCreate, req, &user); err != nil {
		return nil, err
	}

	// フォロー後の関係を取得
	return m.getSingleRelationship(id)
}

// FollowingDeleteRequest : フォロー削除のリクエスト
type FollowingDeleteRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

func (m *Misskey) Unfollow(id string) (*sharedapi.Relationship, error) {
	req := FollowingDeleteRequest{
		I:      m.user.Token,
		UserID: id,
	}

	var user MisskeyUser
	if err := m.post(endpointFollowingDelete, req, &user); err != nil {
		return nil, err
	}

	// フォロー解除後の関係を取得
	return m.getSingleRelationship(id)
}

// BlockingCreateRequest : ブロック作成のリクエスト
type BlockingCreateRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

func (m *Misskey) Block(id string) (*sharedapi.Relationship, error) {
	req := BlockingCreateRequest{
		I:      m.user.Token,
		UserID: id,
	}

	var user MisskeyUser
	if err := m.post(endpointBlockingCreate, req, &user); err != nil {
		return nil, err
	}

	// ブロック後の関係を取得
	return m.getSingleRelationship(id)
}

// BlockingDeleteRequest : ブロック削除のリクエスト
type BlockingDeleteRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

func (m *Misskey) Unblock(id string) (*sharedapi.Relationship, error) {
	req := BlockingDeleteRequest{
		I:      m.user.Token,
		UserID: id,
	}

	var user MisskeyUser
	if err := m.post(endpointBlockingDelete, req, &user); err != nil {
		return nil, err
	}

	// ブロック解除後の関係を取得
	return m.getSingleRelationship(id)
}

// MuteCreateRequest : ミュート作成のリクエスト
type MuteCreateRequest struct {
	I         string     `json:"i"`
	UserID    string     `json:"userId"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

func (m *Misskey) Mute(id string) (*sharedapi.Relationship, error) {
	req := MuteCreateRequest{
		I:      m.user.Token,
		UserID: id,
	}

	if err := m.post(endpointMuteCreate, req, nil); err != nil {
		return nil, err
	}

	// ミュート後の関係を取得
	return m.getSingleRelationship(id)
}

// MuteDeleteRequest : ミュート削除のリクエスト
type MuteDeleteRequest struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

func (m *Misskey) Unmute(id string) (*sharedapi.Relationship, error) {
	req := MuteDeleteRequest{
		I:      m.user.Token,
		UserID: id,
	}

	if err := m.post(endpointMuteDelete, req, nil); err != nil {
		return nil, err
	}

	// ミュート解除後の関係を取得
	return m.getSingleRelationship(id)
}

// getSingleRelationship : 単一のユーザーとの関係を取得
func (m *Misskey) getSingleRelationship(id string) (*sharedapi.Relationship, error) {
	req := UsersRelationRequest{
		I:      m.user.Token,
		UserID: id,
	}

	var relation MisskeyRelation
	if err := m.post(endpointUsersRelation, req, &relation); err != nil {
		return nil, err
	}

	return convertToRelationship(&relation), nil
}

// convertToAccount : MisskeyUserをsharedapi.Accountに変換
func convertToAccount(user *MisskeyUser) *sharedapi.Account {
	username := user.Username
	if user.Host != nil && *user.Host != "" {
		username = user.Username + "@" + *user.Host
	}

	displayName := username
	if user.Name != nil {
		displayName = *user.Name
	}

	bio := ""
	if user.Description != nil {
		bio = *user.Description
	}

	profiles := make([]sharedapi.Profile, len(user.Fields))
	for i, field := range user.Fields {
		profiles[i] = sharedapi.Profile{
			Label: field.Name,
			Value: field.Value,
		}
	}

	return &sharedapi.Account{
		ID:             user.ID,
		Username:       username,
		DisplayName:    displayName,
		Private:        user.IsLocked,
		Bot:            user.IsBot,
		Verified:       user.IsVerified,
		BIO:            bio,
		CreatedAt:      user.CreatedAt,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		PostsCount:     user.NotesCount,
		Profiles:       profiles,
	}
}

// convertToRelationship : MisskeyRelationをsharedapi.Relationshipに変換
func convertToRelationship(relation *MisskeyRelation) *sharedapi.Relationship {
	return &sharedapi.Relationship{
		ID:         relation.ID,
		Following:  relation.IsFollowing,
		FollowedBy: relation.IsFollowed,
		Blocking:   relation.IsBlocking,
		BlockedBy:  relation.IsBlocked,
		Muting:     relation.IsMuted,
		Requested:  relation.HasPendingFollowRequestFromYou,
	}
}
