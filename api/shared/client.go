package shared

import "io"

// Client : API クライアント
type Client interface {
	Authenticate(io.Writer) (*AuthResponse, error)
}

// Config : クライアントの設定
type Config struct {
	AppName      string
	Host         string
	ClientID     string
	ClientSecret string
	Token        string
}
