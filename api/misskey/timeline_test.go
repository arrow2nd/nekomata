package misskey

import (
	"testing"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

func TestTimelineRequestStructure(t *testing.T) {
	// TimelineRequestの構造体が正しく定義されているかテスト
	req := TimelineRequest{
		I: "test_token",
	}

	if req.I != "test_token" {
		t.Errorf("Expected token test_token, got %s", req.I)
	}

	// ポインタフィールドのテスト
	limit := 20
	sinceId := "since123"
	withFiles := true

	req.Limit = &limit
	req.SinceId = &sinceId
	req.WithFiles = &withFiles

	if *req.Limit != 20 {
		t.Errorf("Expected limit 20, got %d", *req.Limit)
	}

	if *req.SinceId != "since123" {
		t.Errorf("Expected sinceId since123, got %s", *req.SinceId)
	}

	if *req.WithFiles != true {
		t.Errorf("Expected withFiles true, got %v", *req.WithFiles)
	}
}

func TestHomeTimelineRequestStructure(t *testing.T) {
	// HomeTimelineRequestの構造体が正しく定義されているかテスト
	req := HomeTimelineRequest{
		I: "test_token",
	}

	if req.I != "test_token" {
		t.Errorf("Expected token test_token, got %s", req.I)
	}

	// ポインタフィールドのテスト
	limit := 30
	includeMyRenotes := false
	withRenotes := true

	req.Limit = &limit
	req.IncludeMyRenotes = &includeMyRenotes
	req.WithRenotes = &withRenotes

	if *req.Limit != 30 {
		t.Errorf("Expected limit 30, got %d", *req.Limit)
	}

	if *req.IncludeMyRenotes != false {
		t.Errorf("Expected includeMyRenotes false, got %v", *req.IncludeMyRenotes)
	}

	if *req.WithRenotes != true {
		t.Errorf("Expected withRenotes true, got %v", *req.WithRenotes)
	}
}

func TestUserListTimelineRequestStructure(t *testing.T) {
	// UserListTimelineRequestの構造体が正しく定義されているかテスト
	req := UserListTimelineRequest{
		I:      "test_token",
		ListId: "list123",
	}

	if req.I != "test_token" {
		t.Errorf("Expected token test_token, got %s", req.I)
	}

	if req.ListId != "list123" {
		t.Errorf("Expected listId list123, got %s", req.ListId)
	}

	// ポインタフィールドのテスト
	limit := 25
	includeLocalRenotes := true

	req.Limit = &limit
	req.IncludeLocalRenotes = &includeLocalRenotes

	if *req.Limit != 25 {
		t.Errorf("Expected limit 25, got %d", *req.Limit)
	}

	if *req.IncludeLocalRenotes != true {
		t.Errorf("Expected includeLocalRenotes true, got %v", *req.IncludeLocalRenotes)
	}
}

func TestMisskeyTimelineMethods(t *testing.T) {
	// Misskeyクライアントのタイムラインメソッドの基本的なテスト
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

	// 実際のAPIコールは行わず、メソッドが存在することを確認
	// GetGlobalTimeline
	posts, err := client.GetGlobalTimeline("", 10)
	// 実際のサーバーがないのでエラーが返るのは正常
	if posts == nil && err == nil {
		t.Error("Expected either posts or error to be non-nil")
	}

	// GetLocalTimeline
	posts, err = client.GetLocalTimeline("since123", 20)
	if posts == nil && err == nil {
		t.Error("Expected either posts or error to be non-nil")
	}

	// GetHomeTimeline
	posts, err = client.GetHomeTimeline("", 15)
	if posts == nil && err == nil {
		t.Error("Expected either posts or error to be non-nil")
	}

	// GetListTimeline
	posts, err = client.GetListTimeline("list123", "since456", 25)
	if posts == nil && err == nil {
		t.Error("Expected either posts or error to be non-nil")
	}
}