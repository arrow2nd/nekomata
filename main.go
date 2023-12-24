package main

import (
	"os"

	"github.com/arrow2nd/nekomata/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cmd.NewCmd()

	if err := app.Run(os.Args); err != nil {
		cli.Exit(err, 1)
	}
}
