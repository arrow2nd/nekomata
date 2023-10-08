package api

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/mastodon"
	"github.com/arrow2nd/nekomata/api/misskey"
	"github.com/arrow2nd/nekomata/api/sharedapi"
)

type Service string

const (
	Mastodon Service = "Mastodon"
	Misskey  Service = "Misskey"
)

// GetAllServices : 利用可能なサービスの一覧
func GetAllServices() []Service {
	return []Service{
		Mastodon,
		// Misskey,
	}
}

// NewClient : サービスを指定してクライアントを作成
func NewClient(service Service, opts *sharedapi.ClientOpts) (sharedapi.Client, error) {
	switch service {
	case Mastodon:
		return mastodon.New(opts), nil
	case Misskey:
		return misskey.New(opts), nil
	}

	return nil, fmt.Errorf("unsupported services: %s", service)
}
