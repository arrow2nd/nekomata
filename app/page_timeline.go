package app

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

	return &timelinePage{
		basePage: newBasePage(tabName),
		kind:     kind,
		posts:    postsView,
	}, nil
}
