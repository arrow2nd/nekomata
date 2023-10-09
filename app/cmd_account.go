package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arrow2nd/nekomata/api"
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/cli"
	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

func (a *App) newAccountCmd() *cli.Command {
	cmd := &cli.Command{
		Name:      "account",
		Shorthand: "a",
		Short:     "Manage your account",
		Validate:  cli.NoArgs(),
	}

	cmd.AddCommand(
		a.newAccountAddCmd(),
		a.newAccountDeleteCmd(),
		a.newAccountListCmd(),
		a.newAccountSwitchCmd(),
	)

	return cmd
}

func (a *App) newAccountAddCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Shorthand: "a",
		Short:     "Add account",
		Hidden:    !global.isCLI,
		Validate:  cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// ログインするサービスを選択
			servicePrompt := promptui.Select{
				Label: "Service",
				Items: api.GetAllServices(),
			}

			_, service, err := servicePrompt.Run()
			if err != nil {
				return err
			}

			// サービスのドメインを入力
			domainPrompt := promptui.Prompt{
				Label:     "Domain",
				Default:   "https://",
				AllowEdit: true,
				Validate: func(u string) error {
					if !strings.HasPrefix(u, "http") {
						return errors.New("must begin with http")
					}
					return nil
				},
			}

			server, err := domainPrompt.Run()
			if err != nil {
				return err
			}

			// クライアントを作成
			userOpts := &sharedapi.UserOpts{
				Server: server,
			}

			client, err := api.NewClient(service, global.conf.Creds.Clients[service], userOpts)
			if err != nil {
				return nil
			}

			// アプリケーション認証
			userToken, err := client.Authenticate(os.Stdout)
			if err != nil {
				return err
			}

			// ログインユーザーを取得
			userOpts.Token = userToken
			user, err := client.GetLoginAccount()
			if err != nil {
				return err
			}

			fmt.Printf("🐱 Logged in: %s (%s)\n", user.DisplayName, user.Username)

			// 保存
			global.conf.Creds.Add(user.Username, userOpts)
			return global.conf.SaveCred()
		},
	}
}

func (a *App) newAccountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Shorthand: "d",
		Short:     "Delete account",
		Long: `Delete account.
If you do not specify an account name, you can select it interactively.`,
		UsageArgs: "[user name]",
		Example:   "delete arrow2nd",
		Hidden:    !global.isCLI,
		Validate:  cli.RangeArgs(0, 1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			target := f.Arg(0)

			// 指定が無い場合、ユーザを選択
			if target == "" {
				prompt := promptui.Select{
					Label: "Account to delete",
					Items: global.conf.Creds.GetAllUsernames(),
				}

				_, seletecd, err := prompt.Run()
				if err != nil {
					return err
				}

				target = seletecd
			}

			if err := global.conf.Creds.Delete(target); err != nil {
				return err
			}

			if err := global.conf.SaveCred(); err != nil {
				return err
			}

			fmt.Printf("successfully deleted: %s\n", f.Arg(0))
			return nil
		},
	}
}

func (a *App) newAccountListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Shorthand: "l",
		Short:     "Show accounts that have been added",
		Hidden:    !global.isCLI,
		Validate:  cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			for _, u := range global.conf.Creds.GetAllUsernames() {
				current := " "

				// TODO: メインユーザーなら * を付ける
				// if u == shared.api.CurrentUser.UserName {
				// 	current = "*"
				// }

				fmt.Printf(" %s %s\n", current, u)
			}

			return nil
		},
	}
}

func (a *App) newAccountSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Shorthand: "s",
		Short:     "Switch the account to be used",
		UsageArgs: "[user name]",
		Example:   "switch arrow2nd",
		Hidden:    global.isCLI,
		Validate:  cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// TODO: 既にログイン中なら切り替えない
			// if f.Arg(0) == shared.api.CurrentUser.UserName {
			// 	return errors.New("account currently logged in")
			// }

			// TODO: ログインする
			// if err := loginAccount(f.Arg(0)); err != nil {
			// 	return err
			// }

			// 初期化
			a.view.Reset()
			a.statusBar.DrawAccountInfo()
			a.initAutocomplate()
			a.execStartupCommands()

			return nil
		},
	}
}
