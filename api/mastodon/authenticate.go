package mastodon

import (
	"fmt"
	"io"
	"net/url"

	"github.com/arrow2nd/nekomata/api/shared"
)

func (m *Mastodon) Authenticate(w io.Writer) (*shared.User, error) {
	return nil, nil
}

func (m *Mastodon) createAuthorizeURL(redirectURL string) string {
	q := url.Values{}

	q.Add("response_type", "code")
	q.Add("client_id", m.opts.ID)
	q.Add("redirect_uri", redirectURL)
	q.Add("scope", "read write follow")

	return fmt.Sprintf("%s/oauth/authorize?%s", m.opts.Server, q.Encode())
}
