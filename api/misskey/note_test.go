package misskey

import (
	"testing"
	"time"
)

func TestConvertToPost(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	text := "Hello, Misskey!"
	cw := "Sensitive content"
	myReaction := "👍"

	user := MisskeyUser{
		ID:       "user123",
		Username: "testuser",
	}

	note := &MisskeyNote{
		ID:        "note123",
		CreatedAt: createdAt,
		Text:      &text,
		CW:        &cw,
		User:      user,
		UserID:    "user123",
		Tags:      []string{"test", "misskey"},
		Visibility: "public",
		Reactions: map[string]int{
			"👍": 5,
			"❤️": 3,
		},
		RenoteCount:  10,
		RepliesCount: 2,
		MyReaction:   &myReaction,
	}

	post := convertToPost(note)

	if post.ID != "note123" {
		t.Errorf("Expected ID note123, got %s", post.ID)
	}

	if post.Text != "Hello, Misskey!" {
		t.Errorf("Expected text 'Hello, Misskey!', got %s", post.Text)
	}

	if post.Visibility != "public" {
		t.Errorf("Expected visibility public, got %s", post.Visibility)
	}

	if post.Sensitive != true {
		t.Errorf("Expected sensitive true (CW is set), got %v", post.Sensitive)
	}

	if post.RepostCount != 10 {
		t.Errorf("Expected repost count 10, got %d", post.RepostCount)
	}

	if len(post.Reactions) != 2 {
		t.Errorf("Expected 2 reactions, got %d", len(post.Reactions))
	}

	// リアクションのテスト
	foundThumbsUp := false
	foundHeart := false
	for _, reaction := range post.Reactions {
		if reaction.Name == "👍" {
			foundThumbsUp = true
			if reaction.Count != 5 {
				t.Errorf("Expected thumbs up count 5, got %d", reaction.Count)
			}
			if reaction.Reacted != true {
				t.Errorf("Expected thumbs up reacted true, got %v", reaction.Reacted)
			}
		}
		if reaction.Name == "❤️" {
			foundHeart = true
			if reaction.Count != 3 {
				t.Errorf("Expected heart count 3, got %d", reaction.Count)
			}
			if reaction.Reacted != false {
				t.Errorf("Expected heart reacted false, got %v", reaction.Reacted)
			}
		}
	}

	if !foundThumbsUp {
		t.Error("Expected to find thumbs up reaction")
	}
	if !foundHeart {
		t.Error("Expected to find heart reaction")
	}

	if len(post.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(post.Tags))
	}

	if post.Tags[0].Name != "test" {
		t.Errorf("Expected first tag 'test', got %s", post.Tags[0].Name)
	}

	if post.Author == nil {
		t.Error("Expected author to be set")
	} else if post.Author.ID != "user123" {
		t.Errorf("Expected author ID user123, got %s", post.Author.ID)
	}
}

func TestConvertToPostWithNilFields(t *testing.T) {
	user := MisskeyUser{
		ID:       "user123",
		Username: "testuser",
	}

	note := &MisskeyNote{
		ID:         "note123",
		Text:       nil,
		CW:         nil,
		User:       user,
		Tags:       []string{},
		Visibility: "home",
		Reactions:  map[string]int{},
		MyReaction: nil,
	}

	post := convertToPost(note)

	if post.Text != "" {
		t.Errorf("Expected empty text, got %s", post.Text)
	}

	if post.Sensitive != false {
		t.Errorf("Expected sensitive false (no CW), got %v", post.Sensitive)
	}

	if len(post.Reactions) != 0 {
		t.Errorf("Expected 0 reactions, got %d", len(post.Reactions))
	}

	if len(post.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(post.Tags))
	}
}

func TestConvertToPostWithRenote(t *testing.T) {
	user := MisskeyUser{
		ID:       "user123",
		Username: "testuser",
	}

	renoteText := "Original note"
	renote := &MisskeyNote{
		ID:   "renote123",
		Text: &renoteText,
		User: user,
	}

	text := "Quote text"
	note := &MisskeyNote{
		ID:     "note123",
		Text:   &text,
		User:   user,
		Renote: renote,
	}

	post := convertToPost(note)

	if post.Reference == nil {
		t.Error("Expected reference to be set")
	} else {
		if post.Reference.ID != "renote123" {
			t.Errorf("Expected reference ID renote123, got %s", post.Reference.ID)
		}
		if post.Reference.Text != "Original note" {
			t.Errorf("Expected reference text 'Original note', got %s", post.Reference.Text)
		}
	}
}

func TestGetVisibilityList(t *testing.T) {
	m := &Misskey{}
	visibilities := m.GetVisibilityList()

	expected := []string{"public", "home", "followers", "specified"}
	if len(visibilities) != len(expected) {
		t.Errorf("Expected %d visibilities, got %d", len(expected), len(visibilities))
	}

	for i, expected := range expected {
		if i >= len(visibilities) || visibilities[i] != expected {
			t.Errorf("Expected visibility %s at index %d, got %s", expected, i, visibilities[i])
		}
	}
}