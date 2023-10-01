package config_test

import (
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/arrow2nd/nekomata/config"
	"github.com/stretchr/testify/assert"
)

func newTestCredentials() config.Credentials {
	return config.Credentials{
		{
			Username: "test_1",
			ClientOpts: &shared.ClientOpts{
				Server:    "https://example.com",
				Name:      "nekomata",
				ID:        "id_1",
				Secret:    "secret_1",
				UserToken: "user_token_1",
			},
		},
		{
			Username: "test_2",
			ClientOpts: &shared.ClientOpts{
				Server:    "https://example.com",
				Name:      "nekomata",
				ID:        "id_2",
				Secret:    "secret_2",
				UserToken: "user_token_2",
			},
		},
	}
}

func TestCredentialGet(t *testing.T) {
	c := newTestCredentials()

	t.Run("指定したユーザーの認証情報が取得できる", func(t *testing.T) {
		got, err := c.Get("test_1")

		assert.NoError(t, err)
		assert.Equal(t, c[0], *got)
	})

	t.Run("見つからなかった際にエラーが返る", func(t *testing.T) {
		_, err := c.Get("hoge")
		assert.EqualError(t, err, "user not found: hoge")
	})
}

func TestGetAllUsernames(t *testing.T) {
	c := newTestCredentials()

	t.Run("取得できる", func(t *testing.T) {
		got := c.GetAllUsernames()
		want := []string{
			"test_1",
			"test_2",
		}

		assert.Equal(t, want, got)
	})
}

func TestWrite(t *testing.T) {
	t.Run("追加できる", func(t *testing.T) {
		c := newTestCredentials()

		want := &config.Credential{
			Username: "hiori",
			ClientOpts: &shared.ClientOpts{
				Server:    "test",
				Name:      "hoge",
				ID:        "fuga",
				Secret:    "piyo",
				UserToken: "mochi",
			},
		}

		c.Write(want)

		got, _ := c.Get("hiori")
		assert.Equal(t, want, got)
	})

	t.Run("同じIDを持つユーザを上書きできる", func(t *testing.T) {
		c := newTestCredentials()

		want := &config.Credential{
			Username: "test_2",
			ClientOpts: &shared.ClientOpts{
				Server:    "test",
				Name:      "hoge",
				ID:        "fuga",
				Secret:    "piyo",
				UserToken: "mochi",
			},
		}

		c.Write(want)

		got, _ := c.Get("test_2")
		assert.Equal(t, want, got)
	})
}

func TestDelete(t *testing.T) {
	t.Run("削除できる", func(t *testing.T) {
		c := newTestCredentials()

		err := c.Delete("test_1")
		assert.NoError(t, err)

		_, err = c.Get("test_1")
		assert.EqualError(t, err, "user not found: test_1")
	})

	t.Run("見つからない場合にエラーが返る", func(t *testing.T) {
		c := newTestCredentials()

		err := c.Delete("hoge")
		assert.EqualError(t, err, "user not found: hoge")
	})
}
