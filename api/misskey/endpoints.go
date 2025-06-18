package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

const (
	endpointMiAuth        sharedapi.Endpoint = "/miauth/:session_id"
	endpointMiAuthCheck   sharedapi.Endpoint = "/api/miauth/:session_id/check"
	endpointAnnouncements sharedapi.Endpoint = "/api/announcements"

	// アカウント関連
	endpointI              sharedapi.Endpoint = "/api/i"
	endpointUsersShow      sharedapi.Endpoint = "/api/users/show"
	endpointUsersSearch    sharedapi.Endpoint = "/api/users/search"
	endpointUsersNotes     sharedapi.Endpoint = "/api/users/notes"
	endpointUsersRelation  sharedapi.Endpoint = "/api/users/relation"
	endpointFollowingCreate sharedapi.Endpoint = "/api/following/create"
	endpointFollowingDelete sharedapi.Endpoint = "/api/following/delete"
	endpointBlockingCreate  sharedapi.Endpoint = "/api/blocking/create"
	endpointBlockingDelete  sharedapi.Endpoint = "/api/blocking/delete"
	endpointMuteCreate      sharedapi.Endpoint = "/api/mute/create"
	endpointMuteDelete      sharedapi.Endpoint = "/api/mute/delete"

	// ノート（投稿）関連
	endpointNotesCreate     sharedapi.Endpoint = "/api/notes/create"
	endpointNotesDelete     sharedapi.Endpoint = "/api/notes/delete"
	endpointNotesReactions  sharedapi.Endpoint = "/api/notes/reactions"
	endpointNotesReactionsCreate sharedapi.Endpoint = "/api/notes/reactions/create"
	endpointNotesReactionsDelete sharedapi.Endpoint = "/api/notes/reactions/delete"
	endpointNotesUnrenote   sharedapi.Endpoint = "/api/notes/unrenote"
	endpointNotesFavoritesCreate sharedapi.Endpoint = "/api/notes/favorites/create"
	endpointNotesFavoritesDelete sharedapi.Endpoint = "/api/notes/favorites/delete"

	// メディア関連
	endpointDriveFilesCreate sharedapi.Endpoint = "/api/drive/files/create"

	// タイムライン関連
	endpointNotesTimeline       sharedapi.Endpoint = "/api/notes/timeline"
	endpointNotesLocalTimeline  sharedapi.Endpoint = "/api/notes/local-timeline"
	endpointNotesGlobalTimeline sharedapi.Endpoint = "/api/notes/global-timeline"
	endpointNotesUserListTimeline sharedapi.Endpoint = "/api/notes/user-list-timeline"
)
