package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/cli"
	"github.com/spf13/pflag"
)

func (a *App) newTimelineCmd() *cli.Command {
	timelineCmd := &cli.Command{
		Name:   "timeline",
		Short:  "Add timeline page",
		Hidden: global.isCLI,
	}

	timelineCmd.AddCommand(a.newTimelineSubCmds()...)

	return timelineCmd
}

func (a *App) newTimelineSubCmds() []*cli.Command {
	action := func(t timelineKind, c *cli.Command, f *pflag.FlagSet) error {
		page, err := newTimelinePage(t)
		if err != nil {
			return err
		}

		openBackground, _ := f.GetBool("background")
		if err := a.view.AddPage(page, !openBackground); err != nil {
			return err
		}

		go func() {
			if err := page.Load(); err != nil {
				label := fmt.Sprintf("load (%s)", t)
				global.SetErrorStatus(label, err.Error())
			}
		}()

		return nil
	}

	home := &cli.Command{
		Name:    "home",
		Short:   "Add home timeline page",
		SetFlag: setOpenBackgroundFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return action(homeTimeline, c, f)
		},
	}

	global := &cli.Command{
		Name:    "global",
		Short:   "Add global timeline page",
		SetFlag: setOpenBackgroundFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return action(globalTimeline, c, f)
		},
	}

	local := &cli.Command{
		Name:    "local",
		Short:   "Add local timeline page",
		SetFlag: setOpenBackgroundFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return action(localTimeline, c, f)
		},
	}

	return []*cli.Command{
		home,
		global,
		local,
	}
}
