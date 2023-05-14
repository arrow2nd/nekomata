package shared_test

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestRecieveAuthenticateCode(t *testing.T) {
	lazyPost := func(t *testing.T, code string) {
		time.Sleep(1 * time.Second)

		req, _ := http.NewRequest(http.MethodPost, shared.AuthCallbackURL, nil)
		req.URL.RawQuery += "code=" + code

		_, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
	}

	t.Run("正しく認証コードを受け取れるか", func(t *testing.T) {
		go lazyPost(t, "hoge")

		code, err := shared.RecieveAuthenticateCode("code", func(s string) bool {
			return s != ""
		})

		assert.NoError(t, err)
		assert.Equal(t, "hoge", code, "受け取った認証コードが一致するか")
	})

	t.Run("受け取った認証コードが空文字の場合エラーになるか", func(t *testing.T) {
		go lazyPost(t, "")

		_, err := shared.RecieveAuthenticateCode("code", func(s string) bool {
			return s != ""
		})

		assert.ErrorContains(t, err, "failed to recieve code")
	})
}

func TestPrintAuthenticateURL(t *testing.T) {
	buf := &bytes.Buffer{}
	shared.PrintAuthenticateURL(buf, "hoge")
	assert.Contains(t, buf.String(), "hoge", "渡した文字列が含まれているか")
}
