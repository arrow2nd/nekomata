package mastodon

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (m *Mastodon) IsStreamingSupported() bool {
	return true
}

type event struct {
	Stream  []string `json:"stream"`
	Event   string
	Payload string
}

func (m *Mastodon) handleWebSocket(q url.Values, opts *sharedapi.StreamingTimelineOpts) error {
	endpoint := endpointStreaming.URL(m.user.Server, nil)
	url, err := sharedapi.ConvertHttpToWebsocket(endpoint)
	if err != nil {
		return err
	}

	q.Add("access_token", m.user.Token)
	url.RawQuery = q.Encode()

	conn, _, err := websocket.Dial(context.Background(), url.String(), nil)
	if err != nil {
		return err
	}

	go func() {
		<-opts.Context.Done()
		conn.CloseNow()
	}()

	for {
		select {
		case <-opts.Context.Done():
			return nil // 終了
		default:
		}

		var res event
		if err := wsjson.Read(opts.Context, conn, &res); err != nil {
			if !errors.Is(err, context.Canceled) {
				opts.OnError(err)
			}
			break // 再接続
		}

		switch res.Event {
		case "update":
			var status status
			if err := json.Unmarshal([]byte(res.Payload), &status); err != nil {
				opts.OnError(err)
				continue // 継続
			}
			opts.OnUpdate(status.ToShared())

		case "delete":
			opts.OnDelete(strings.TrimSpace(res.Payload))
		}
	}

	return nil
}

func (m *Mastodon) StreamingGlobalTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "public")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingLocalTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "public:local")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingHomeTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "user")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingListTimeline(opts *sharedapi.StreamingListTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "list")
	q.Add("list", opts.ListID)

	return m.handleWebSocket(q, opts.StreamingTimelineOpts)
}
