package mastodon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

func (m *Mastodon) request(method string, endpoint shared.Endpoint, q url.Values, auth bool, out interface{}) error {
	url := endpoint.URL(m.opts.Server)
	req, err := http.NewRequest(method, url, strings.NewReader(q.Encode()))
	if err != nil {
		return fmt.Errorf("create request error (%s): %w", endpoint, err)
	}

	if auth {
		req.Header.Set("Authorization", "Bearer "+m.opts.UserToken)
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return &shared.RequestError{
			Endpoint: endpoint,
			Err:      err,
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
			Endpoint: endpoint,
			Err:      err,
		}
	}

	return nil
}
