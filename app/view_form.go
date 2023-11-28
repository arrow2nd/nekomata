package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *view) ShowPostForm(onSubmit func()) {
	form := tview.NewForm().
		AddDropDown("Visibility", global.client.GetVisibilityList(), 0, nil). // FIXME: focus時に "l" を入力するとフリーズするかも
		AddCheckbox("NSFW", false, nil).
		AddTextArea("Post", "", 0, 0, 0, nil).
		AddButton("Cancel", func() {
			v.HiddenPostForm()
		}).
		AddButton("Submit", func() {
			v.HiddenPostForm()
		}).
		SetCancelFunc(func() {
			v.HiddenPostForm()
		}).
		SetFieldBackgroundColor(tcell.GetColor("#1c1c1c"))

	form.
		SetBorder(true).
		SetTitle(" Press ESC to close ").
		SetTitleAlign(tview.AlignLeft)

	formPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(form, 0, 1, true)

	global.app.QueueUpdateDraw(func() {
		v.pages.AddPage("postForm", formPage, true, true)
	})

	global.SetDisableViewKeyEvent(true)
}

func (v *view) HiddenPostForm() {
	v.pages.RemovePage("postForm")
	global.SetDisableViewKeyEvent(false)
}
