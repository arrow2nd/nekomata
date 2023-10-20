package app

import "github.com/arrow2nd/nekomata/api/sharedapi"

type timelineKind string

const (
	homeTimeline   timelineKind = "home"
	globalTimeline timelineKind = "global"
	localTimeline  timelineKind = "local"
)

type timelinePage struct {
	*basePage
	kind  timelineKind
	posts *posts
}

func newTimelinePage(kind timelineKind) (*timelinePage, error) {
	tabName := global.conf.Pref.Text.TabHome

	switch kind {
	case globalTimeline:
		tabName = global.conf.Pref.Text.TabGlobal
	case localTimeline:
		tabName = global.conf.Pref.Text.TabLocal
	}

	postsView, err := newPostsView()
	if err != nil {
		return nil, err
	}

	page := &timelinePage{
		basePage: newBasePage(tabName),
		kind:     kind,
		posts:    postsView,
	}

	page.SetFrame(postsView.view)

	return page, nil
}

func (t *timelinePage) Load() error {
	var (
		sinceID = t.posts.GetSinceId()
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

	t.posts.Update(posts)

	return nil
}
