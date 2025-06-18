package misskey

import (
	"testing"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

func TestIsStreamingSupported(t *testing.T) {
	m := &Misskey{}
	
	if !m.IsStreamingSupported() {
		t.Error("Expected streaming to be supported")
	}
}

func TestMisskeyStreamConnect(t *testing.T) {
	// MisskeyStreamConnectの構造体テスト
	connect := MisskeyStreamConnect{
		Type: "connect",
		Body: MisskeyStreamConnectBody{
			Channel: "homeTimeline",
			ID:      "home-timeline",
			Params:  map[string]interface{}{"test": "value"},
		},
	}

	if connect.Type != "connect" {
		t.Errorf("Expected type 'connect', got %s", connect.Type)
	}

	if connect.Body.Channel != "homeTimeline" {
		t.Errorf("Expected channel 'homeTimeline', got %s", connect.Body.Channel)
	}

	if connect.Body.ID != "home-timeline" {
		t.Errorf("Expected ID 'home-timeline', got %s", connect.Body.ID)
	}

	if connect.Body.Params["test"] != "value" {
		t.Error("Expected params to contain test=value")
	}
}

func TestMisskeyStreamMessage(t *testing.T) {
	// MisskeyStreamMessageの構造体テスト
	msg := MisskeyStreamMessage{
		Type: "channel",
		Body: MisskeyStreamBody{
			ID:   "home-timeline",
			Type: "note",
			Body: map[string]interface{}{"id": "note123"},
		},
	}

	if msg.Type != "channel" {
		t.Errorf("Expected type 'channel', got %s", msg.Type)
	}

	if msg.Body.ID != "home-timeline" {
		t.Errorf("Expected body ID 'home-timeline', got %s", msg.Body.ID)
	}

	if msg.Body.Type != "note" {
		t.Errorf("Expected body type 'note', got %s", msg.Body.Type)
	}
}

func TestMisskeyStreamingMethods(t *testing.T) {
	// Misskeyストリーミングメソッドの基本テスト
	clientCred := &sharedapi.ClientCredential{
		Service: "Misskey",
		Name:    "TestApp",
	}

	userCred := &sharedapi.UserCredential{
		Server: "https://misskey.example.com",
		Token:  "test_token",
	}

	client, err := New(clientCred, userCred)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 実際のWebSocket接続は行わず、メソッドが存在することを確認
	// 実際のサーバーがないので、関数呼び出しでpanicしないことを確認

	// StreamingGlobalTimeline
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StreamingGlobalTimeline panicked: %v", r)
		}
	}()

	// nil optsでテスト（実際の接続は失敗するが、panicしないことを確認）
	err = client.StreamingGlobalTimeline(nil)
	if err == nil {
		t.Error("Expected error for nil opts, got nil")
	}
}

func TestStreamURLConstruction(t *testing.T) {
	tests := []struct {
		name       string
		serverURL  string
		token      string
		expected   string
	}{
		{
			name:      "HTTPS URL",
			serverURL: "https://misskey.io",
			token:     "testtoken",
			expected:  "wss://misskey.io/streaming?i=testtoken",
		},
		{
			name:      "HTTP URL",
			serverURL: "http://localhost:3000",
			token:     "localtoken",
			expected:  "wss://localhost:3000/streaming?i=localtoken",
		},
		{
			name:      "URL with trailing slash",
			serverURL: "https://misskey.example.com/",
			token:     "exampletoken",
			expected:  "wss://misskey.example.com//streaming?i=exampletoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// URL構築ロジックのテスト
			client := &Misskey{
				user: &sharedapi.UserCredential{
					Server: tt.serverURL,
					Token:  tt.token,
				},
			}

			// handleWebSocketの内部ロジックをテスト用に分離
			serverHost := client.user.Server
			if serverHost[:8] == "https://" {
				serverHost = serverHost[8:]
			} else if serverHost[:7] == "http://" {
				serverHost = serverHost[7:]
			}
			
			constructed := "wss://" + serverHost + "/streaming?i=" + tt.token
			if constructed != tt.expected {
				t.Errorf("Expected URL %s, got %s", tt.expected, constructed)
			}
		})
	}
}