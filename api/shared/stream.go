package shared

import "context"

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
