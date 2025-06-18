package misskey

import (
	"testing"
	"time"
)

func TestConvertToAccount(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	host := "example.com"
	name := "Test User"
	description := "Test description"

	user := &MisskeyUser{
		ID:             "test123",
		Username:       "testuser",
		Host:           &host,
		Name:           &name,
		Description:    &description,
		IsBot:          true,
		IsLocked:       false,
		IsVerified:     true,
		CreatedAt:      createdAt,
		FollowersCount: 100,
		FollowingCount: 50,
		NotesCount:     200,
		Fields: []MisskeyField{
			{Name: "Website", Value: "https://example.com"},
			{Name: "Location", Value: "Tokyo"},
		},
	}

	account := convertToAccount(user)

	if account.ID != "test123" {
		t.Errorf("Expected ID test123, got %s", account.ID)
	}

	if account.Username != "testuser@example.com" {
		t.Errorf("Expected username testuser@example.com, got %s", account.Username)
	}

	if account.DisplayName != "Test User" {
		t.Errorf("Expected display name Test User, got %s", account.DisplayName)
	}

	if account.Bot != true {
		t.Errorf("Expected bot true, got %v", account.Bot)
	}

	if account.Private != false {
		t.Errorf("Expected private false, got %v", account.Private)
	}

	if account.Verified != true {
		t.Errorf("Expected verified true, got %v", account.Verified)
	}

	if account.BIO != "Test description" {
		t.Errorf("Expected bio Test description, got %s", account.BIO)
	}

	if account.FollowersCount != 100 {
		t.Errorf("Expected followers count 100, got %d", account.FollowersCount)
	}

	if account.FollowingCount != 50 {
		t.Errorf("Expected following count 50, got %d", account.FollowingCount)
	}

	if account.PostsCount != 200 {
		t.Errorf("Expected posts count 200, got %d", account.PostsCount)
	}

	if len(account.Profiles) != 2 {
		t.Errorf("Expected 2 profile fields, got %d", len(account.Profiles))
	}

	if account.Profiles[0].Label != "Website" || account.Profiles[0].Value != "https://example.com" {
		t.Errorf("Profile field 0 mismatch")
	}
}

func TestConvertToRelationship(t *testing.T) {
	relation := &MisskeyRelation{
		ID:                                "test123",
		IsFollowing:                       true,
		HasPendingFollowRequestFromYou:    false,
		HasPendingFollowRequestToYou:      true,
		IsFollowed:                        false,
		IsBlocking:                        true,
		IsBlocked:                         false,
		IsMuted:                           true,
		IsRenoteMuted:                     false,
	}

	relationship := convertToRelationship(relation)

	if relationship.ID != "test123" {
		t.Errorf("Expected ID test123, got %s", relationship.ID)
	}

	if relationship.Following != true {
		t.Errorf("Expected following true, got %v", relationship.Following)
	}

	if relationship.FollowedBy != false {
		t.Errorf("Expected followed by false, got %v", relationship.FollowedBy)
	}

	if relationship.Blocking != true {
		t.Errorf("Expected blocking true, got %v", relationship.Blocking)
	}

	if relationship.BlockedBy != false {
		t.Errorf("Expected blocked by false, got %v", relationship.BlockedBy)
	}

	if relationship.Muting != true {
		t.Errorf("Expected muting true, got %v", relationship.Muting)
	}

	if relationship.Requested != false {
		t.Errorf("Expected requested false, got %v", relationship.Requested)
	}
}

func TestConvertToAccountWithNilFields(t *testing.T) {
	user := &MisskeyUser{
		ID:             "test123",
		Username:       "testuser",
		Host:           nil,
		Name:           nil,
		Description:    nil,
		IsBot:          false,
		IsLocked:       false,
		IsVerified:     false,
		FollowersCount: 0,
		FollowingCount: 0,
		NotesCount:     0,
		Fields:         []MisskeyField{},
	}

	account := convertToAccount(user)

	if account.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", account.Username)
	}

	if account.DisplayName != "testuser" {
		t.Errorf("Expected display name testuser, got %s", account.DisplayName)
	}

	if account.BIO != "" {
		t.Errorf("Expected empty bio, got %s", account.BIO)
	}

	if len(account.Profiles) != 0 {
		t.Errorf("Expected 0 profile fields, got %d", len(account.Profiles))
	}
}