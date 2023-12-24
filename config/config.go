package config

import "github.com/arrow2nd/nekomata/api/sharedapi"

type Config struct {
	Credential  *Credential
	Preferences *Preferences
	DirPath     string
}

func New() (*Config, error) {
	path, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		Credential: &Credential{
			Clients: map[string]*sharedapi.ClientCredential{},
			Users:   map[string]*sharedapi.UserCredential{},
		},
		Preferences: defaultPreferences(),
		DirPath:     path,
	}, nil
}
