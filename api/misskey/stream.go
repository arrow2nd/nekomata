package misskey

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (m *Misskey) IsStreamingSupported() bool {
	return true
}

// MisskeyStreamMessage : Misskeyストリーミングメッセージ
type MisskeyStreamMessage struct {
	Type string                 `json:"type"`
	Body MisskeyStreamBody      `json:"body"`
}

// MisskeyStreamBody : Misskeyストリーミングボディ
type MisskeyStreamBody struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

// MisskeyStreamConnect : Misskeyストリーミング接続メッセージ
type MisskeyStreamConnect struct {
	Type string                    `json:"type"`
	Body MisskeyStreamConnectBody  `json:"body"`
}

// MisskeyStreamConnectBody : Misskeyストリーミング接続ボディ
type MisskeyStreamConnectBody struct {
	Channel string                 `json:"channel"`
	ID      string                 `json:"id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

func (m *Misskey) handleWebSocket(channel string, channelID string, params map[string]interface{}, opts *sharedapi.StreamingTimelineOpts) error {
	// WebSocketのURLを構築
	serverHost := strings.TrimPrefix(m.user.Server, "https://")
	serverHost = strings.TrimPrefix(serverHost, "http://")
	streamingURL := fmt.Sprintf("wss://%s/streaming?i=%s", 
		serverHost,
		m.user.Token)

	// WebSocket接続
	conn, _, err := websocket.Dial(context.Background(), streamingURL, nil)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	// 終了時の処理
	go func() {
		<-opts.Context.Done()
		conn.CloseNow()
	}()

	// チャンネルに接続
	connectMsg := MisskeyStreamConnect{
		Type: "connect",
		Body: MisskeyStreamConnectBody{
			Channel: channel,
			ID:      channelID,
			Params:  params,
		},
	}

	if err := wsjson.Write(opts.Context, conn, connectMsg); err != nil {
		return fmt.Errorf("websocket write error: %w", err)
	}

	// メッセージ受信ループ
	for {
		select {
		case <-opts.Context.Done():
			return nil
		default:
		}

		var msg MisskeyStreamMessage
		if err := wsjson.Read(opts.Context, conn, &msg); err != nil {
			if !errors.Is(err, context.Canceled) {
				opts.OnError(err)
			}
			break
		}

		// チャンネルメッセージの処理
		if msg.Type == "channel" && msg.Body.ID == channelID {
			switch msg.Body.Type {
			case "note":
				// ノート受信時の処理
				noteData, err := json.Marshal(msg.Body.Body)
				if err != nil {
					opts.OnError(err)
					continue
				}

				var note MisskeyNote
				if err := json.Unmarshal(noteData, &note); err != nil {
					opts.OnError(err)
					continue
				}

				opts.OnUpdate(convertToPost(&note))

			case "noteDeleted":
				// ノート削除時の処理
				if deleteData, ok := msg.Body.Body.(map[string]interface{}); ok {
					if noteID, ok := deleteData["noteId"].(string); ok {
						opts.OnDelete(noteID)
					}
				}
			}
		}
	}

	return nil
}

func (m *Misskey) StreamingGlobalTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	return m.handleWebSocket("globalTimeline", "global-timeline", nil, opts)
}

func (m *Misskey) StreamingLocalTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	return m.handleWebSocket("localTimeline", "local-timeline", nil, opts)
}

func (m *Misskey) StreamingHomeTimeline(opts *sharedapi.StreamingTimelineOpts) error {
	return m.handleWebSocket("homeTimeline", "home-timeline", nil, opts)
}

func (m *Misskey) StreamingListTimeline(opts *sharedapi.StreamingListTimelineOpts) error {
	params := map[string]interface{}{
		"listId": opts.ListID,
	}
	return m.handleWebSocket("userList", "list-timeline", params, opts.StreamingTimelineOpts)
}
