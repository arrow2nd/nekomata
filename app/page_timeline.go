package app

import (
	"context"
	"fmt"

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
	kind         timelineKind
	postList     *postList
	staremCancel context.CancelFunc
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
		Text:         &global.conf.Pref.Text,
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

	// タイムラインを取得
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

	// 表示に反映
	if err := t.postList.Update(posts); err != nil {
		return err
	}

	return err
}

func (t *timelinePage) OnDelete() {
	t.staremCancel()
}

func (t *timelinePage) StreamingRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.staremCancel = cancel

	// TODO: 非対応の場合も考慮したい
	opts := &sharedapi.StreamingTimelineOpts{
		Context: ctx,
		OnUpdate: func(post *sharedapi.Post) {
			if err := t.postList.Update([]*sharedapi.Post{post}); err != nil {
				label := fmt.Sprintf("stream (%s)", t.name)
				global.SetErrorStatus(label, err.Error())
			}
		},
		OnDelete: func(id string) {
			// TODO: 後で対応
		},
		OnError: func(err error) {
			label := fmt.Sprintf("stream (%s)", t.name)
			global.SetErrorStatus(label, err.Error())
		},
	}

	switch t.kind {
	case homeTimeline:
		return global.client.StreamingHomeTimeline(opts)
	case globalTimeline:
		return global.client.StreamingGlobalTimeline(opts)
	case localTimeline:
		return global.client.StreamingLocalTimeline(opts)
	}

	return nil
}
