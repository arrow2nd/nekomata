package misskey

import "github.com/arrow2nd/nekomata/api/sharedapi"

const (
	endpointMiAuth        sharedapi.Endpoint = "/miauth/:session_id"
	endpointMiAuthCheck   sharedapi.Endpoint = "/api/miauth/:session_id/check"
	endpointAnnouncements sharedapi.Endpoint = "/api/announcements"
)
