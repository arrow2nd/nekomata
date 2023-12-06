package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *view) ShowPostForm(onSubmit func()) {
	close := func() {
		v.pages.RemovePage("postForm")
		global.enableAppKeybind = true
	}

	form := tview.NewForm().
		AddDropDown("Visibility", global.client.GetVisibilityList(), 0, nil).
		AddCheckbox("NSFW", false, nil).
		AddInputField("Contemt warning", "", 0, nil, nil).
		AddTextArea("Body", "", 0, 5, 0, nil).
		AddButton("Cancel", close).
		AddButton("Submit", nil).
		SetFieldBackgroundColor(tcell.GetColor("#1c1c1c")).
		SetCancelFunc(close)

	form.
		SetTitle(" New post ").
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
