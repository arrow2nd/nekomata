package misskey

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiAuthRun(t *testing.T) {
	// 取得できるか
}

func TestCreateAuthURL(t *testing.T) {
	m := &miAuth{
		Name: "test_app",
		Host: "example.com",
	}

	p := []string{
		"aaaa",
		"bbbb",
	}

	u := m.createAuthURL(p)
	r := regexp.MustCompile("https://example.com/miauth/.+callback=http%3A%2F%2Flocalhost%3A3000%2Fcallback&name=test_app&permission=aaaa%2Cbbbb")

	assert.Regexp(t, r, u)
}

func TestHandleCallback(t *testing.T) {
	// セッションIDが一致しない

	// チャンネル経由でセッションIDが受け取れるか
}

func TestRecieveToken(t *testing.T) {
	// Postに失敗

	// HTTPエラー

	// JSONデコードエラー

	// トークン取得失敗

	// 取得できるか
}
