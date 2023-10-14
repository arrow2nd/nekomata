package app

import (
	"github.com/arrow2nd/nekomata/cli"
	"github.com/spf13/pflag"
)

func (a *App) newQuitCmd() *cli.Command {
	return &cli.Command{
		Name:     "quit",
		Short:    "Quit the application",
		Validate: cli.NoArgs(),
		Hidden:   global.isCLI,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			a.quitApp()
			return nil
		},
	}
}
