package app

import (
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/config"
	"github.com/rivo/tview"
)

// Global : 全体共有
type Global struct {
	client                *sharedapi.Client
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
		exit(message)
	}

	go func() {
		g.chStatus <- message
	}()
}

// SetErrorStatus : エラーメッセージを設定
func (g *Global) SetErrorStatus(label, errStatus string) {
	if g.isCLI {
		exitError(createStatusMessage(label, errStatus), exitCodeErr)
	}

	g.SetStatus("ERR: "+label, errStatus)
}

// SetIndicator : インジケータを設定
func (g *Global) SetIndicator(indicator string) {
	go func() {
		g.chIndicator <- indicator
	}()
}

// SetDisableViewKeyEvent : ビューのキーイベントを無効化
func (g *Global) SetDisableViewKeyEvent(b bool) {
	go func() {
		g.chDisableViewKeyEvent <- b
	}()
}

// ReqestPopupModal : モーダルの表示をリクエスト
func (g *Global) ReqestPopupModal(o *ModalOpts) {
	go func() {
		g.chPopupModal <- o
	}()
}

// RequestExecCommand : コマンドの実行をリクエスト
func (g *Global) RequestExecCommand(c string) {
	go func() {
		g.chExecCommand <- c
	}()
}

// RequestInputCommand : コマンドの入力をリクエスト
func (g *Global) RequestInputCommand(c string) {
	go func() {
		g.chInputCommand <- c
	}()
}

// RequestFocusPrimitive : 指定したプリミティブへのフォーカスを要求
func (g *Global) RequestFocusPrimitive(p tview.Primitive) {
	go func() {
		g.chFocusPrimitive <- &p
	}()
}

// RequestFocusView : ビューへのフォーカスを要求
func (g *Global) RequestFocusView() {
	go func() {
		g.chFocusView <- true
	}()
}
