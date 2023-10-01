package config

// Config : 設定
type Config struct {
	// Cred : 認証情報
	Creds *Credentials
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

	return &Config{
		Creds:   &Credentials{},
		Pref:    defaultPreferences(),
		Style:   defaultStyle(),
		DirPath: path,
	}, nil
}
