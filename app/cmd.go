package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/cli"
	"github.com/spf13/pflag"
)

func newCmd() *cli.Command {
	return &cli.Command{
		Name:  "nekomata",
		Short: "🐱 Multi-SNS client with TUI",
		Long:  "TUI social networking client for multiple services",
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("user", "u", shared.conf.Pref.Feature.MainUser, "username of the account to log in to")
			f.BoolP("version", "v", false, "Display version")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// コマンドラインからの実行ならバージョンを表示して終了
			if ver, _ := f.GetBool("version"); shared.isCLI && ver {
				exit(fmt.Sprintf("🐱 %s for v.%s", c.Name, version))
			}

			// エラーメッセージを組み立てる
			arg := f.Arg(0)
			if arg != "" {
				arg = ": " + arg
			}

			return fmt.Errorf("unavailable or not found command%s", arg)
		},
	}
}

func setUnfocusFlag(f *pflag.FlagSet) {
	f.BoolP("unfocus", "u", false, "no focus on page")
}
