package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *view) ShowPostForm(title string, opts *submitPostOpts) {
	close := func() {
		v.pages.RemovePage("postForm")
		global.enableAppKeybind = true
	}

	visibilities := global.client.GetVisibilityList()

	defaultIndex, _ := find(visibilities, func(v string) bool {
		return v == opts.post.Visibility
	})

	if defaultIndex < 0 {
		defaultIndex = 0
	}

	form := tview.NewForm().
		AddDropDown("Visibility", global.client.GetVisibilityList(), defaultIndex, func(option string, _ int) {
			opts.post.Visibility = option
		}).
		AddCheckbox("NSFW", opts.post.Sensitive, func(checked bool) {
			opts.post.Sensitive = checked
		}).
		AddInputField("Contemt warning", "", 0, nil, nil).
		AddTextArea("Body", opts.post.Text, 0, 5, 0, func(text string) {
			opts.post.Text = text
		}).
		AddButton("Cancel", close).
		AddButton("Submit", func() {
			submitPost(opts)
			close()
		}).
		SetFieldBackgroundColor(tcell.GetColor("#1c1c1c")).
		SetCancelFunc(close)

	form.
		SetTitle(" Compose " + title + " ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	modal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(form, 0, 1, true)

	global.RequestQueueUpdateDraw(func() {
		v.pages.AddPage("postForm", modal, true, true)
		global.enableAppKeybind = false
	})
}
