package app

import "github.com/rivo/tview"

type docsPage struct {
	*basePage
	textView *tview.TextView
}

func newDocsPage(name, text string) *docsPage {
	tabName := replaceTabTemplateName(global.conf.Pref.Template.TabDocument, name)

	textView := tview.NewTextView().
		SetWordWrap(true).
		SetText(text)

	p := &docsPage{
		basePage: newBasePage(tabName),
		textView: textView,
	}

	p.SetFrame(p.textView)

	return p
}
