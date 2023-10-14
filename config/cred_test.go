package config

import (
	"testing"

	"github.com/arrow2nd/nekomata/api"
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

func newTestCredentials() Credential {
	return Credential{
		Clients: map[string]*sharedapi.ClientCredential{
			api.ServiceMastodon: {
				Name:   "nekomata",
				ID:     "id_1",
				Secret: "secret_1",
			},
		},
		Users: map[string]*sharedapi.UserCredential{
			"test_1": {
				Server: "https://example.com",
				Token:  "user_token_1",
			},
			"test_2": {
				Server: "https://example.com",
				Token:  "user_token_2",
			},
		},
	}
}

func TestCredentialGet(t *testing.T) {
	c := newTestCredentials()

	t.Run("指定したユーザーの資格情報が取得できる", func(t *testing.T) {
		got, err := c.GetUser("test_1")

		assert.NoError(t, err)
		assert.Equal(t, c.Users["test_1"], got)
	})

	t.Run("見つからなかった際にエラーが返る", func(t *testing.T) {
		_, err := c.GetUser("hoge")
		assert.EqualError(t, err, "user not found: hoge")
	})
}

func TestGetAllUsernames(t *testing.T) {
	c := newTestCredentials()

	t.Run("取得できる", func(t *testing.T) {
		got := c.GetAllUsernames()

		assert.Contains(t, got, "test_1")
		assert.Contains(t, got, "test_2")
	})
}

func TestWrite(t *testing.T) {
	t.Run("追加できる", func(t *testing.T) {
		c := newTestCredentials()

		want := &sharedapi.UserCredential{
			Server: "test",
			Token:  "mochi",
		}

		c.AddUser("hiori", want)

		got, _ := c.GetUser("hiori")
		assert.Equal(t, want, got)
	})

	t.Run("同じIDを持つユーザを上書きできる", func(t *testing.T) {
		c := newTestCredentials()

		want := &sharedapi.UserCredential{
			Server: "test",
			Token:  "mochi",
		}

		c.AddUser("test_2", want)

		got, _ := c.GetUser("test_2")
		assert.Equal(t, want, got)
	})
}

func TestDelete(t *testing.T) {
	t.Run("削除できる", func(t *testing.T) {
		c := newTestCredentials()

		err := c.DeleteUser("test_1")
		assert.NoError(t, err)

		_, err = c.GetUser("test_1")
		assert.EqualError(t, err, "user not found: test_1")
	})

	t.Run("見つからない場合にエラーが返る", func(t *testing.T) {
		c := newTestCredentials()

		err := c.DeleteUser("hoge")
		assert.EqualError(t, err, "user not found: hoge")
	})
}
