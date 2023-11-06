package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/cli"
	"github.com/arrow2nd/nekomata/config"
	"github.com/spf13/pflag"
)

func (a *App) newDocsCmd() *cli.Command {
	cmd := &cli.Command{
		Name:     "docs",
		Short:    "Show documentation",
		Hidden:   global.isCLI,
		Validate: cli.NoArgs(),
	}

	cmd.AddCommand(a.newDocsKeybindingsCmd())

	return cmd
}

func (a *App) newDocsKeybindingsCmd() *cli.Command {
	k := global.conf.Pref.Keybindings

	global := fmt.Sprintf(
		`[Global]
  %-20s Quit application

`,
		k.Global.GetString(config.ActionQuit),
	)

	view := fmt.Sprintf(
		`[View]
  %-20s Select previous tab
  %-20s Select next tab
  %-20s Close current page
  %-20s Redraw screen (window width changes are not reflected)
  %-20s Focus the command line
  %-20s Show documentation for keybindings

`,
		k.View.GetString(config.ActionSelectPrevTab),
		k.View.GetString(config.ActionSelectNextTab),
		k.View.GetString(config.ActionClosePage),
		k.View.GetString(config.ActionRedraw),
		k.View.GetString(config.ActionFocusCmdLine),
		k.View.GetString(config.ActionShowHelp),
	)

	page := fmt.Sprintf(
		`[Common Page]
  %-20s Reload page

`,
		k.Page.GetString(config.ActionReloadPage),
	)

	post := fmt.Sprintf(
		`[Post list]
  %-20s Scroll up
  %-20s Scroll down
  %-20s Move cursor up
  %-20s Move cursor down
  %-20s Move cursor top
  %-20s Move cursor bottom

  %-20s Reaction
  %-20s Unreaction
  %-20s Repost
  %-20s Unrepost
  %-20s Bookmark
  %-20s Unbookmark
  %-20s New post
  %-20s Reply
  %-20s Delete
  %-20s Open in browser
  %-20s Open user timeline page
  %-20s Copy link to clipboard
`,
		k.Posts.GetString(config.ActionScrollUp),
		k.Posts.GetString(config.ActionScrollDown),
		k.Posts.GetString(config.ActionCursorUp),
		k.Posts.GetString(config.ActionCursorDown),
		k.Posts.GetString(config.ActionCursorTop),
		k.Posts.GetString(config.ActionCursorBottom),
		k.Posts.GetString(config.ActionReaction),
		k.Posts.GetString(config.ActionUnreaction),
		k.Posts.GetString(config.ActionRepost),
		k.Posts.GetString(config.ActionUnrepost),
		k.Posts.GetString(config.ActionBookmark),
		k.Posts.GetString(config.ActionUnbookmark),
		k.Posts.GetString(config.ActionPost),
		k.Posts.GetString(config.ActionReply),
		k.Posts.GetString(config.ActionDelete),
		k.Posts.GetString(config.ActionOpenBrowser),
		k.Posts.GetString(config.ActionOpenUserPage),
		k.Posts.GetString(config.ActionCopyUrl),
	)

	user := fmt.Sprintf(
		`[User page]
  %-20s Follow
  %-20s Unfollow
  %-20s Mute
  %-20s Unmute
  %-20s Block
  %-20s Unblock
  %-20s Open user likes page
`,
		k.Posts.GetString(config.ActionFollow),
		k.Posts.GetString(config.ActionUnfollow),
		k.Posts.GetString(config.ActionMute),
		k.Posts.GetString(config.ActionUnmute),
		k.Posts.GetString(config.ActionBlock),
		k.Posts.GetString(config.ActionUnblock),
		k.Posts.GetString(config.ActionOpenUserPage),
	)

	text := global + view + page + post + user

	return &cli.Command{
		Name:     "keybindings",
		Short:    "Documentation for keybindings",
		Validate: cli.NoArgs(),
		SetFlag:  setOpenBackgroundFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			openBackground, _ := f.GetBool("background")
			return a.view.AddPage(newDocsPage("Keybindings", text), !openBackground)
		},
	}
}
