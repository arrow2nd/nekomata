package app

import (
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/app/layout"
)

type timelineKind string

const (
	homeTimeline   timelineKind = "home"
	globalTimeline timelineKind = "global"
	localTimeline  timelineKind = "local"
)

type timelinePage struct {
	*basePage
	kind     timelineKind
	postList *postList
}

func newTimelinePage(kind timelineKind) (*timelinePage, error) {
	tabName := global.conf.Pref.Text.TabHome

	switch kind {
	case globalTimeline:
		tabName = global.conf.Pref.Text.TabGlobal
	case localTimeline:
		tabName = global.conf.Pref.Text.TabLocal
	}

	layout := &layout.Layout{
		Width:        getWindowWidth(),
		Template:     &global.conf.Pref.Template,
		Appearancene: &global.conf.Pref.Appearance,
		Icon:         &global.conf.Pref.Icon,
		Style:        global.conf.Style,
	}

	postsView, err := newPostsView(layout)
	if err != nil {
		return nil, err
	}

	page := &timelinePage{
		basePage: newBasePage(tabName),
		kind:     kind,
		postList: postsView,
	}

	page.SetFrame(postsView.textView)

	return page, nil
}

func (t *timelinePage) Load() error {
	var (
		sinceID = t.postList.GetSinceId()
		limit   = global.conf.Pref.Feature.LoadTweetsLimit
		posts   []*sharedapi.Post
		err     error
	)

	switch t.kind {
	case homeTimeline:
		posts, err = global.client.GetHomeTimeline(sinceID, limit)
	case globalTimeline:
		posts, err = global.client.GetGlobalTimeline(sinceID, limit)
	case localTimeline:
		posts, err = global.client.GetLocalTimeline(sinceID, limit)
	}

	if err != nil {
		return err
	}

	return t.postList.Update(posts)
}
