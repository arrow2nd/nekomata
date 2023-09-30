package mastodon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arrow2nd/nekomata/api/shared"
)

type Mastodon struct {
	opts *shared.ClientOpts
}

func New(c *shared.ClientOpts) *Mastodon {
	return &Mastodon{
		opts: c,
	}
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
		req.Header.Set("Authorization", "Bearer "+m.opts.UserToken)
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
		return &shared.RequestError{
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
			return shared.NewHTTPError(res)
		}
		return e
	}

	if out == nil {
		return nil
	}

	if err := decorder.Decode(out); err != nil {
		return &shared.DecodeError{
			URL: opts.url,
			Err: err,
		}
	}

	return nil
}
