package app

import (
	"fmt"
	"strings"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/config"
)

// actionUser : 投稿ユーザーに対してアクション
func (p *postList) actionUser(action string) {
	target := p.getSelectPost()
	if target == nil {
		return
	}

	author := target.Author

	f := func() {
		err := fmt.Errorf("unknown action: %s", action)

		switch action {
		case config.ActionFollow:
			_, err = global.client.Follow(author.ID)
		case config.ActionUnfollow:
			_, err = global.client.Unfollow(author.ID)
		case config.ActionBlock:
			_, err = global.client.Block(author.ID)
		case config.ActionUnblock:
			_, err = global.client.Unblock(author.ID)
		case config.ActionMute:
			_, err = global.client.Mute(author.ID)
		case config.ActionUnmute:
			_, err = global.client.Unmute(author.ID)
		}

		if err != nil {
			global.SetErrorStatus(action, err.Error())
			return
		}

		if !strings.HasSuffix(action, "e") {
			action += "e"
		}

		global.SetStatus(action+"d", createUserSummary(author))
	}

	// 確認画面が不要ならそのまま実行
	if !global.conf.Pref.Confirm[strings.ToLower(action)] {
		f()
		return
	}

	title := fmt.Sprintf(
		"Do you want to [%s]%s[-:-:-] this user?",
		global.conf.Style.App.EmphasisText,
		strings.ToLower(action),
	)

	global.ReqestPopupModal(&ModalOpts{title, "", f})
}

// createUserSummary : ユーザー情報の要約を作成
func createUserSummary(u *sharedapi.Account) string {
	return fmt.Sprintf("%s @%s", u.DisplayName, u.Username)
}
