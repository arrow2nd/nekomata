package app

import (
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/app/exit"
	"github.com/arrow2nd/nekomata/config"
	"github.com/rivo/tview"
)

// Global : 全体共有
type Global struct {
	currentUsername       string
	client                sharedapi.Client
	conf                  *config.Config
	isCLI                 bool
	chStatus              chan string
	chIndicator           chan string
	chPopupModal          chan *ModalOpts
	chExecCommand         chan string
	chInputCommand        chan string
	chFocusView           chan bool
	chFocusPrimitive      chan *tview.Primitive
	chDisableViewKeyEvent chan bool
}

var global = Global{
	conf:                  nil,
	isCLI:                 false,
	chStatus:              make(chan string, 1),
	chIndicator:           make(chan string, 1),
	chPopupModal:          make(chan *ModalOpts, 1),
	chExecCommand:         make(chan string, 1),
	chInputCommand:        make(chan string, 1),
	chFocusView:           make(chan bool, 1),
	chFocusPrimitive:      make(chan *tview.Primitive, 1),
	chDisableViewKeyEvent: make(chan bool, 1),
}

// SetStatus : ステータスメッセージを設定
func (g *Global) SetStatus(label, status string) {
	message := createStatusMessage(label, status)

	if g.isCLI {
		exit.OK(message)
	}

	g.chStatus <- message
}

// SetErrorStatus : エラーメッセージを設定
func (g *Global) SetErrorStatus(label, errStatus string) {
	if g.isCLI {
		exit.Error(createStatusMessage(label, errStatus), exit.CodeErr)
	}

	g.SetStatus("ERR: "+label, errStatus)
}

// SetIndicator : インジケータを設定
func (g *Global) SetIndicator(indicator string) {
	g.chIndicator <- indicator
}

// SetDisableViewKeyEvent : ビューのキーイベントを無効化
func (g *Global) SetDisableViewKeyEvent(b bool) {
	g.chDisableViewKeyEvent <- b
}

// ReqestPopupModal : モーダルの表示をリクエスト
func (g *Global) ReqestPopupModal(o *ModalOpts) {
	g.chPopupModal <- o
}

// RequestExecCommand : コマンドの実行をリクエスト
func (g *Global) RequestExecCommand(c string) {
	g.chExecCommand <- c
}

// RequestInputCommand : コマンドの入力をリクエスト
func (g *Global) RequestInputCommand(c string) {
	g.chInputCommand <- c
}

// RequestFocusPrimitive : 指定したプリミティブへのフォーカスを要求
func (g *Global) RequestFocusPrimitive(p tview.Primitive) {
	g.chFocusPrimitive <- &p
}

// RequestFocusView : ビューへのフォーカスを要求
func (g *Global) RequestFocusView() {
	g.chFocusView <- true
}
