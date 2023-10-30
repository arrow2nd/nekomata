package app

import (
	"fmt"

	"github.com/arrow2nd/nekomata/app/layout"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type statusBar struct {
	flex      *tview.Flex
	leftView  *tview.TextView
	rightView *tview.TextView
}

func newStatusBar() *statusBar {
	return &statusBar{
		flex:      tview.NewFlex(),
		leftView:  tview.NewTextView(),
		rightView: tview.NewTextView(),
	}
}

// Init : 初期化
func (s *statusBar) Init() {
	bgColor := global.conf.Style.StatusBar.BackgroundColor.ToColor()

	s.leftView.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextStyle(tcell.StyleDefault.Background(bgColor))

	s.rightView.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextStyle(tcell.StyleDefault.Background(bgColor))

	s.flex.
		SetDirection(tview.FlexColumn).
		AddItem(s.leftView, 0, 1, false).
		AddItem(s.rightView, 0, 1, false)
}

// DrawAccountInfo : ログイン中のアカウント情報を描画
func (s *statusBar) DrawAccountInfo() {
	s.leftView.Clear()
	fmt.Fprintf(s.leftView, layout.CreateStyledText(global.conf.Style.StatusBar.Text, " @"+global.currentUsername, ""))
}

// DrawPageIndicator : 現在のページのインジケータを描画
func (s *statusBar) DrawPageIndicator(d string) {
	s.rightView.Clear()
	fmt.Fprintf(s.rightView, layout.CreateStyledText(global.conf.Style.StatusBar.Text, d, ""))
}
