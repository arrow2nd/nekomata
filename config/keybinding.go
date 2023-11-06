package config

import (
	"fmt"
	"strings"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/gdamore/tcell/v2"
)

const (
	// アプリ全体のアクション
	ActionQuit = "quit"

	// ページビューのアクション
	ActionSelectPrevTab = "select_prev_tab"
	ActionSelectNextTab = "select_next_tab"
	ActionClosePage     = "close_page"
	ActionRedraw        = "redraw"
	ActionFocusCmdLine  = "focus_cmdline"
	ActionShowHelp      = "show_help"

	// ページ共通のアクション
	ActionReloadPage = "reload_page"

	// 投稿一覧のアクション
	ActionScrollUp     = "scroll_up"
	ActionScrollDown   = "scroll_down"
	ActionCursorUp     = "cursor_up"
	ActionCursorDown   = "cursor_down"
	ActionCursorTop    = "cursor_top"
	ActionCursorBottom = "cursor_bottom"
	ActionPost         = "post"
	ActionReaction     = "reaction"
	ActionUnreaction   = "unreaction"
	ActionRepost       = "repost"
	ActionUnrepost     = "unrepost"
	ActionDelete       = "delete"
	ActionBookmark     = "bookmark"
	ActionUnbookmark   = "unbookmark"
	ActionReply        = "reply"
	ActionOpenBrowser  = "open_browser"
	ActionOpenUserPage = "open_user"
	ActionCopyUrl      = "copy_url"

	// ユーザーページのアクション
	ActionFollow        = "follow"
	ActionUnfollow      = "unfollow"
	ActionBlock         = "block"
	ActionUnblock       = "unblock"
	ActionMute          = "mute"
	ActionUnmute        = "unmute"
	ActionOpenReactions = "open_reactions"
)

type keybinding map[string][]string

// GetString : キーバインド文字列を取得
func (k keybinding) GetString(key string) string {
	s := strings.Join(k[key], ", ")

	if s == "" {
		return "*No assignment*"
	}

	return s
}

// MappingEventHandler : キーバインドにイベントハンドラをマッピング
func (k keybinding) MappingEventHandler(handlers map[string]func()) (*cbind.Configuration, error) {
	c := cbind.NewConfiguration()

	for action, keys := range k {
		f, ok := handlers[action]
		if !ok {
			return nil, fmt.Errorf("unknown action: %s", action)
		}

		handler := func(_ *tcell.EventKey) *tcell.EventKey {
			f()
			return nil
		}

		for _, key := range keys {
			key = strings.TrimSpace(key)

			if key == "" {
				continue
			}

			if err := c.Set(key, handler); err != nil {
				return nil, fmt.Errorf("failed to set key bindings: %w", err)
			}
		}

	}

	return c, nil
}
