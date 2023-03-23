package mastodon

import (
	"encoding/json"
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

func (m *Mastodon) post(endpoint shared.Endpoint, q url.Values, out interface{}) error {
	url := endpoint.URL(m.opts.Server)
	res, err := http.PostForm(url, q)
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
