package app

import (
	"context"
	"errors"
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
	streamCancel context.CancelFunc
}

func newTimelinePage(kind timelineKind) (*timelinePage, error) {
	tabTemplate := global.conf.Pref.Template.TabHome

	switch kind {
	case globalTimeline:
		tabTemplate = global.conf.Pref.Template.TabGlobal
	case localTimeline:
		tabTemplate = global.conf.Pref.Template.TabLocal
	}

	layout := &layout.Layout{
		Template:   &global.conf.Pref.Template,
		Appearance: &global.conf.Pref.Appearance,
		Text:       &global.conf.Pref.Text,
		Style:      global.conf.Style,
	}

	postsView, err := newPostsView(layout)
	if err != nil {
		return nil, err
	}

	page := &timelinePage{
		basePage:     newBasePage(tabTemplate),
		kind:         kind,
		postList:     postsView,
		streamCancel: nil,
	}

	page.SetFrame(postsView.textView)

	handler, err := createCommonPageKeyHandler(page)
	if err != nil {
		return nil, err
	}

	page.frame.SetInputCapture(handler)

	return page, nil
}

func (t *timelinePage) Load() error {
	var (
		sinceID = t.postList.GetSinceId()
		limit   = global.conf.Pref.Feature.LoadPostCount
		posts   []*sharedapi.Post
		err     error
	)

	// ストリーミング中は手動で読み込みさせない
	if t.streamCancel != nil {
		return errors.New("cannot load manually due to streaming")
	}

	// 読み込み中
	t.postList.DrawMessage(global.conf.Pref.Text.Loading)

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

	// 読み込み完了
	if n := len(posts); n > 0 {
		global.SetStatus(t.name, fmt.Sprintf("%d posts loaded", n))
	} else {
		global.SetStatus(t.name, global.conf.Pref.Text.NoPosts)
	}

	// ストリーミングを開始
	if global.client.IsStreamingSupported() && t.streamCancel == nil {
		go t.Streaming()
	}

	// 表示に反映
	return t.postList.Update(posts)
}

func (t *timelinePage) OnDelete() {
	if t.streamCancel != nil {
		t.streamCancel()
	}
}

func (t *timelinePage) Streaming() {
	ctx, cancel := context.WithCancel(context.Background())
	t.streamCancel = cancel

	defer func() {
		t.streamCancel = nil
	}()

	opts := &sharedapi.StreamingTimelineOpts{
		Context: ctx,
		OnUpdate: func(post *sharedapi.Post) {
			if err := t.postList.Update([]*sharedapi.Post{post}); err != nil {
				t.showError(err)
			}
		},
		OnDelete: func(id string) {
			// TODO: 後で対応
		},
		OnError: func(err error) {
			t.showError(err)
		},
	}

	var err error = nil

	switch t.kind {
	case homeTimeline:
		err = global.client.StreamingHomeTimeline(opts)
	case globalTimeline:
		err = global.client.StreamingGlobalTimeline(opts)
	case localTimeline:
		err = global.client.StreamingLocalTimeline(opts)
	default:
		err = fmt.Errorf("invalid kind: %s", t.kind)
	}

	if err != nil {
		t.showError(err)
	}
}

func (t *timelinePage) showError(e error) {
	label := fmt.Sprintf("streaming (%s)", t.name)
	global.SetErrorStatus(label, e.Error())
}
