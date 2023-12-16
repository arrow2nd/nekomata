package app

import (
	"os"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/arrow2nd/nekomata/cli"
	"github.com/arrow2nd/nekomata/config"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

// App : アプリケーション
type App struct {
	cmd         *cli.Command
	view        *view
	statusBar   *statusBar
	commandLine *commandLine
}

// New : 新規作成
func New() *App {
	global.app = tview.NewApplication()

	return &App{
		cmd:         newCmd(),
		view:        nil,
		statusBar:   nil,
		commandLine: nil,
	}
}

// Init : 初期化
func (a *App) Init() error {
	if err := a.loadConfig(); err != nil {
		return err
	}

	// 実行時の引数をパース
	isSkipLogin, username, err := a.parseRuntimeArgs()
	if err != nil {
		return err
	}

	// アカウントにログイン
	if !isSkipLogin {
		if err := login(username); err != nil {
			return err
		}
	}

	// コマンドを初期化
	a.initCommands()

	// コマンドラインモードならUIの初期化をスキップ
	if global.isCLI {
		return nil
	}

	// UI準備
	a.setAppStyles()
	a.view = newView()
	a.statusBar = newStatusBar()
	a.commandLine = newCommandLine()

	// 日本語環境等での罫線の乱れ対策
	// LINK: https://github.com/mattn/go-runewidth/issues/14
	runewidth.DefaultCondition.EastAsianWidth = !global.conf.Pref.Feature.IsLocaleCJK

	// Ctrl+K/Jの再マッピングを無効化
	cbind.UnifyEnterKeys = false

	// キーバインドを設定
	if err := a.setGlobalKeybindings(); err != nil {
		return err
	}
	if err := a.setViewKeybindings(); err != nil {
		return err
	}

	// ステータスバー初期化
	a.statusBar.Init()
	a.statusBar.DrawAccountInfo()

	// コマンドライン初期化
	a.commandLine.Init()
	a.initAutocomplate()

	// 画面レイアウト
	layout := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetBorders(false).
		AddItem(a.view.tabBar, 0, 0, 1, 1, 0, 0, false).
		AddItem(a.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(a.commandLine.inputField, 3, 0, 1, 1, 0, 0, false).
		AddItem(a.view.flex, 1, 0, 1, 1, 0, 0, true)

	global.app.SetRoot(layout, true)

	a.execStartupCommands()

	return nil
}

// loadConfig : 設定を読み込む
func (a *App) loadConfig() error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	global.conf = conf

	// 環境設定
	if err := global.conf.LoadPreferences(); err != nil {
		return err
	}

	// スタイル定義
	if err := global.conf.LoadStyle(); err != nil {
		return err
	}

	// 資格情報
	return global.conf.LoadCred()
}

// parseRuntimeArgs : 実行時の引数をパースして、ログインユーザを返す
func (a *App) parseRuntimeArgs() (bool, string, error) {
	f := a.cmd.NewFlagSet()
	f.ParseErrorsWhitelist.UnknownFlags = true

	if err := f.Parse(os.Args[1:]); err != nil {
		return false, "", err
	}

	// ログインをスキップするか
	arg := f.Arg(0)
	isSkipLogin := f.Changed("help") || f.Changed("version") || arg == "edit"

	// コマンドラインモードか
	global.isCLI = f.NArg() > 0 || isSkipLogin

	account, _ := f.GetString("account")

	return isSkipLogin, account, nil
}

// setAppStyles : アプリ全体のスタイルを設定
func (a *App) setAppStyles() {
	app := global.conf.Style.App

	bgColor := app.BackgroundColor.ToColor()
	textColor := app.TextColor.ToColor()
	borderColor := app.BorderColor.ToColor()

	// 背景色
	tview.Styles.PrimitiveBackgroundColor = bgColor
	tview.Styles.ContrastBackgroundColor = bgColor

	// TODO: Dropdownの背景色
	// tview.Styles.MoreContrastBackgroundColor = app.BackgroundColor.ToColor()

	// テキスト色
	tview.Styles.PrimaryTextColor = textColor
	tview.Styles.ContrastSecondaryTextColor = textColor
	tview.Styles.TitleColor = textColor
	tview.Styles.TertiaryTextColor = app.SubTextColor.ToColor()

	// ボーダー色
	tview.Styles.BorderColor = borderColor
	tview.Styles.GraphicsColor = borderColor

	// ボーダー
	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft
}

// setGlobalKeybindings : アプリ全体のキーバインドを設定
func (a *App) setGlobalKeybindings() error {
	handlers := map[string]func(){
		config.ActionQuit: func() {
			a.quitApp()
		},
	}

	c, err := global.conf.Pref.Keybindings.Global.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	global.app.SetInputCapture(a.warpKeyEventHandler(c))

	return nil
}

// setViewKeybindings : ビューのキーバインドを設定
func (a *App) setViewKeybindings() error {
	handlers := map[string]func(){
		config.ActionSelectPrevTab: func() {
			a.view.MoveTab(TabMovePrev)
		},
		config.ActionSelectNextTab: func() {
			a.view.MoveTab(TabMoveNext)
		},
		config.ActionRedraw: func() {
			global.app.Sync()
		},
		config.ActionFocusCmdLine: func() {
			global.app.SetFocus(a.commandLine.inputField)
		},
		config.ActionShowHelp: func() {
			global.RequestExecCommand("docs keybindings")
		},
		config.ActionClosePage: func() {
			a.view.CloseCurrentPage()
		},
	}

	c, err := global.conf.Pref.Keybindings.View.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	a.view.SetInputCapture(a.warpKeyEventHandler(c))

	return nil
}

// warpKeyEventHandler : イベントハンドラのラップ関数
func (a *App) warpKeyEventHandler(c *cbind.Configuration) func(*tcell.EventKey) *tcell.EventKey {
	return func(ev *tcell.EventKey) *tcell.EventKey {
		// 操作が無効
		if !global.enableAppKeybind {
			return ev
		}

		return c.Capture(ev)
	}
}

// initCommands : コマンドを初期化
func (a *App) initCommands() {
	a.cmd.AddCommand(
		a.newAccountCmd(),
		a.newQuitCmd(),
		a.newDocsCmd(),
		a.newEditCmd(),
		a.newPostCmd(),
		a.newTimelineCmd(),
	)

	if global.isCLI {
		return
	}

	// ヘルプの出力を新規ページに割り当てる
	a.cmd.Help = func(c *cli.Command, h string) {
		a.view.AddPage(newDocsPage(c.Name, h), true)
	}
}

// initAutocomplate : 入力補完を初期化
func (a *App) initAutocomplate() {
	cmds := a.cmd.GetChildrenNames(true)

	if err := a.commandLine.SetAutocompleteItems(cmds); err != nil {
		global.SetErrorStatus("Init - CommandLine", err.Error())
	}
}

// execStartupCommands : 起動時に実行するコマンドを一括で実行
func (a *App) execStartupCommands() {
	for _, c := range global.conf.Pref.Feature.StartupCmds {
		if err := a.ExecCommnad(c); err != nil {
			global.SetErrorStatus("Command", err.Error())
		}
	}
}

// ExecCommnad : コマンドを実行
func (a *App) ExecCommnad(cmd string) error {
	args, err := split(cmd)
	if err != nil {
		return err
	}

	return a.cmd.Execute(args)
}

// Run : アプリを実行
func (a *App) Run() error {
	// コマンドラインモード
	if global.isCLI {
		return a.cmd.Execute(os.Args[1:])
	}

	go a.eventReceiver()

	return global.app.Run()
}

// eventReceiver : イベントレシーバ
func (a *App) eventReceiver() {
	for {
		select {
		case status := <-global.chStatus:
			// ステータスメッセージを表示
			global.app.QueueUpdateDraw(func() {
				a.commandLine.ShowStatusMessage(status)
			})

		case indicator := <-global.chIndicator:
			// インジケータを更新
			global.app.QueueUpdateDraw(func() {
				a.statusBar.DrawPageIndicator(indicator)
			})

		case opt := <-global.chPopupModal:
			// モーダルを表示
			global.app.QueueUpdateDraw(func() {
				a.view.PopupModal(opt)
			})

		case cmd := <-global.chExecCommand:
			// コマンドを実行`
			if err := a.ExecCommnad(cmd); err != nil {
				global.SetErrorStatus("Command", err.Error())
			}

		case cmd := <-global.chInputCommand:
			// コマンドを入力
			global.app.SetFocus(a.commandLine.inputField)
			global.app.QueueUpdateDraw(func() {
				a.commandLine.SetText(cmd)
			})

		case <-global.chFocusView:
			// ビューにフォーカス
			global.app.QueueUpdateDraw(func() {
				global.app.SetFocus(a.view.flex)
			})
		}
	}
}

// quitApp : アプリを終了
func (a *App) quitApp() {
	// 確認画面が不要ならそのまま終了
	if !global.conf.Pref.Confirm["quit"] {
		global.app.Stop()
		return
	}

	a.view.PopupModal(&ModalOpts{
		title:  "Do you want to quit the app?",
		onDone: global.app.Stop,
	})
}
