package shared

import "io"

type Config struct {
	AppName      string
	Host         string
	ClientID     string
	ClientSecret string
	Token        string
}

// Client : API クライアント
type Client interface {
	Authenticate(io.Writer) (*AuthResponse, error)
}
