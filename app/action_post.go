package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/config"
	"github.com/pkg/browser"
)

func (p *postList) action(action string) {
	target := p.getSelectPost()
	if target == nil {
		return
	}

	id := target.ID

	f := func() {
		var (
			err    error           = fmt.Errorf("unknown action: %s", action)
			result *sharedapi.Post = nil
		)

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
			result, err = global.client.DeletePost(id)
		}

		if err != nil {
			global.SetErrorStatus(action, err.Error())
			return
		}

		// ステータスを反映
		switch action {
		case config.ActionReaction:
			// TODO: リアクション種別が複数ある場合どうにかする (Misskey)
			if target.Reactions[0].Count == result.Reactions[0].Count {
				result.Reactions[0].Count++
			}
			target.Reactions = result.Reactions

		case config.ActionUnreaction:
			// TODO: リアクション種別が複数ある場合どうにかする (Misskey)
			if target.Reactions[0].Count == result.Reactions[0].Count {
				result.Reactions[0].Count--
			}
			target.Reactions = result.Reactions

		case config.ActionRepost:
			if target.RepostCount == result.RepostCount {
				result.RepostCount++
			}
			target.Reposted = result.Reposted
			target.RepostCount = result.RepostCount

		case config.ActionUnrepost:
			if target.RepostCount == result.RepostCount {
				result.RepostCount--
			}
			target.Reposted = result.Reposted
			target.RepostCount = result.RepostCount

		case config.ActionBookmark, config.ActionUnbookmark:
			target.Bookmarked = result.Bookmarked

		case config.ActionDelete:
			p.DeletePost(id)
		}

		// 再描画
		p.draw(p.getCurrentCursorPos())

		if !strings.HasSuffix(action, "e") {
			action += "e"
		}

		global.SetStatus(action+"d", createPostSummary(result))
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

func createPostSummary(p *sharedapi.Post) string {
	return fmt.Sprintf("%s | %s", createUserSummary(p.Author), p.Text)
}

func createUserSummary(u *sharedapi.Account) string {
	return fmt.Sprintf("%s @%s", u.DisplayName, u.Username)
}

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
