package app

import (
	"errors"
	"os"
	"regexp"

	"github.com/arrow2nd/nekomata/api"
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/manifoldco/promptui"
)

func login(username string) error {
	// ユーザー名が空ならアプリケーション認証から
	if username == "" {
		account, err := authenticateAndSaveCredential()
		if err != nil {
			return err
		}

		// メインユーザーに設定
		global.conf.Pref.Feature.MainUser = account.Username
		if err := global.conf.SavePreferences(); err != nil {
			return err
		}

		username = account.Username
	}

	// ログインユーザーの資格情報を取得
	userCred, err := global.conf.Creds.GetUser(username)
	if err != nil {
		return err
	}

	clientCred, err := global.conf.Creds.GetClient(userCred.Server)
	if err != nil {
		return err
	}

	client, err := api.NewClient(clientCred, userCred)
	if err != nil {
		return err
	}

	global.client = client
	global.currentUsername = username

	return nil
}

func authenticateAndSaveCredential() (*sharedapi.Account, error) {
	// ログインするサービスを選択
	servicePrompt := promptui.Select{
		Label: "Select the service you wish to login",
		Items: api.GetAllServices(),
	}

	_, service, err := servicePrompt.Run()
	if err != nil {
		return nil, err
	}

	// サービスのドメインを入力
	domainPrompt := promptui.Prompt{
		Label:     "Enter the domain of the service",
		Default:   "https://",
		AllowEdit: true,
		Validate: func(d string) error {
			ok, err := regexp.MatchString(`^https?://[a-zA-Z0-9-_.]+\.[a-z]+(/.*)?$`, d)

			if !ok || err != nil {
				return errors.New("url format is invalid")
			}

			return nil
		},
	}

	server, err := domainPrompt.Run()
	if err != nil {
		return nil, err
	}

	// TODO: 入力されたドメインが選択したサービスのものか確認してもよさそう

	clientCred, _ := global.conf.Creds.GetClient(server)
	if clientCred == nil {
		clientCred = &sharedapi.ClientCredential{
			Name:    global.name,
			Service: service,
		}
	}

	userCred := &sharedapi.UserCredential{
		Server: server,
	}

	// クライアントを作成
	client, err := api.NewClient(clientCred, userCred)
	if err != nil {
		return nil, err
	}

	// 資格情報が空ならアプリケーションをサーバーに登録
	if clientCred.IsUncertified() {
		id, secret, err := client.RegisterNewApplication()
		if err != nil {
			return nil, err
		}

		clientCred.ID = id
		clientCred.Secret = secret

		global.conf.Creds.AddClient(server, clientCred)
	}

	// アプリケーション認証
	userToken, err := client.Authenticate(os.Stdout)
	if err != nil {
		return nil, err
	}

	// ログインユーザーを取得
	userCred.Token = userToken
	account, err := client.GetLoginAccount()
	if err != nil {
		return nil, err
	}

	global.conf.Creds.AddUser(account.Username, userCred)

	if err := global.conf.SaveCred(); err != nil {
		return nil, err
	}

	return account, nil
}
