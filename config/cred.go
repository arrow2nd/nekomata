package config

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/shared"
)

type Credentials map[string]*shared.ClientOpts

// Get : 取得
func (c Credentials) Get(username string) (*shared.ClientOpts, error) {
	for u, cred := range c {
		if u == username {
			return cred, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", username)
}

// GetAllNames : 全てのユーザ名を取得
func (c Credentials) GetAllUsernames() []string {
	ls := []string{}

	for username := range c {
		ls = append(ls, username)
	}

	return ls
}

// Add : 追加
func (c *Credentials) Add(username string, opts *shared.ClientOpts) {
	(*c)[username] = opts
}

// Delete : 削除
func (c *Credentials) Delete(username string) error {
	if _, ok := (*c)[username]; !ok {
		return fmt.Errorf("user not found: %s", username)
	}

	delete(*c, username)

	return nil
}
