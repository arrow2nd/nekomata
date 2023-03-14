package misskey

import (
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMiAuthRun(t *testing.T) {
	// 取得できるか
}

func TestMiAuthCreateURL(t *testing.T) {
	m := &miAuth{
		Name: "test_app",
		Host: "example.com",
	}

	p := []string{
		"aaaa",
		"bbbb",
	}

	u, _ := m.createURL(p)
	r := regexp.MustCompile("https://example.com/miauth/.+callback=http%3A%2F%2Flocalhost%3A3000%2Fcallback&name=test_app&permission=aaaa%2Cbbbb")

	assert.Regexp(t, r, u, "正しい形式か")
}

func TestMiAuthRecieveSessionID(t *testing.T) {
	type result struct {
		id  string
		err error
	}

	run := func(r chan *result, id string) {
		m := &miAuth{Name: "test_app", Host: "example.com"}
		recieveID, err := m.recieveSessionID(id)
		r <- &result{id: recieveID, err: err}
	}

	postCallback := func(id string) (*http.Response, error) {
		q := url.Values{}
		q.Set("session", id)

		url := shared.CreateURL("http", listenAddr, "callback").String()
		return http.Post(url+"?"+q.Encode(), "", nil)
	}

	t.Run("セッションIDが受け取れるか", func(t *testing.T) {
		id, _ := uuid.NewUUID()
		wantSessionID := id.String()

		result := make(chan *result, 1)
		go run(result, wantSessionID)

		res, err := postCallback(wantSessionID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.NoError(t, r.err)
		assert.Equal(t, wantSessionID, r.id, "受け取ったセッションIDが生成したものと一致するか")
	})

	t.Run("セッションIDが一致しない場合エラーを返すか", func(t *testing.T) {
		id, _ := uuid.NewUUID()
		wantSessionID := id.String()

		result := make(chan *result, 1)
		go run(result, wantSessionID)

		res, err := postCallback("hogehoge")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "呼び出し元に適切なステータスコードが返っているか")

		r := <-result
		assert.ErrorContains(t, r.err, "failed to recieve session id")
	})
}

func TestMiAuthRecieveToken(t *testing.T) {
	// Postに失敗

	// HTTPエラー

	// JSONデコードエラー

	// トークン取得失敗

	// 取得できるか
}
