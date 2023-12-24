package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

func configCmd() *cli.Command {
	return &cli.Command{
		Name:      "config",
		Usage:     "Edit configuration file",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "editor",
				Aliases: []string{"e"},
				Usage:   "specify which editor to use",
				Value:   os.Getenv("EDITOR"),
			},
		},
	}
}
