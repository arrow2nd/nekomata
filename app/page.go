package app

import (
	"fmt"
	"strings"

	"github.com/arrow2nd/nekomata/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type page interface {
	GetName() string
	GetPrimitive() tview.Primitive
	Load() error
	OnActive()
	OnInactive()
	OnDelete()
}

type basePage struct {
	page
	name      string
	indicator string
	frame     *tview.Frame
	isActive  bool
}

func newBasePage(name string) *basePage {
	return &basePage{
		name:      truncate(name, global.conf.Pref.Appearance.TabMaxWidth),
		indicator: "",
		frame:     nil,
		isActive:  false,
	}
}

// GetName : ページ名を取得
func (b *basePage) GetName() string {
	return b.name
}

// GetPrimitive : プリミティブを取得
func (b *basePage) GetPrimitive() tview.Primitive {
	return b.frame
}

// SetFrame : フレームを設定
func (b *basePage) SetFrame(p tview.Primitive) {
	b.frame = tview.NewFrame(p)
	b.frame.SetBorders(1, 1, 0, 0, 1, 1)
}

// Load : 読み込み
func (b *basePage) Load() error {
	return nil
}

// OnActive : ページがアクティブになった
func (b *basePage) OnActive() {
	b.isActive = true

	// 以前のインジケータの内容を反映
	global.SetIndicator(b.indicator)
}

// OnInactive : ページが非アクティブになった
func (b *basePage) OnInactive() {
	b.isActive = false
}

// OnDelete : ページが破棄された
func (b *basePage) OnDelete() {}

// createCommonPageKeyHandler : ページ共通のキーハンドラを作成
func createCommonPageKeyHandler(p page) (func(*tcell.EventKey) *tcell.EventKey, error) {
	handler := map[string]func(){
		config.ActionReloadPage: func() {
			go func() {
				if err := p.Load(); err != nil {
					label := fmt.Sprintf("load (%s)", p.GetName())
					global.SetErrorStatus(label, err.Error())
				}
			}()
		},
	}

	c, err := global.conf.Pref.Keybindings.Page.MappingEventHandler(handler)
	if err != nil {
		return nil, err
	}

	return c.Capture, nil
}

// replaceTabTemplateName : タブテンプレートの名前を置換する
func replaceTabTemplateName(template, name string) string {
	return strings.Replace(template, "{{ name }}", name, 1)
}
