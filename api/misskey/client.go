package misskey

import (
	"net/url"

	"github.com/arrow2nd/nekomata/api/shared"
)

type Misskey struct {
	shared.Client
	config  *shared.Config
	baseURL string
}

func Init(c *shared.Config) *Misskey {
	baseURL := url.URL{}
	baseURL.Scheme = "https"
	baseURL = *baseURL.JoinPath(c.Host, "api")

	return &Misskey{
		config:  c,
		baseURL: baseURL.String(),
	}
}
