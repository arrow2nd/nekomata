package mastodon

import (
	"github.com/arrow2nd/nekomata/api/sharedapi"
)

const (
	endpointOauthAuthorize    sharedapi.Endpoint = "/oauth/authorize"
	endpointOauthToken        sharedapi.Endpoint = "/oauth/token"
	endpointApps              sharedapi.Endpoint = "/api/v1/apps"
	endpointAnnouncements     sharedapi.Endpoint = "/api/v1/announcements"
	endpointStatuses          sharedapi.Endpoint = "/api/v1/statuses"
	endpointFavourite         sharedapi.Endpoint = "/api/v1/statuses/:id/favourite"
	endpointUnfavourite       sharedapi.Endpoint = "/api/v1/statuses/:id/unfavourite"
	endpointReblog            sharedapi.Endpoint = "/api/v1/statuses/:id/reblog"
	endpointUnreblog          sharedapi.Endpoint = "/api/v1/statuses/:id/unreblog"
	endpointBookmark          sharedapi.Endpoint = "/api/v1/statuses/:id/bookmark"
	endpointUnbookmark        sharedapi.Endpoint = "/api/v1/statuses/:id/unbookmark"
	endpointAccounts          sharedapi.Endpoint = "/api/v1/accounts/:id"
	endpointVerifyCredentials sharedapi.Endpoint = "/api/v1/accounts/verify_credentials"
	endpointRelationships     sharedapi.Endpoint = "/api/v1/accounts/relationships"
	endpointAccountsSearch    sharedapi.Endpoint = "/api/v1/accounts/search"
	endpointAccountsStatuses  sharedapi.Endpoint = "/api/v1/accounts/:id/statuses"
	endpointFollow            sharedapi.Endpoint = "/api/v1/accounts/:id/follow"
	endpointUnfollow          sharedapi.Endpoint = "/api/v1/accounts/:id/unfollow"
	endpointBlock             sharedapi.Endpoint = "/api/v1/accounts/:id/block"
	endpointUnblock           sharedapi.Endpoint = "/api/v1/accounts/:id/unblock"
	endpointMute              sharedapi.Endpoint = "/api/v1/accounts/:id/mute"
	endpointUnmute            sharedapi.Endpoint = "/api/v1/accounts/:id/unmute"
	endpointTimelinePublic    sharedapi.Endpoint = "/api/v1/timelines/public"
	endpointTimelineHome      sharedapi.Endpoint = "/api/v1/timelines/home"
	endpointTimelineList      sharedapi.Endpoint = "/api/v1/timelines/list/:list_id"
	endpointStreaming         sharedapi.Endpoint = "/api/v1/streaming"
	endpointMedia             sharedapi.Endpoint = "/api/v1/media/:id"
	endpointMediaUploadAsync  sharedapi.Endpoint = "/api/v2/media"
)
