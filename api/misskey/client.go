package misskey

import "github.com/arrow2nd/nekomata/api/shared"

// Misskey : みすきー
type Misskey struct {
	config  *shared.Config
	baseURL string
}

// New : 新しいクライアントを生成
func New(c *shared.Config) *Misskey {
	baseURL := shared.CreateURL("https", c.Host, "api")

	return &Misskey{
		config:  c,
		baseURL: baseURL.String(),
	}
}
