package config

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// Credentials : 認証情報
type Credentials struct {
	// Clients : クライアントのカスタム認証情報
	Clients map[string]*sharedapi.ClientOpts
	// Users : ユーザー
	Users map[string]*sharedapi.UserOpts
}

// Get : 取得
func (c Credentials) Get(username string) (*sharedapi.UserOpts, error) {
	for u, cred := range c.Users {
		if u == username {
			return cred, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", username)
}

// GetAllNames : 全てのユーザ名を取得
func (c Credentials) GetAllUsernames() []string {
	ls := []string{}

	for username := range c.Users {
		ls = append(ls, username)
	}

	return ls
}

// Add : 追加
func (c *Credentials) Add(username string, user *sharedapi.UserOpts) {
	c.Users[username] = user
}

// Delete : 削除
func (c *Credentials) Delete(username string) error {
	if _, ok := c.Users[username]; !ok {
		return fmt.Errorf("user not found: %s", username)
	}

	delete(c.Users, username)

	return nil
}
