package config

import (
	"github.com/arrow2nd/nekomata/api"
	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// Config : 設定
type Config struct {
	// Cred : 資格情報
	Creds *Credential
	// Pref : 環境設定
	Pref *Preferences
	// Style : スタイル定義
	Style *Style
	// DirPath : 設定ディレクトリのパス
	DirPath string
}

func New() (*Config, error) {
	path, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	defaultClients := map[string]*sharedapi.ClientCredential{}
	for _, s := range api.GetAllServices() {
		defaultClients[s] = &sharedapi.ClientCredential{}
	}

	return &Config{
		Creds: &Credential{
			Clients: defaultClients,
			Users:   map[string]*sharedapi.UserCredential{},
		},
		Pref:    defaultPreferences(),
		Style:   defaultStyle(),
		DirPath: path,
	}, nil
}
