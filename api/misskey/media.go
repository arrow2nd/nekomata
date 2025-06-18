package misskey

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"encoding/json"
)

// DriveFile : Misskeyのドライブファイル
type DriveFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	MD5       string `json:"md5"`
	Size      int    `json:"size"`
	URL       string `json:"url"`
	CreatedAt string `json:"createdAt"`
}

func (m *Misskey) UploadMedia(filename string, src io.Reader) (string, error) {
	// マルチパートフォームデータを作成
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// ファイルフィールドを追加
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file error: %w", err)
	}

	// ファイル内容をコピー
	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return "", fmt.Errorf("copy file error: %w", err)
	}

	// アクセストークンを追加
	err = writer.WriteField("i", m.user.Token)
	if err != nil {
		return "", fmt.Errorf("write token field error: %w", err)
	}

	// マルチパートフォームを閉じる
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("close writer error: %w", err)
	}

	// HTTPリクエストを作成
	endpointURL := endpointDriveFilesCreate.URL(m.user.Server, nil)
	req, err := http.NewRequest(http.MethodPost, endpointURL, &buf)
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// リクエストを送信
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status: %d", res.StatusCode)
	}

	// レスポンスをパース
	var file DriveFile
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&file); err != nil {
		return "", fmt.Errorf("decode response error: %w", err)
	}

	return file.ID, nil
}
