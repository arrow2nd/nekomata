package misskey

import "github.com/arrow2nd/nekomata/api"

const (
	endpointMiAuth        api.Endpoint = "/miauth/:session_id"
	endpointMiAuthCheck   api.Endpoint = "/api/miauth/:session_id/check"
	endpointAnnouncements api.Endpoint = "/api/announcements"
)
