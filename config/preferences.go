package config

const PreferencesVersion = 1

type Feature struct {
	// MainAccount : メインのアカウント
	MainAccount string `toml:"main_account"`
}

type Preferences struct {
	Version int     `toml:"version"`
	Feature Feature `toml:"feature"`
}

func defaultPreferences() *Preferences {
	return &Preferences{
		Version: PreferencesVersion,
		Feature: Feature{
			MainAccount: "",
		},
	}
}
