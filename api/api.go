package api

import (
	"io"

	"github.com/arrow2nd/nekomata/api/misskey"
	"github.com/arrow2nd/nekomata/api/shared"
)

type Service string

const (
	ServiceMastodon Service = "Mastodon"
	ServiceMisskey  Service = "Misskey"
)

func NewClient(w io.Writer, s Service, c *shared.Config) (shared.Client, error) {
	// TODO: サービス毎にそれぞれクライアントを返す
	return misskey.Init(c), nil
}
