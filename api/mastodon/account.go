package mastodon

import "time"

type account struct {
	ID             string          `json:"id"`
	Username       string          `json:"username"`
	Acct           string          `json:"acct"`
	DisplayName    string          `json:"display_name"`
	Locked         bool            `json:"locked"`
	Bot            bool            `json:"bot"`
	CreatedAt      time.Time       `json:"created_at"`
	Note           string          `json:"note"`
	FollowersCount int             `json:"followers_count"`
	FollowingCount int             `json:"following_count"`
	StatusesCount  int             `json:"statuses_count"`
	Fields         []accountFields `json:"fields"`
}

type accountFields struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
