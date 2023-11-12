package mastodon

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// media : メディア
type media struct {
	// ID : メディアID
	ID string `json:"id"`
	// Type : 種類 (unknown, image, gif, video, audio)
	Type string `json:"type"`
	// URL : オリジナルのメディアを指すURL
	URL string `json:"url"`
	// PreviewURL : スケールダウンされたメディアを指すURL
	PreviewURL string `json:"preview_url"`
	// Meta : メタ情報
	Meta mediaMeta `json:"meta"`
}

// mediaMeta : メディアのメタ情報
type mediaMeta struct {
	// Original : オリジナルサイズ
	Original mediaSize `json:"original"`
	// Small : 縮小サイズ
	Small mediaSize `json:"small"`
}

// mediaSize : メディアサイズ
type mediaSize struct {
	// Width : 幅
	Width int `json:"width"`
	// Height : 高さ
	Height int `json:"height"`
}

func (m *Mastodon) UploadMedia(filename string, src io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, src); err != nil {
		return "", err
	}

	writer.Close()

	opts := &requestOpts{
		method:      http.MethodPost,
		contentType: writer.FormDataContentType(),
		url:         endpointMediaUploadAsync.URL(m.user.Server, nil),
		q:           nil,
		body:        body,
		isAuth:      true,
	}

	res := media{}
	if err := m.request(opts, &res); err != nil {
		return "", err
	}

	return res.ID, nil
}
