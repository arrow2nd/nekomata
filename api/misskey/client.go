package misskey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arrow2nd/nekomata/api"
)

type Misskey struct {
	opts *api.ClientOpts
}

// New : 新しいクライアントを生成
func New(c *api.ClientOpts) *Misskey {
	return &Misskey{
		opts: c,
	}
}

func (m *Misskey) post(endpoint api.Endpoint, in, out interface{}) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("create payload error (%s): %w", endpoint, err)
	}

	endpointURL := endpoint.URL(m.opts.Server, nil)
	req, err := http.NewRequest(http.MethodPost, endpointURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create request error (%s): %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &api.RequestError{
			URL: endpointURL,
			Err: err,
		}
	}

	defer res.Body.Close()

	// TODO: 200以外も返ってきてた気がするので修正する
	if res.StatusCode != http.StatusOK {
		return api.NewHTTPError(res)
	}

	decorder := json.NewDecoder(res.Body)
	if err := decorder.Decode(out); err != nil {
		return &api.DecodeError{
			URL: endpointURL,
			Err: err,
		}
	}

	return nil
}
