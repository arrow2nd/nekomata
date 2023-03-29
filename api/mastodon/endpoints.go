package mastodon

import "github.com/arrow2nd/nekomata/api/shared"

const (
	oauthAuthorizeEndpoint shared.Endpoint = "oauth/authorize"
	oauthTokenEndpoint     shared.Endpoint = "oauth/token"
	announcementsEndpoint  shared.Endpoint = "api/v1/announcements"
	postNewStatusEndpoint  shared.Endpoint = "api/v1/statuses"
)
