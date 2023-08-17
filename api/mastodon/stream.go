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

	url.RawQuery = q.Encode()

	ws, err := websocket.Dial(url.String(), "", m.opts.Server)
	if err != nil {
		return err
	}

	go func() {
		defer ws.Close()

		var res event

		for {
			select {
			case <-opts.Context.Done():
				// 終了
				return
			default:
			}

			if err := websocket.JSON.Receive(ws, &res); err != nil {
				opts.OnError(err)
				return
			}

			switch res.Event {
			case "update":
				var status status
				if err := json.Unmarshal([]byte(res.Payload), &status); err != nil {
					opts.OnError(err)
					return
				}
				opts.OnUpdate(status.ToShared())

			case "delete":
				opts.OnDelete(strings.TrimSpace(res.Payload))
			}
		}
	}()

	return nil
}

func (m *Mastodon) StreamingGlobalTimeline(opts *shared.StreamingTimelineOpts) error {
	q := url.Values{}
	q.Add("access_token", m.opts.UserToken)
	q.Add("stream", "public")

	return m.handleWebSocket(q, opts)
}

func (m *Mastodon) StreamingLocalTimeline(opts *shared.StreamingTimelineOpts) error {
	return nil
}

func (m *Mastodon) StreamingHomeTimeline(opts *shared.StreamingTimelineOpts) error {
	return nil
}

func (m *Mastodon) StreamingListTimeline(opts *shared.StreamingListTimelineOpts) error {
	return nil
}
