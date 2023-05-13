package misskey

import "github.com/arrow2nd/nekomata/api/shared"

const (
	endpointMiAuth        shared.Endpoint = "/miauth/:session_id"
	endpointMiAuthCheck   shared.Endpoint = "/api/miauth/:session_id/check"
	endpointAnnouncements shared.Endpoint = "/api/announcements"
)
