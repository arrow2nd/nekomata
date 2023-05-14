package mastodon

import "github.com/arrow2nd/nekomata/api/shared"

const (
	endpointOauthAuthorize shared.Endpoint = "/oauth/authorize"
	endpointOauthToken     shared.Endpoint = "/oauth/token"
	endpointAnnouncements  shared.Endpoint = "/api/v1/announcements"
	endpointStatuses       shared.Endpoint = "/api/v1/statuses"
	endpointFavourite      shared.Endpoint = "/api/v1/statuses/:id/favourite"
	endpointUnfavourite    shared.Endpoint = "/api/v1/statuses/:id/unfavourite"
	endpointReblog         shared.Endpoint = "/api/v1/statuses/:id/reblog"
	endpointUnreblog       shared.Endpoint = "/api/v1/statuses/:id/unreblog"
	endpointBookmark       shared.Endpoint = "/api/v1/statuses/:id/bookmark"
	endpointUnbookmark     shared.Endpoint = "/api/v1/statuses/:id/unbookmark"
	endpointFollow         shared.Endpoint = "/api/v1/accounts/:id/follow"
	endpointUnfollow       shared.Endpoint = "/api/v1/accounts/:id/unfollow"
	endpointBlock          shared.Endpoint = "/api/v1/accounts/:id/block"
	endpointUnblock        shared.Endpoint = "/api/v1/accounts/:id/unblock"
	endpointMute           shared.Endpoint = "/api/v1/accounts/:id/mute"
	endpointUnmute         shared.Endpoint = "/api/v1/accounts/:id/unmute"
)
