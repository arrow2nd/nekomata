package mastodon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

type Mastodon struct {
	client *sharedapi.ClientCredential
	user   *sharedapi.UserCredential
}

func New(c *sharedapi.ClientCredential, u *sharedapi.UserCredential) (*Mastodon, error) {
	if c == nil {
		return nil, errors.New("mastodon: client credential is empty")
	}

	return &Mastodon{
		client: c,
		user:   u,
	}, nil
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
