package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/config"
	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

// actionPost : 投稿に対してアクション
func (p *postList) actionPost(action string) {
	target := p.getSelectPost()
	if target == nil {
		return
	}

	id := target.ID

	f := func() {
		var result *sharedapi.Post = nil
		err := fmt.Errorf("unknown action: %s", action)

		switch action {
		case config.ActionReaction:
			// TODO: リアクション種別が複数ある場合どうにかする (Misskey)
			result, err = global.client.Reaction(id, "")
		case config.ActionUnreaction:
			result, err = global.client.Unreaction(id)
		case config.ActionRepost:
			result, err = global.client.Repost(id)
		case config.ActionUnrepost:
			result, err = global.client.Unrepost(id)
		case config.ActionBookmark:
			result, err = global.client.Bookmark(id)
		case config.ActionUnbookmark:
			result, err = global.client.Unbookmark(id)
		case config.ActionDelete:
			err = global.client.DeletePost(id)
		}

		if err != nil {
			global.SetErrorStatus(action, err.Error())
			return
		}

		// ステータスを反映
		switch action {
		case config.ActionReaction:
			// TODO: リアクション種別が複数ある場合どうにかする (Misskey)
			synchronizeResponseCounts(target.Reactions[0].Count, &result.Reactions[0].Count, 1)
			target.Reactions = result.Reactions
		case config.ActionUnreaction:
			// TODO: リアクション種別が複数ある場合どうにかする (Misskey)
			synchronizeResponseCounts(target.Reactions[0].Count, &result.Reactions[0].Count, -1)
			target.Reactions = result.Reactions
		case config.ActionRepost:
			synchronizeResponseCounts(target.RepostCount, &result.RepostCount, 1)
			target.Reposted = result.Reposted
			target.RepostCount = result.RepostCount
		case config.ActionUnrepost:
			synchronizeResponseCounts(target.RepostCount, &result.RepostCount, -1)
			target.Reposted = result.Reposted
			target.RepostCount = result.RepostCount
		case config.ActionBookmark, config.ActionUnbookmark:
			target.Bookmarked = result.Bookmarked
		case config.ActionDelete:
			if err := p.DeletePost(id); err != nil {
				global.SetErrorStatus("Delete", err.Error())
				return
			}
		}

		// 再描画
		p.draw(p.getCurrentCursorPos())

		if !strings.HasSuffix(action, "e") {
			action += "e"
		}

		global.SetStatus(action+"d", createPostSummary(target))
	}

	// 確認画面が不要ならそのまま実行
	if !global.conf.Pref.Confirm[strings.ToLower(action)] {
		f()
		return
	}

	title := fmt.Sprintf(
		"Do you want to [%s]%s[-:-:-] this post?",
		global.conf.Style.App.EmphasisText,
		strings.ToLower(action),
	)

	global.ReqestPopupModal(&ModalOpts{title, "", f})
}

// synchronizeResponseCounts : レスポンスの数値と実際の状態とを同期させる
func synchronizeResponseCounts(prev int, next *int, add int) {
	if prev == *next {
		*next = *next + add
	}
}

// createPostSummary : 投稿の要約を作成
func createPostSummary(p *sharedapi.Post) string {
	text := p.Text
	if text == "" {
		text = "<empty>"
	}
	return fmt.Sprintf("%s | %s", createUserSummary(p.Author), text)
}

// openBrowser : 投稿をブラウザで表示
func (p *postList) openBrowser() {
	post := p.getSelectPost()
	if post == nil {
		return
	}

	u, err := global.client.CreatePostURL(post)
	if err != nil {
		global.SetErrorStatus("Open", err.Error())
		return
	}

	browser.Stdout = io.Discard
	browser.Stderr = io.Discard

	if err := browser.OpenURL(u); err != nil {
		global.SetErrorStatus("Open", err.Error())
		return
	}

	global.SetStatus("Opened", createPostSummary(post))
}

// copyToClipboard : 投稿のURLをクリップボードにコピー
func (p *postList) copyToClipboard() {
	post := p.getSelectPost()
	if post == nil {
		return
	}

	u, err := global.client.CreatePostURL(post)
	if err != nil {
		global.SetErrorStatus("Copy", err.Error())
		return
	}

	if err := clipboard.WriteAll(u); err != nil {
		global.SetErrorStatus("Copy", err.Error())
		return
	}

	global.SetStatus("Copied", createPostSummary(post))
}
