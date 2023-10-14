package config

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// Credential : 資格情報
type Credential struct {
	// Clients : クライアントのカスタム資格情報
	Clients map[string]*sharedapi.ClientCredential `toml:"clients"`
	// Users : ユーザー
	Users map[string]*sharedapi.UserCredential `tonl:"users"`
}

// GetClient : クライアントの資格情報を取得
func (c *Credential) GetClient(service string) (*sharedapi.ClientCredential, error) {
	cred, ok := c.Clients[service]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", service)
	}

	return cred, nil
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
func (c *Credential) AddUser(username string, user *sharedapi.UserCredential) {
	c.Users[username] = user
}

// DeleteUser : ユーザーの資格情報を削除
func (c *Credential) DeleteUser(username string) error {
	if _, ok := c.Users[username]; !ok {
		return fmt.Errorf("user not found: %s", username)
	}

	delete(c.Users, username)

	return nil
}
