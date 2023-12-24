package cmd

import (
	"os"
	"os/exec"

	"github.com/arrow2nd/nekomata/config"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
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
		Action: func(ctx *cli.Context) error {
			fp := filepicker.New()
			fp.AllowedTypes = []string{".toml"}
			fp.CurrentDirectory, _ = config.GetConfigDir()

			m := configModel{
				editor:     ctx.String("editor"),
				filepicker: fp,
			}

			_, err := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
			return err
		},
	}
}

type configModel struct {
	editor     string
	filepicker filepicker.Model
	quitting   bool
	err        error
}

type editorFinishedMsg struct {
	err error
}

func (m configModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key := msg.String(); key == "ctrl+c" || key == "q" {
			m.quitting = true
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.quitting = true
		}
		return m, tea.Quit
	}

	// エディタを起動
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		cmd := exec.Command(m.editor, path)

		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return editorFinishedMsg{err}
		})
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	return m, cmd
}

func (m configModel) View() string {
	if m.err != nil {
		return m.filepicker.Styles.DisabledFile.Render(m.err.Error())
	}

	if m.quitting {
		return ""
	}

	return "\n Select configuration file to edit" + "\n\n" + m.filepicker.View() + "\n"
}
