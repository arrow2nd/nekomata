package mastodon

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/arrow2nd/nekomata/api/shared"
	"golang.org/x/net/websocket"
)

type event struct {
	Stream  []string `json:"stream"`
	Event   string
	Payload string
}

func (m *Mastodon) handleWebSocket(q url.Values, opts *shared.StreamingTimelineOpts) error {
	endpoint := endpointStreaming.URL(m.opts.Server, nil)
	url, err := shared.ConvertHttpToWebsocket(endpoint)
	if err != nil {
		return err
	}

	q.Add("access_token", m.opts.UserToken)
	url.RawQuery = q.Encode()

	conn, err := websocket.Dial(url.String(), "", m.opts.Server)
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
			return opts.Context.Err() // 終了
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

func (m *Mastodon) StreamingGlobalTimeline(opts *shared.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "public")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingLocalTimeline(opts *shared.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "public:local")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingHomeTimeline(opts *shared.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "user")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingListTimeline(opts *shared.StreamingListTimelineOpts) error {
	q := url.Values{}
	q.Add("stream", "list")
	q.Add("list", opts.ListID)

	return m.handleWebSocket(q, opts.StreamingTimelineOpts)
}
