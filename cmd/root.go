package cmd

import "github.com/urfave/cli/v2"

func NewCmd() *cli.App {
	return &cli.App{
		Name:        "nekometa",
		Usage:       "🐱 Multi-SNS client with TUI",
		Description: "TUI social networking client for multiple services",
		Commands: []*cli.Command{
			accountCmd(),
			configCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "account",
				Aliases: []string{"a"},
				Usage:   "user account to login",
			},
		},
	}
}
