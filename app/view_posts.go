package app

import (
	"fmt"
	"sync"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/app/layout"
	"github.com/arrow2nd/nekomata/config"
	"github.com/rivo/tview"
	"golang.org/x/net/context"
)

// cursorMove : カーソルの移動量
const (
	cursorMoveUp   int = -1
	cursorMoveDown int = 1
)

type postList struct {
	textView     *tview.TextView
	pinnedPosts  []*sharedapi.Post
	posts        []*sharedapi.Post
	mu           sync.Mutex
	layout       *layout.Layout
	streamCancel context.CancelFunc
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
			p.highlightCursor(0)
			p.textView.ScrollToHighlight()
		},
		config.ActionCursorBottom: func() {
			lastIndex := p.GetPostsCount() - 1
			p.highlightCursor(lastIndex)
			p.textView.ScrollToHighlight()
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

	p.highlightCursor(idx + int(c))
	p.textView.ScrollToHighlight()
}

// highlightCursor : カーソルをハイライト
func (p *postList) highlightCursor(i int) {
	// 範囲内に丸める
	if max := p.GetPostsCount(); i < 0 {
		i = 0
	} else if i >= max {
		i = max - 1
	}

	p.textView.Highlight(layout.CreatePostHighlightTag(i))
}

// GetPostsCount : 投稿数を取得
func (p *postList) GetPostsCount() int {
	c := len(p.posts)

	if l := len(p.pinnedPosts); l > 0 {
		c += l
	}

	return c
}

// GetSinceId : 最新の投稿のIDを取得
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

// Update : 投稿一覧を更新
func (p *postList) Update(posts []*sharedapi.Post) error {
	addedPostsCount := p.addPosts(posts)
	cursorPos := p.getCurrentCursorPos()

	// 先頭ではない投稿を選択中の場合、更新後もその投稿を選択したままにする
	if cursorPos != 0 {
		cursorPos += addedPostsCount
	}

	var err error = nil

	global.app.QueueUpdateDraw(func() {
		err = p.draw(cursorPos)
	})

	return err
}

// getCurrentCursorPos : 現在のカーソル位置を取得
func (p *postList) getCurrentCursorPos() int {
	pos := getHighlightId(p.textView.GetHighlights())

	if pos == -1 {
		pos = 0
	}

	return pos
}

// addPosts : 投稿を追加
func (p *postList) addPosts(posts []*sharedapi.Post) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	size := len(p.posts)
	addSize := len(posts)
	allSize := size + addSize
	maxSize := global.conf.Pref.Feature.MaxPostCount

	// 最大蓄積数を超えていたら古いものから削除
	if allSize > maxSize {
		size -= allSize - maxSize
	}

	p.posts = append(posts, p.posts[:size]...)

	return addSize
}

// DeletePost : 投稿を削除
func (p *postList) DeletePost(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, ok := find(p.posts, func(c *sharedapi.Post) bool {
		// 引用元のIDを参照
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

	// カーソルが流れる可能性があるか
	scrollOffsetRow, _ := p.textView.GetScrollOffset()
	isCursorFlowing := cursorPos != 0 && scrollOffsetRow == 0

	// カーソル行数の計算用
	isCalculatingLines := isCursorFlowing
	cursorLineNum := 0

	p.textView.Clear()

	// 表示する投稿が無いなら描画を中断
	if p.GetPostsCount() == 0 {
		p.DrawMessage(global.conf.Pref.Text.NoPosts)
		return nil
	}

	width := getWindowWidth()
	contents := p.posts

	// ピン留めがある場合、先頭に追加
	if len(p.pinnedPosts) > 0 {
		contents = append(p.pinnedPosts, p.posts...)
	}

	for i, post := range contents {
		postLayout, err := p.layout.CreatePost(i, post)
		if err != nil {
			return err
		}

		fmt.Fprintln(p.textView, postLayout)

		// 末尾の投稿ではないならセパレータを挿入
		insertSeparator := !appearance.HideTweetSeparator || i < p.GetPostsCount()-1
		if insertSeparator {
			fmt.Fprintln(p.textView, p.layout.CreatePostSeparator(appearance.TweetSeparator, width))
		}

		// カーソルの行数を計算する必要がないならスキップ
		if !isCalculatingLines {
			continue
		}

		// カーソルの当たっている投稿なら計算を終了
		if i == cursorPos {
			cursorLineNum++
			isCalculatingLines = false
			continue
		}

		cursorLineNum += getStringDisplayRow(postLayout, width)
		if insertSeparator {
			cursorLineNum++
		}
	}

	p.highlightCursor(cursorPos)

	// カーソルが流れる & 位置がTextViewの半分より上 or
	// 既にスクロール済みならカーソル位置までスクロールさせる
	// NOTE: 無条件にScrollToHighlight()を呼ぶと画面の下半分が描画されないことがあるため
	_, _, _, innerHeight := p.textView.GetInnerRect()
	if (isCursorFlowing && cursorLineNum >= innerHeight/2) || (cursorPos != 0 && scrollOffsetRow != 0) {
		p.textView.ScrollToHighlight()
	}

	return nil
}

// DrawMessage : ビューにメッセージを表示
func (p *postList) DrawMessage(s string) {
	p.textView.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}
