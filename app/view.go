package app

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// TabMovePrev : 前のタブに移動
	TabMovePrev int = -1
	// TabMoveNext : 次のタブに移動
	TabMoveNext int = 1
)

// ModalOpt : モーダルの設定
type ModalOpt struct {
	title  string
	text   string
	onDone func()
}

// tab : タブアイテム
type tab struct {
	id   string
	name string
}

// view : ページの表示域
type view struct {
	flex            *tview.Flex
	pages           *tview.Pages
	tabBar          *tview.TextView
	textArea        *tview.TextArea
	modal           *tview.Modal
	pageItems       map[string]page
	tabs            []string
	currentTabIndex int
	mu              sync.Mutex
}

func newView() *view {
	v := &view{
		flex:            tview.NewFlex(),
		pages:           tview.NewPages(),
		tabBar:          tview.NewTextView(),
		textArea:        tview.NewTextArea(),
		modal:           tview.NewModal(),
		pageItems:       map[string]page{},
		tabs:            []string{},
		currentTabIndex: 0,
	}

	v.flex.
		SetDirection(tview.FlexRow).
		AddItem(v.pages, 0, 1, true).
		AddItem(v.textArea, 0, 0, false)

	tabBgColor := shared.conf.Style.Tab.BackgroundColor.ToColor()
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

// createPageTag : ページ管理用のタグ文字列を作成
func createPageTag(id int) string {
	return fmt.Sprintf("page_%d", id)
}

// drawTab : タブを描画
func (v *view) drawTab() {
	v.tabBar.Clear()

	for i, tab := range v.tabs {
		fmt.Fprintf(v.tabBar, `[%s]["%s"] %s [""][-:-:-]`, shared.conf.Style.Tab.Text, tab, tab)

		// タブが2個以上あるならセパレータを挿入
		if i < len(v.tabs)-1 {
			fmt.Fprint(v.tabBar, shared.conf.Pref.Appearance.TabSeparator)
		}
	}
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

	// ページを追加
	v.pageItems[newTab] = p
	v.pages.AddPage(newTab, p.GetPrimivite(), true, focus)

	// フォーカスが当たっているならタブをハイライト
	if focus {
		v.tabBar.Highlight(newTab)
		v.currentTabIndex = v.pages.GetPageCount() - 1
	}

	// タブを追加
	v.tabs = append(v.tabs, newTab)
	v.drawTab()

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
	v.tabs = []string{}
	v.tabBar.SetText("")
	v.currentTabIndex = 0
}

// CloseCurrentPage : 現在のページを削除
func (v *view) CloseCurrentPage() {
	// ページが1つのみなら削除しない
	if v.pages.GetPageCount() == 1 {
		shared.SetErrorStatus("App", "last page cannot be closed")
		return
	}

	id, _ := v.pages.GetFrontPage()
	name := v.pageItems[id].GetName()

	newTabs := []string{}

	// タブを削除
	for _, tab := range v.tabs {
		if tab != name {
			newTabs = append(newTabs, tab)
		}
	}

	v.tabs = newTabs

	// 再描画して反映
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

// MoveTab : タブを移動する
func (v *view) MoveTab(move int) {
	maxTabIndex := v.pages.GetPageCount()
	if maxTabIndex == 0 {
		return
	}

	prevTabIndex := v.currentTabIndex
	nextTabIndex := v.currentTabIndex + move

	// 範囲内に丸める
	if nextTabIndex < 0 {
		nextTabIndex = maxTabIndex - 1
	} else if nextTabIndex >= maxTabIndex {
		nextTabIndex = 0
	}

	// 移動前と同じなら中断
	if nextTabIndex == prevTabIndex {
		return
	}

	v.currentTabIndex = nextTabIndex
	v.tabBar.Highlight(v.tabs[nextTabIndex])
}

// handleTabHighlight : タブがハイライトされたときのコールバック
func (v *view) handleTabHighlight(added, removed, remaining []string) {
	// FIXME: 1つ目のタブを追加した or startupCommand でタブを追加した時にエラーになる
	//        tview 内部の t.lineIndex の要素数が0の場合があるらしい

	// ハイライトされたタブまでスクロール
	// v.tabBar.ScrollToHighlight()

	// 前のページを非アクティブにする
	if len(removed) > 0 {
		if page, ok := v.pageItems[removed[0]]; ok {
			page.OnInactive()
		}
	}

	// ページを切り替え
	v.pages.SwitchToPage(added[0])
	v.pageItems[added[0]].OnActive()
}

// PopupModal : モーダルを表示
func (v *view) PopupModal(o *ModalOpt) {
	message := o.title

	// メッセージがあるなら追加
	if o.text != "" {
		message = fmt.Sprintf("%s\n\n%s", o.title, o.text)
	}

	f := func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			o.onDone()
		}

		v.pages.RemovePage("modal")
		shared.SetDisableViewKeyEvent(false)
	}

	v.modal.
		SetFocus(0).
		SetText(message).
		SetDoneFunc(f)

	v.pages.AddPage("modal", v.modal, true, true)

	shared.RequestFocusPrimitive(v.modal)
	shared.SetDisableViewKeyEvent(true)
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

	shared.RequestFocusPrimitive(v.textArea)
	shared.SetDisableViewKeyEvent(true)
}

// HiddenTextArea : テキストエリアを非表示
func (v *view) HiddenTextArea() {
	v.flex.ResizeItem(v.textArea, 0, 0)

	shared.RequestFocusPrimitive(v.pages)
	shared.SetDisableViewKeyEvent(false)
}
