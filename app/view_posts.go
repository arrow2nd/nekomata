package app

import (
	"sync"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/app/layout"
	"github.com/arrow2nd/nekomata/config"
	"github.com/rivo/tview"
)

// cursorMove : カーソルの移動量
const (
	cursorMoveUp   int = -1
	cursorMoveDown int = 1
)

type postList struct {
	textView    *tview.TextView
	pinnedPosts []*sharedapi.Post
	posts       []*sharedapi.Post
	layout      *layout.Layout
	mu          sync.Mutex
}

func newPostsView(l *layout.Layout) (*postList, error) {
	p := &postList{
		textView:    tview.NewTextView(),
		pinnedPosts: []*sharedapi.Post{},
		posts:       []*sharedapi.Post{},
		layout:      l,
	}

	p.layout.Writer = p.textView

	p.textView.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true)

	if err := p.setKeybindings(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *postList) setKeybindings() error {
	handlers := map[string]func(){
		config.ActionScrollUp: func() {
			r, c := p.textView.GetScrollOffset()
			p.textView.ScrollTo(r+1, c)
		},
		config.ActionScrollDown: func() {
			r, c := p.textView.GetScrollOffset()
			p.textView.ScrollTo(r-1, c)
		},
		config.ActionCursorUp: func() {
			p.moveCursor(cursorMoveUp)
		},
		config.ActionCursorDown: func() {
			p.moveCursor(cursorMoveDown)
		},
		config.ActionCursorTop: func() {
			p.scrollToPost(0)
		},
		config.ActionCursorBottom: func() {
			lastIndex := p.GetPostsCount() - 1
			p.scrollToPost(lastIndex)
		},
		config.ActionTweetLike: func() {
		},
		config.ActionTweetUnlike: func() {
		},
		config.ActionTweetRetweet: func() {
		},
		config.ActionTweetUnretweet: func() {
		},
		config.ActionTweetDelete: func() {
		},
		config.ActionUserFollow: func() {
		},
		config.ActionUserUnfollow: func() {
		},
		config.ActionUserBlock: func() {
		},
		config.ActionUserUnblock: func() {
		},
		config.ActionUserMute: func() {
		},
		config.ActionUserUnmute: func() {
		},
		config.ActionOpenUserPage: func() {
		},
		config.ActionOpenUserLikes: func() {
		},
		config.ActionTweet: func() {
		},
		config.ActionQuote: func() {
		},
		config.ActionReply: func() {
		},
		config.ActionOpenBrowser: func() {
		},
		config.ActionCopyUrl: func() {
		},
	}

	c, err := global.conf.Pref.Keybindings.TweetView.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	p.textView.SetInputCapture(c.Capture)

	return nil
}

// moveCursor : カーソルを移動
func (p *postList) moveCursor(c int) {
	idx := getHighlightId(p.textView.GetHighlights())
	if idx == -1 {
		return
	}

	p.scrollToPost(idx + int(c))
}

// scrollToPost : 指定したポストまでスクロール
func (p *postList) scrollToPost(i int) {
	// 範囲内に丸める
	if max := p.GetPostsCount(); i < 0 {
		i = 0
	} else if i >= max {
		i = max - 1
	}

	p.textView.Highlight(layout.CreatePostHighlightTag(i))
	p.textView.ScrollToHighlight()
}

// GetPostsCount : ポスト数を取得
func (p *postList) GetPostsCount() int {
	c := len(p.posts)

	if l := len(p.pinnedPosts); l > 0 {
		c += l
	}

	return c
}

// GetSinceId : 最新のポストのIDを取得
func (p *postList) GetSinceId() string {
	if len(p.posts) == 0 {
		return ""
	}

	return p.posts[0].ID
}

// SetPinned : ピン留めを登録
func (p *postList) SetPinned(pinned []*sharedapi.Post) {
	p.pinnedPosts = []*sharedapi.Post{}

	for i, post := range pinned {
		p.pinnedPosts[i] = post
	}
}

