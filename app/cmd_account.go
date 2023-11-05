package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/cli"
	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

// getTargetUsername : 対象のユーザー名を引数 or 入力から取得
func (a *App) getTargetUsername(label string, f *pflag.FlagSet) (string, error) {
	username := f.Arg(0)

	// 指定が無い場合選択
	if username == "" {
		prompt := promptui.Select{
			Label: label,
			Items: global.conf.Creds.GetAllUsernames(),
		}

		_, selected, err := prompt.Run()
		if err != nil {
			return "", err
		}

		return selected, nil
	}

	return username, nil
}

func (a *App) newAccountCmd() *cli.Command {
	cmd := &cli.Command{
		Name:     "account",
		Short:    "Manage your account",
		Validate: cli.NoArgs(),
	}

	cmd.AddCommand(
		a.newAccountAddCmd(),
		a.newAccountDeleteCmd(),
		a.newAccountListCmd(),
		a.newAccountSetCmd(),
		a.newAccountSwitchCmd(),
	)

	return cmd
}

func (a *App) newAccountAddCmd() *cli.Command {
	return &cli.Command{
		Name:     "add",
		Short:    "Add account",
		Hidden:   !global.isCLI,
		Validate: cli.NoArgs(),
		SetFlag: func(f *pflag.FlagSet) {
			f.BoolP("main", "m", false, "set as main user")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			account, err := authenticateAndSaveCredential()
			if err != nil {
				return err
			}

			// メインユーザーに設定
			if main, _ := f.GetBool("main"); main {
				if err := global.conf.SavePreferences(); err != nil {
					return err
				}
			}

			fmt.Printf("🐱 Logged in: %s (%s)\n", account.DisplayName, account.Username)
			return nil
		},
	}
}

func (a *App) newAccountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Short: "Delete account",
		Long: `Delete account.
If you do not specify an account name, you can select it interactively.`,
		UsageArgs: "[user name]",
		Example:   "delete arrow2nd",
		Hidden:    !global.isCLI,
		Validate:  cli.RangeArgs(0, 1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			username, err := a.getTargetUsername("Select the account to delete", f)
			if err != nil {
				return err
			}

			if err := global.conf.Creds.DeleteUser(username); err != nil {
				return err
			}

			if err := global.conf.SaveCred(); err != nil {
				return err
			}

			fmt.Printf("🐱 Deleted: %s\n", f.Arg(0))
			return nil
		},
	}
}

func (a *App) newAccountListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Short:    "Show accounts that have been added",
		Hidden:   !global.isCLI,
		Validate: cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			for _, u := range global.conf.Creds.GetAllUsernames() {
				current := " "

				if u == global.currentUsername {
					current = "*"
				}

				fmt.Printf(" %s %s\n", current, u)
			}

			return nil
		},
	}
}

func (a *App) newAccountSetCmd() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Short:     "Set main account",
		UsageArgs: "[user name]",
		Hidden:    !global.isCLI,
		Validate:  cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			username, err := a.getTargetUsername("Select the account to be set as main", f)
			if err != nil {
				return err
			}

			global.conf.Pref.Feature.MainAccount = username
			if err := global.conf.SavePreferences(); err != nil {
				return err
			}

			fmt.Printf("🐱 Set up for main account: %s\n", username)
			return nil
		},
	}
}

func (a *App) newAccountSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
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
