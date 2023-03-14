package misskey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arrow2nd/nekomata/api/shared"
)

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

func (m *Misskey) post(endpoint string, in, out interface{}) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	u, _ := url.JoinPath(m.baseURL, endpoint)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	defer res.Body.Close()

	// TODO: 200以外も返ってきてた気がするので修正する
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http status error: %s", res.Status)
	}

	decorder := json.NewDecoder(res.Body)
	if err := decorder.Decode(out); err != nil {
		return fmt.Errorf("failed to decord json: %w", err)
	}

	return nil
}