// Update : ポストを更新
func (p *postList) Update(posts []*sharedapi.Post) error {
	addedPostsCount := p.addPosts(posts)
	cursorPos := p.getCurrentCursorPos()

	// 先頭以外のツイートを選択中の場合、更新後もそのツイートを選択したままにする
	// NOTE: "先頭以外" なのは、ストリームモードで放置した時にカーソルが段々下に下がってしまうのを防ぐため
	if cursorPos != 0 {
		cursorPos += addedPostsCount
	}

	return p.draw(cursorPos)
}

// getCurrentCursorPos : 現在のカーソル位置を取得
func (p *postList) getCurrentCursorPos() int {
	pos := getHighlightId(p.textView.GetHighlights())

	if pos == -1 {
		pos = 0
	}

	return pos
}

// addPosts : ポストを追加
func (p *postList) addPosts(posts []*sharedapi.Post) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	size := len(p.posts)
	addSize := len(posts)
	allSize := size + addSize
	maxSize := global.conf.Pref.Feature.AccmulateTweetsLimit

	// 最大蓄積数を超えていたら古いものから削除
	if allSize > maxSize {
		size -= allSize - maxSize
	}

	p.posts = append(posts, p.posts[:size]...)

	return addSize
}

// DeletePost : ポストを削除
func (p *postList) DeletePost(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, ok := find(p.posts, func(c *sharedapi.Post) bool {
		// リポスト元のIDを参照
		if ref := c.Reference; ref != nil {
			return ref.ID == id
		}

		return c.ID == id
	})

	if !ok {
		return nil
	}

	// i番目の要素を削除
	p.posts = p.posts[:i+copy(p.posts[i:], p.posts[i+1:])]

	// 再描画して反映
	return p.draw(p.getCurrentCursorPos())
}

// draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (p *postList) draw(cursorPos int) error {
	// icon := global.conf.Pref.Icon
	appearance := global.conf.Pref.Appearance
	p.textView.
		SetTextAlign(tview.AlignLeft).
		Clear()

	// 表示するポストが無いなら描画を中断
	if p.GetPostsCount() == 0 {
		p.DrawMessage(global.conf.Pref.Text.NoPosts)
		return nil
	}

	contents := p.posts

	// ピン留めがある場合、先頭に追加
	if len(p.pinnedPosts) > 0 {
		contents = append(p.pinnedPosts, p.posts...)
	}

	for i, post := range contents {
		// var quoted *sharedapi.Post
		// annotation := ""

		// 参照ツイートを確認
		// for _, rc := range content.ReferencedTweets {
		// 	switch rc.Reference.Type {
		// 	case "retweeted":
		// 		annotation += createAnnotation("RT by", content.Author)
		// 		content = rc.TweetDictionary
		// 	case "replied_to":
		// 		annotation += createAnnotation("Reply to", rc.TweetDictionary.Author)
		// 	case "quoted":
		// 		quotedTweet = rc.TweetDictionary
		// 	}
		// }

		// ピン留めツイート
		// if i == 0 && t.pinned != nil {
		// 	annotation += fmt.Sprintf("[gray:-:-]%s Pinned Tweet[-:-:-]", icon.Pinned)
		// }
		if err := p.layout.Post(i, post); err != nil {
			return err
		}

		// fmt.Fprintf(p.textView, "%s\n", post.Text)

		// 引用元ツイートを表示
		// if quotedTweet != nil {
		// 	if !appearance.HideQuoteTweetSeparator {
		// 		fmt.Fprintln(t.view, createSeparator(appearance.QuoteTweetSeparator, width))
		// 	}
		//
		// 	fmt.Fprintln(t.view, createTweetLayout("", quotedTweet, -1, width))
		// }

		// セパレータを挿入しない
		if appearance.HideTweetSeparator {
			continue
		}

		// 末尾のツイート以外ならセパレータを挿入
		if i < p.GetPostsCount()-1 {
			p.layout.PrintSeparator(appearance.TweetSeparator)
		}
	}

	p.scrollToPost(cursorPos)

	return nil
}

// DrawMessage : ビューにメッセージを表示
func (p *postList) DrawMessage(s string) {
	p.textView.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}
