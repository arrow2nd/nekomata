package config

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

type Credential struct {
	// Clients : クライアント
	Clients map[string]*sharedapi.ClientCredential `toml:"clients"`
	// Users : ユーザー
	Users map[string]*sharedapi.UserCredential `toml:"users"`
}

// GetClient : クライアントの資格情報を取得
func (c *Credential) GetClient(server string) (*sharedapi.ClientCredential, error) {
	cred, ok := c.Clients[server]
	if !ok {
		return nil, fmt.Errorf("client credential not found: %s", server)
	}

	return cred, nil
}

// AddClient : クライアントの資格情報を追加
func (c *Credential) AddClient(server string, cred *sharedapi.ClientCredential) {
	c.Clients[server] = cred
}

// GetUser : ユーザーの資格情報を取得
func (c *Credential) GetUser(username string) (*sharedapi.UserCredential, error) {
	cred, ok := c.Users[username]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", username)
	}

	return cred, nil
}

// GetAllUsernames : 全てのユーザ名を取得
func (c *Credential) GetAllUsernames() []string {
	ls := []string{}

	for username := range c.Users {
		ls = append(ls, username)
	}

	return ls
}

// AddUser : ユーザーの資格情報を追加
func (c *Credential) AddUser(username string, cred *sharedapi.UserCredential) {
	c.Users[username] = cred
}

// DeleteUser : ユーザーの資格情報を削除
func (c *Credential) DeleteUser(username string) error {
	if _, ok := c.Users[username]; !ok {
		return fmt.Errorf("user not found: %s", username)
	}

	delete(c.Users, username)

	return nil
}
