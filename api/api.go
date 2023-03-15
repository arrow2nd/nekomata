package api

import (
	"fmt"
	"io"

	"github.com/arrow2nd/nekomata/api/mastodon"
	"github.com/arrow2nd/nekomata/api/misskey"
	"github.com/arrow2nd/nekomata/api/shared"
)

type Service string

const (
	ServiceMastodon Service = "Mastodon"
	ServiceMisskey  Service = "Misskey"
)

func NewClient(w io.Writer, s Service, c *shared.ClientOpts) (shared.Client, error) {
	switch s {
	case ServiceMisskey:
		return misskey.New(c), nil
	case ServiceMastodon:
		return mastodon.New(c), nil
	}

	return nil, fmt.Errorf("Incorrect service: %s", s)
}
