package mastodon

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"golang.org/x/net/websocket"
)

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

	conn, err := websocket.Dial(url.String(), "", m.user.Server)
	if err != nil {
		return err
	}

	go func() {
		<-opts.Context.Done()
		conn.Close()
	}()

	for {
		var res event
		err := websocket.JSON.Receive(conn, &res)

		select {
		case <-opts.Context.Done():
			return nil // 終了
		default:
		}

		if err != nil {
			opts.OnError(err)
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
