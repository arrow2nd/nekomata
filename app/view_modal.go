package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// ModalOpts : モーダルの設定
type ModalOpts struct {
	title  string
	text   string
	onDone func()
}

// PopupModal : モーダルを表示
func (v *view) PopupModal(opts *ModalOpts) {
	message := opts.title

	// テキストがあるなら追加
	if opts.text != "" {
		message = fmt.Sprintf("%s\n\n%s", opts.title, opts.text)
	}

	f := func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			opts.onDone()
		}

		v.pages.RemovePage("modal")
		global.enableAppKeybind = true
	}

	v.modal.
		SetFocus(0).
		SetText(message).
		SetDoneFunc(f)

	v.pages.AddPage("modal", v.modal, true, true)
	global.enableAppKeybind = false
}

// handleModalKeyEvent : モーダルのキーイベントハンドラ
func (v *view) handleModalKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// hjを左キーの入力イベントに置換
	if keyRune == 'h' || keyRune == 'j' {
		return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
	}

	// klを右キーの入力イベントに置換
	if keyRune == 'k' || keyRune == 'l' {
		return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
	}

	return event
}
