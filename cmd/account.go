package cmd

import "github.com/urfave/cli/v2"

func accountCmd() *cli.Command {
	return &cli.Command{
		Name:    "account",
		Aliases: []string{},
		Usage:   "Manage your account",
		Subcommands: []*cli.Command{
			accountAddCmd(),
			accountDeleteCmd(),
			accountSetCmd(),
			accountListCmd(),
		},
	}
}

func accountAddCmd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "main",
				Aliases: []string{"m"},
				Usage:   "set as main user",
			},
		},
	}
}

func accountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete account",
		ArgsUsage: "[user name]",
	}
}

func accountListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "Show accounts that have been added",
	}
}

func accountSetCmd() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set main account",
		ArgsUsage: "[user name]",
	}
}
