package api

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/mastodon"
	"github.com/arrow2nd/nekomata/api/misskey"
	"github.com/arrow2nd/nekomata/api/sharedapi"
)

const (
	ServiceMastodon = "Mastodon"
	ServiceMisskey  = "Misskey"
)

// GetAllServices : 利用可能なサービスの一覧
func GetAllServices() []string {
	return []string{
		ServiceMastodon,
		// Misskey,
	}
}

// NewClient : サービスを指定してクライアントを作成
func NewClient(c *sharedapi.ClientCredential, u *sharedapi.UserCredential) (sharedapi.Client, error) {
	switch u.Service {
	case ServiceMastodon:
		return mastodon.New(c, u)
	case ServiceMisskey:
		return misskey.New(c, u)
	}

	return nil, fmt.Errorf("unsupported services: %s", u.Service)
}
