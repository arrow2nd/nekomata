package shared

import "io"

// Client : API クライアント
type Client interface {
	Authenticate(io.Writer) (*User, error)
	GetAnnouncements() ([]*Announcement, error)
}

// ClientOpts : クライアントの設定
type ClientOpts struct {
	Server    string
	Name      string
	ID        string
	Secret    string
	UserToken string
}
