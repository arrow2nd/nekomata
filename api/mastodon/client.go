package mastodon

import (
	"encoding/json"
	"fmt"
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

func (m *Mastodon) request(method, url string, q url.Values, auth bool, out interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("create request error (%s): %w", url, err)
	}

	if auth {
		req.Header.Set("Authorization", "Bearer "+m.opts.UserToken)
	}

	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return &shared.RequestError{
			URL: url,
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

	if err := decorder.Decode(out); err != nil {
		return &shared.DecodeError{
			URL: url,
			Err: err,
		}
	}

	return nil
}
