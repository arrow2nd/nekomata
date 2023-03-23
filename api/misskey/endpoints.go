package misskey

import "github.com/arrow2nd/nekomata/api/shared"

const (
	miAuthEndpoint        shared.Endpoint = "miauth/:session_id"
	miAuthCheckEndpoint   shared.Endpoint = "api/miauth/:session_id/check"
	announcementsEndpoint shared.Endpoint = "api/announcements"
)
