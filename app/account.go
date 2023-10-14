package app

import "github.com/arrow2nd/nekomata/api"

func login(username string) error {
	// ユーザー名が空ならアプリケーション認証から
	if username == "" {
	}

	// ログインユーザーの資格情報を取得
	userCred, err := global.conf.Creds.GetUser(username)
	if err != nil {
		return err
	}

	clientCred, err := global.conf.Creds.GetClient(userCred.Service)
	if err != nil {
		return err
	}

	client, err := api.NewClient(clientCred, userCred)
	if err != nil {
		return err
	}

	global.client = &client
	global.loginUsername = username

	return nil
}
