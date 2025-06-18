package app

import "fmt"

const (
	// TabMovePrev : 前のタブに移動
	TabMovePrev int = -1
	// TabMoveNext : 次のタブに移動
	TabMoveNext int = 1
)

// MoveTab : タブを移動する
func (v *view) MoveTab(move int) {
	maxIndex := v.pages.GetPageCount()
	if maxIndex == 0 {
		return
	}

	prevIndex := v.currentTabIndex
	nextIndex := v.currentTabIndex + move

	// 範囲内に丸める
	if nextIndex < 0 {
		nextIndex = maxIndex - 1
	} else if nextIndex >= maxIndex {
		nextIndex = 0
	}

	// 移動前と同じなら中断
	if nextIndex == prevIndex {
		return
	}

	v.currentTabIndex = nextIndex
	v.tabBar.Highlight(v.tabs[nextIndex])
}

// addTab : タブを追加
func (v *view) addTab(name string) {
	v.tabs = append(v.tabs, name)
	v.drawTab()
}

// removeTab : タブを削除 (対象を指定しない場合すべて削除)
func (v *view) removeTab(name string) {
	newTabs := []string{}

	if name != "" {
		for _, tab := range v.tabs {
			if tab != name {
				newTabs = append(newTabs, tab)
			}
		}
	}

	v.tabs = newTabs
}

// drawTab : タブを描画
func (v *view) drawTab() {
	v.tabBar.Clear()

	for i, tab := range v.tabs {
		fmt.Fprintf(v.tabBar, `[%s]["%s"] %s [""][-:-:-]`, global.conf.Style.Tab.Text, tab, tab)

		// タブが2個以上あるならセパレータを挿入
		if i < len(v.tabs)-1 {
			fmt.Fprint(v.tabBar, global.conf.Pref.Appearance.TabSeparator)
		}
	}
}

// handleTabHighlight : タブがハイライトされたときのコールバック
func (v *view) handleTabHighlight(added, removed, remaining []string) {
	// 追加されたタブがない場合は何もしない（初期化時のエラー回避）
	if len(added) == 0 {
		return
	}

	// ハイライトされたタブまでスクロール（エラー回避のため条件付きで実行）
	if len(v.tabs) > 1 {
		defer func() {
			if r := recover(); r != nil {
				// tview内部のlineIndexエラーをキャッチ
				// エラーログを出力せずに処理を継続
			}
		}()
		v.tabBar.ScrollToHighlight()
	}

	// 前のページを非アクティブにする
	if len(removed) > 0 {
		if page, ok := v.pageItems[removed[0]]; ok {
			page.OnInactive()
		}
	}

	// ページの存在チェック
	if page, ok := v.pageItems[added[0]]; ok {
		// ページを切り替え
		v.pages.SwitchToPage(added[0])
		page.OnActive()
	}
}
