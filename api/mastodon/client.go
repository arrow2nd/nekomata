package mastodon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

var (
	defaultName   = ""
	defaultID     = ""
	defaultSecret = ""
)

type Mastodon struct {
	client *sharedapi.ClientOpts
	user   *sharedapi.UserOpts
}

func New(c *sharedapi.ClientOpts, u *sharedapi.UserOpts) *Mastodon {
	mastodon := &Mastodon{
		client: c,
		user:   u,
	}

	if c == nil || c.Name == "" || c.ID == "" || c.Secret == "" {
		mastodon.client = &sharedapi.ClientOpts{
			Name:   defaultName,
			ID:     defaultID,
			Secret: defaultSecret,
		}
	}

	return mastodon
}

type requestOpts struct {
	method      string
	contentType string
	url         string
	q           url.Values
	body        io.Reader
	isAuth      bool
}

func (m *Mastodon) request(opts *requestOpts, out interface{}) error {
	req, err := http.NewRequest(opts.method, opts.url, opts.body)
	if err != nil {
		return fmt.Errorf("create request error (%s): %w", opts.url, err)
	}

	if opts.isAuth {
		req.Header.Set("Authorization", "Bearer "+m.user.Token)
	}

	if opts.contentType != "" {
		req.Header.Set("Content-Type", opts.contentType)
	}

	if opts.q != nil {
		req.URL.RawQuery = opts.q.Encode()
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return &sharedapi.RequestError{
			URL: opts.url,
			Err: err,
		}
	}

	defer res.Body.Close()
	decorder := json.NewDecoder(res.Body)

	// エラーレスポンスをデコードして返す
	if res.StatusCode != http.StatusOK {
		e := &errorResponse{}
		if err := decorder.Decode(e); err != nil {
			return sharedapi.NewHTTPError(res)
		}
		return e
	}

	if out == nil {
		return nil
	}

	if err := decorder.Decode(out); err != nil {
		return &sharedapi.DecodeError{
			URL: opts.url,
			Err: err,
		}
	}

	return nil
}
