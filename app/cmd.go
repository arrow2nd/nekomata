package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/app/exit"
	"github.com/arrow2nd/nekomata/cli"
	"github.com/spf13/pflag"
)

func newCmd() *cli.Command {
	return &cli.Command{
		Name:  global.name,
		Short: "🐱 Multi-SNS client with TUI",
		Long:  "TUI social networking client for multiple services",
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("account", "u", global.conf.Pref.Feature.MainAccount, "user account to login")
			f.BoolP("version", "v", false, "Display version")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// コマンドラインからの実行ならバージョンを表示して終了
			if ver, _ := f.GetBool("version"); global.isCLI && ver {
				exit.OK(fmt.Sprintf("🐱 %s for v.%s", c.Name, global.version))
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

func setOpenBackgroundFlag(f *pflag.FlagSet) {
	f.BoolP("background", "u", false, "open the page in the background")
}
