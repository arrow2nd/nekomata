package app

import (
	"fmt"
	"sync"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/config"
	"github.com/rivo/tview"
)

// cursorMove : カーソルの移動量
const (
	cursorMoveUp   int = -1
	cursorMoveDown int = 1
)

type posts struct {
	view     *tview.TextView
	pinned   []*sharedapi.Post
	contents []*sharedapi.Post
	mu       sync.Mutex
}

func newPostsView() (*posts, error) {
	p := &posts{
		view:     tview.NewTextView(),
		pinned:   []*sharedapi.Post{},
		contents: []*sharedapi.Post{},
	}

	p.view.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true)

	p.view.SetHighlightedFunc(func(_, _, _ []string) {
		p.view.ScrollToHighlight()
	})

	if err := p.setKeybindings(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *posts) setKeybindings() error {
	handlers := map[string]func(){
		config.ActionScrollUp: func() {
			r, c := p.view.GetScrollOffset()
			p.view.ScrollTo(r+1, c)
		},
		config.ActionScrollDown: func() {
			r, c := p.view.GetScrollOffset()
			p.view.ScrollTo(r-1, c)
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
		// config.ActionTweetLike: func() {
		// 	p.actionForTweet(tweetActionLike)
		// },
		// config.ActionTweetUnlike: func() {
		// 	p.actionForTweet(tweetActionUnlike)
		// },
		// config.ActionTweetRetweet: func() {
		// 	p.actionForTweet(tweetActionRetweet)
		// },
		// config.ActionTweetUnretweet: func() {
		// 	p.actionForTweet(tweetActionUnretweet)
		// },
		// config.ActionTweetDelete: func() {
		// 	p.actionForTweet(tweetActionDelete)
		// },
		// config.ActionUserFollow: func() {
		// 	p.actionForUser(userActionFollow)
		// },
		// config.ActionUserUnfollow: func() {
		// 	p.actionForUser(userActionUnfollow)
		// },
		// config.ActionUserBlock: func() {
		// 	p.actionForUser(userActionBlock)
		// },
		// config.ActionUserUnblock: func() {
		// 	p.actionForUser(userActionUnblock)
		// },
		// config.ActionUserMute: func() {
		// 	p.actionForUser(userActionMute)
		// },
		// config.ActionUserUnmute: func() {
		// 	p.actionForUser(userActionUnmute)
		// },
		// config.ActionOpenUserPage: func() {
		// 	p.openUserPage()
		// },
		// config.ActionOpenUserLikes: func() {
		// 	p.openUserLikes()
		// },
		// config.ActionTweet: func() {
		// 	shared.RequestExecCommand("tweet")
		// },
		// config.ActionQuote: func() {
		// 	p.insertQuoteCommand()
		// },
		// config.ActionReply: func() {
		// 	p.insertReplyCommand()
		// },
		// config.ActionOpenBrowser: func() {
		// 	p.openBrower()
		// },
		// config.ActionCopyUrl: func() {
		// 	p.copyLinkToClipBoard()
		// },
	}

	c, err := global.conf.Pref.Keybindings.TweetView.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	p.view.SetInputCapture(c.Capture)

	return nil
}

// moveCursor : カーソルを移動
func (p *posts) moveCursor(c int) {
	idx := getHighlightId(p.view.GetHighlights())
	if idx == -1 {
		return
	}

	p.scrollToPost(idx + int(c))
}

// scrollToPost : 指定したポストまでスクロール
func (p *posts) scrollToPost(i int) {
	// 範囲内に丸める
	if max := p.GetPostsCount(); i < 0 {
		i = 0
	} else if i >= max {
		i = max - 1
	}

	p.view.Highlight(createPostTag(i))
}

// GetPostsCount : ポスト数を取得
func (p *posts) GetPostsCount() int {
	c := len(p.contents)

	if l := len(p.pinned); l > 0 {
		c += l
	}

	return c
}

// createPostTag : ポスト追跡用のタグを作成
func createPostTag(id int) string {
	return fmt.Sprintf("post_%d", id)
}

// SetPinned : ピン留めを登録
func (p *posts) SetPinned(pinned []*sharedapi.Post) {
	p.pinned = []*sharedapi.Post{}

	for i, post := range pinned {
		p.pinned[i] = post
	}
}

// Update : ポストを更新
func (p *posts) Update(posts []*sharedapi.Post) {
	addedPostsCount := p.addPosts(posts)
	cursorPos := p.getCurrentCursorPos()

	// 先頭以外のツイートを選択中の場合、更新後もそのツイートを選択したままにする
	// NOTE: "先頭以外" なのは、ストリームモードで放置した時にカーソルが段々下に下がってしまうのを防ぐため
	if cursorPos != 0 {
		cursorPos += addedPostsCount
	}

	// t.draw(cursorPos)
}

// getCurrentCursorPos : 現在のカーソル位置を取得
func (p *posts) getCurrentCursorPos() int {
	pos := getHighlightId(p.view.GetHighlights())

	if pos == -1 {
		pos = 0
	}

	return pos
}

// addPosts : ポストを追加
func (p *posts) addPosts(posts []*sharedapi.Post) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	size := len(p.contents)
	addSize := len(posts)
	allSize := size + addSize
	maxSize := global.conf.Pref.Feature.AccmulateTweetsLimit

	// 最大蓄積数を超えていたら古いものから削除
	if allSize > maxSize {
		size -= allSize - maxSize
	}

	p.contents = append(posts, p.contents[:size]...)

	return addSize
}

// DeletePost : ポストを削除
func (p *posts) DeletePost(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, ok := find(p.contents, func(c *sharedapi.Post) bool {
		// リポスト元のIDを参照
		if ref := c.Reference; ref != nil {
			return ref.ID == id
		}

		return c.ID == id
	})

	if !ok {
		return
	}

	// i番目の要素を削除
	p.contents = p.contents[:i+copy(p.contents[i:], p.contents[i+1:])]

	// 再描画して反映
	// t.draw(t.getCurrentCursorPos())
}

// draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (p *posts) draw(cursorPos int) {
	// icon := global.conf.Pref.Icon
	appearance := global.conf.Pref.Appearance
	// width := getWindowWidth()

	p.view.
		SetTextAlign(tview.AlignLeft).
		Clear()

	// 表示するポストが無いなら描画を中断
	if p.GetPostsCount() == 0 {
		p.DrawMessage(global.conf.Pref.Text.NoPosts)
		return
	}

	contents := p.contents

	// ピン留めがある場合、先頭に追加
	if len(p.pinned) > 0 {
		contents = append(p.pinned, p.contents...)
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

		// fmt.Fprintln(t.view, createTweetLayout(annotation, content, i, width))
		fmt.Fprintf(p.view, "%d : %s\n", i, post.Text)

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
		// if i < t.GetTweetsCount()-1 {
		// 	fmt.Fprintln(t.view, createSeparator(appearance.TweetSeparator, width))
		// }
	}

	p.scrollToPost(cursorPos)
}

// DrawMessage : ビューにメッセージを表示
func (p *posts) DrawMessage(s string) {
	p.view.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}
