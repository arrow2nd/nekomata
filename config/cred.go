package config

import (
	"fmt"

	"github.com/arrow2nd/nekomata/api/shared"
)

type Credential struct {
	Username   string             `toml:"username"`
	ClientOpts *shared.ClientOpts `toml:"client_options"`
}

type Credentials []Credential

// Get : 取得
func (c Credentials) Get(username string) (*Credential, error) {
	for _, cred := range c {
		if cred.Username == username {
			return &cred, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", username)
}

// GetAllNames : 全てのユーザ名を取得
func (c Credentials) GetAllUsernames() []string {
	ls := []string{}

	for _, cred := range c {
		ls = append(ls, cred.Username)
	}

	return ls
}

// Write : 書込む
func (c *Credentials) Write(newCred *Credential) {
	// 同じユーザが居れば上書きする
	for i, cred := range *c {
		if cred.Username == newCred.Username {
			(*c)[i] = *newCred
			return
		}
	}

	// 新規追加
	*c = append(*c, *newCred)
}

// Delete : 削除
func (c *Credentials) Delete(username string) error {
	err := fmt.Errorf("user not found: %s", username)
	tmp := []Credential{}

	for _, cred := range *c {
		if cred.Username == username {
			err = nil
			continue
		}

		tmp = append(tmp, cred)
	}

	if err != nil {
		return err
	}

	*c = tmp
	return nil
}
