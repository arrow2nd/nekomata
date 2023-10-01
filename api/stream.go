package api

import (
	"context"
	"net/url"
)

// StreamingTimelineOpts : TLストリーミング関数の共通オプション
type StreamingTimelineOpts struct {
	Context  context.Context
	OnUpdate func(post *Post)
	OnDelete func(id string)
	OnError  func(err error)
}

// StreamingListTimelineOpts : リストTLストリーミング関数のオプション
type StreamingListTimelineOpts struct {
	*StreamingTimelineOpts
	ListID string
}

// ConvertHttpToWebsocket : スキーマをWSに変換
func ConvertHttpToWebsocket(rawUrl string) (*url.URL, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	return u, nil
}
