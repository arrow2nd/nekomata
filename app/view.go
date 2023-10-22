package app

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// view : ページの表示域
type view struct {
	flex            *tview.Flex
	pages           *tview.Pages
	pageItems       map[string]page
	tabBar          *tview.TextView
	tabs            []string
	currentTabIndex int
	textArea        *tview.TextArea
	modal           *tview.Modal
	mu              sync.Mutex
}

func newView() *view {
	v := &view{
		flex:            tview.NewFlex(),
		pages:           tview.NewPages(),
		pageItems:       map[string]page{},
		tabBar:          tview.NewTextView(),
		tabs:            []string{},
		currentTabIndex: 0,
		textArea:        tview.NewTextArea(),
		modal:           tview.NewModal(),
	}

	v.flex.
		SetDirection(tview.FlexRow).
		AddItem(v.pages, 0, 1, true).
		AddItem(v.textArea, 0, 0, false)

	tabBgColor := global.conf.Style.Tab.BackgroundColor.ToColor()
	v.tabBar.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft).
		SetHighlightedFunc(v.handleTabHighlight).
		SetTextStyle(tcell.StyleDefault.Background(tabBgColor))

	v.modal.
		AddButtons([]string{"No", "Yes"}).
		SetInputCapture(v.handleModalKeyEvent)

	v.textArea.
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1).
		SetBorder(true)

	return v
}

// SetInputCapture : キーハンドラを設定
func (v *view) SetInputCapture(f func(*tcell.EventKey) *tcell.EventKey) {
	v.flex.SetInputCapture(f)
}

// AddPage : ページを追加
func (v *view) AddPage(p page, focus bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	newTab := p.GetName()

	// ページが重複する場合、既にあるページに移動
	if _, ok := v.pageItems[newTab]; ok {
		tabIndex, found := find(v.tabs, func(tab string) bool { return tab == newTab })
		if !found {
			return fmt.Errorf("Failed to add page (%s)", newTab)
		}

		v.tabBar.Highlight(newTab)
		v.currentTabIndex = tabIndex

		return fmt.Errorf("This page already exists (%s)", newTab)
	}

	// ページ・タブを追加
	v.pageItems[newTab] = p
	v.pages.AddPage(newTab, p.GetPrimivite(), true, focus)
	v.addTab(newTab)

	// フォーカスが当たっているならタブをハイライト
	if focus {
		v.tabBar.Highlight(newTab)
		v.currentTabIndex = v.pages.GetPageCount() - 1
	}

	return nil
}

// Reset : リセット
func (v *view) Reset() {
	// ページを削除
	for id := range v.pageItems {
		v.pages.RemovePage(id)
	}
	v.pageItems = map[string]page{}

	// タブを削除
	v.removeTab("")
	v.tabBar.SetText("")
	v.currentTabIndex = 0
}

// CloseCurrentPage : 現在のページを削除
func (v *view) CloseCurrentPage() {
	// ページが1つのみなら削除しない
	if v.pages.GetPageCount() == 1 {
		global.SetErrorStatus("App", "last page cannot be closed")
		return
	}

	id, _ := v.pages.GetFrontPage()
	name := v.pageItems[id].GetName()

	// タブを削除
	v.removeTab(name)
	v.drawTab()

	// ページを削除
	v.pages.RemovePage(id)
	v.pageItems[id].OnDelete()
	delete(v.pageItems, id)

	// 前のタブを選択
	if v.currentTabIndex--; v.currentTabIndex < 0 {
		v.currentTabIndex = 0
	}

	v.tabBar.Highlight(v.tabs[v.currentTabIndex])
}

// ShowTextArea : テキストエリアを表示
func (v *view) ShowTextArea(hint string, onSubmit func(s string)) {
	f := func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		// 閉じる
		if key == tcell.KeyEsc {
			v.HiddenTextArea()
			return nil
		}

		// 入力確定
		if key == tcell.KeyCtrlP {
			v.HiddenTextArea()
			onSubmit(v.textArea.GetText())
			return nil
		}

		return event
	}

	v.textArea.
		SetText("", false).
		SetPlaceholder(hint).
		SetTitle(" Press ESC to close, press Ctrl-p to post ").
		SetInputCapture(f)

	v.flex.ResizeItem(v.textArea, 0, 1)

	global.RequestFocusPrimitive(v.textArea)
	global.SetDisableViewKeyEvent(true)
}

// HiddenTextArea : テキストエリアを非表示
func (v *view) HiddenTextArea() {
	v.flex.ResizeItem(v.textArea, 0, 0)

	global.RequestFocusPrimitive(v.pages)
	global.SetDisableViewKeyEvent(false)
}
