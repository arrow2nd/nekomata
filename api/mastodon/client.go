package mastodon

import "github.com/arrow2nd/nekomata/api/shared"

type Mastodon struct {
	opts *shared.ClientOpts
}

func New(c *shared.ClientOpts) *Mastodon {
	return &Mastodon{
		opts: c,
	}
}
