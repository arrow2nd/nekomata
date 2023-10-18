package app

import (
	"github.com/arrow2nd/nekomata/cli"
	"github.com/spf13/pflag"
)

func (a *App) newTimelineCmd() *cli.Command {
	timelineCmd := &cli.Command{
		Name:   "timeline",
		Short:  "Add timeline page",
		Hidden: global.isCLI,
	}

	timelineCmd.AddCommand(
		a.newHomeTimelineCmd(),
	)

	return timelineCmd
}

func (a *App) newHomeTimelineCmd() *cli.Command {
	return &cli.Command{
		Name:  "home",
		Short: "Add home timeline page",
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			page, err := newTimelinePage(homeTimeline)
			if err != nil {
				return err
			}

			unfocus, _ := f.GetBool("unfocus")
			if err := a.view.AddPage(page, !unfocus); err != nil {
				return err
			}

			page.posts.view.SetChangedFunc(func() {
				a.app.Draw()
			})

			go func() {
				if err := page.Load(); err != nil {
					global.SetErrorStatus("timeline", err.Error())
				}
			}()

			return nil
		},
	}
}
