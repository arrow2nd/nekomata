package config

const PreferencesVersion = 1

// Feature : 機能
type Feature struct {
	// MainAccount : メインで使用するアカウント
	MainAccount string `toml:"main_account"`
	// LoadPostCount : 1度に読み込む投稿数
	LoadPostCount int `toml:"load_post_count"`
	// MaxPostCount : 投稿の最大蓄積数
	MaxPostCount int `toml:"max_post_count"`
	// IsLocaleCJK : ロケールがCJKか
	IsLocaleCJK bool `toml:"is_locale_cjk"`
	// StartupCmds : 起動時に実行するコマンド
	StartupCmds []string `toml:"startup_cmds"`
}

// Appearance : 外観
type Appearance struct {
	// StylePath : 配色テーマファイルのパス
	StylePath string `toml:"style_path"`
	// FormatDate : 日付のフォーマット
	FormatDate string `toml:"format_date"`
	// FormatTime : 時刻のフォーマット
	FormatTime string `toml:"format_time"`
	// UserBIOMaxRow : ユーザBIOの最大表示行数
	UserBIOMaxRow int `toml:"user_bio_max_row"`
	// UserDetailSeparator : ユーザ詳細のセパレータ
	UserDetailSeparator string `toml:"user_detail_separator"`
	// PostSeparator : 投稿のセパレータ
	PostSeparator string `toml:"post_separator"`
	// HidePostSeparator : 投稿間のセパレータを非表示
	HidePostSeparator bool `toml:"hide_post_separator"`
	// GraphChar : 投票グラフの表示に使用する文字
	GraphChar string `toml:"graph_char"`
	// GraphMaxWidth : 投票グラフの最大表示幅
	GraphMaxWidth int `toml:"graph_max_width"`
	// TabSeparator : タブのセパレータ
	TabSeparator string `toml:"tab_separator"`
	// TabMaxWidth : タブの最大表示幅
	TabMaxWidth int `toml:"tab_max_width"`
}

// Template : 表示テンプレート
type Template struct {
	// Post : 投稿
	Post string `toml:"post"`
	// PostAnnotation : 投稿のアノテーション
	PostAnnotation string `toml:"post_annotation"`
	// PostDetail : 投稿詳細
	PostDetail string `toml:"post_detail"`
	// PostPoll : 投票
	PostPoll string `toml:"post_poll"`
	// PostPollGraph : 投票グラフ
	PostPollGraph string `toml:"post_poll_graph"`
	// PostPollDetail : 投票詳細
	PostPollDetail string `toml:"post_poll_detail"`
	// User : ユーザ
	User string `toml:"user"`
	// UserDetail : ユーザ詳細
	UserDetail string `toml:"user_detail"`
	// TabHome : ホームタイムラインタブ
	TabHome string `toml:"tab_home"`
	// TabGlobal : グローバスタイムラインタブ
	TabGlobal string `toml:"tab_global"`
	// TabLocal : ローカルタイムラインタブ
	TabLocal string `toml:"tab_local"`
	// TabList : リストタブ
	TabList string `toml:"tab_list"`
	// TabMention : メンションタブ
	TabMention string `toml:"tab_mention"`
	// TabBookmarks : ブックマークタブ
	TabBookmarks string `toml:"tab_bookmark"`
	// TabUser : ユーザタブ
	TabUser string `toml:"tab_user"`
	// TabSearch : 検索タブ
	TabSearch string `toml:"tab_search"`
	// TabReactions : いいねタブ
	TabReactions string `toml:"tab_reactions"`
	// TabAnnouncements : アナウンスタブ
	TabAnnouncements string `toml:"announcements"`
	// TabDocument : ドキュメントタブ
	TabDocument string `toml:"tab_document"`
}

// Text : 表示テキスト
type Text struct {
	// Geo : 位置情報
	Geo string `toml:"geo"`
	// Link : リンク
	Link string `toml:"link"`
	// Pinned : ピン留め
	Pinned string `toml:"pinned"`
	// Verified : 認証バッジ
	Verified string `toml:"verified"`
	// Private : 非公開バッジ
	Private string `toml:"private"`
	// Bookmarked : ブックマーク済み
	Bookmarked string `toml:"bookmarked"`
	// Reaction : リアクション数
	Reaction string `toml:"reaction"`
	// Repost : リポスト数
	Repost string `toml:"repost"`
	// Loading : 読み込み中
	Loading string `toml:"loading"`
	// NoPosts : 投稿なし
	NoPosts string `toml:"no_posts"`
}

// Keybindings : キーバインド
type Keybindings struct {
	// Global : アプリ全体のキーバインド
	Global keybinding `toml:"global"`
	// View : メインビューのキーバインド
	View keybinding `toml:"view"`
	// Page : ページ共通のキーバインド
	Page keybinding `toml:"page"`
	// Posts : 投稿一覧のキーバインド
	Posts keybinding `toml:"posts"`
	// User : ユーザーページのキーバインド
	User keybinding `toml:"user"`
}

// Preferences : 環境設定
type Preferences struct {
	Version     int             `toml:"version"`
	Feature     Feature         `toml:"feature"`
	Confirm     map[string]bool `toml:"comfirm"`
	Appearance  Appearance      `toml:"appearance"`
	Template    Template        `toml:"template"`
	Text        Text            `toml:"text"`
	Keybindings Keybindings     `toml:"keybinding"`
}

// defaultPreferences : デフォルト設定
func defaultPreferences() *Preferences {
	return &Preferences{
		Version: PreferencesVersion,
		Feature: Feature{
			MainAccount:   "",
			LoadPostCount: 25,
			MaxPostCount:  250,
			IsLocaleCJK:   true,
			StartupCmds: []string{
				"timeline home",
			},
		},
		Confirm: map[string]bool{
			ActionReaction:   true,
			ActionUnreaction: true,
			ActionRepost:     true,
			ActionUnrepost:   true,
			ActionBookmark:   true,
			ActionUnbookmark: true,
			ActionDelete:     true,
			ActionFollow:     true,
			ActionUnfollow:   true,
			ActionBlock:      true,
			ActionUnblock:    true,
			ActionMute:       true,
			ActionUnmute:     true,
			ActionPost:       true,
			ActionQuit:       true,
		},
		Appearance: Appearance{
			StylePath:           "style_default.toml",
			FormatDate:          "2006/01/02",
			FormatTime:          "15:04:05",
			UserBIOMaxRow:       3,
			UserDetailSeparator: " | ",
			HidePostSeparator:   false,
			PostSeparator:       "─",
			GraphChar:           "\u2588",
			GraphMaxWidth:       30,
			TabSeparator:        "|",
			TabMaxWidth:         20,
		},
		Template: Template{
			Post:             "{{ author }}\n{{ text }}\n{{ detail }}\n{{ metrics }}",
			PostAnnotation:   "{text} {author_name} {author_username}",
			PostDetail:       "{{ createdAt }}{{ if .Via }} | via {{ .Via }}{{ end }}",
			PostPoll:         "",
			PostPollGraph:    "",
			PostPollDetail:   "",
			User:             "{{ displayName }} {{ username }} {{ badges }}",
			UserDetail:       "",
			TabHome:          "Home",
			TabGlobal:        "Global",
			TabLocal:         "Local",
			TabList:          "List: {{ name }}",
			TabMention:       "Mention",
			TabUser:          "User: @{{ name }}",
			TabSearch:        "Search: {{ name }}",
			TabBookmarks:     "Bookmarks",
			TabReactions:     "Reactions: @{{ name }}",
			TabAnnouncements: "Announcements",
			TabDocument:      "Document: {{ name }}",
		},
		Text: Text{
			Geo:        "📍",
			Link:       "🔗",
			Pinned:     "📌",
			Verified:   "✅",
			Private:    "🔒",
			Bookmarked: "Bookmarked",
			Reaction:   "Fav",
			Repost:     "Repost",
			Loading:    "Loading...",
			NoPosts:    "No posts ฅ^-ω-^ฅ",
		},
		Keybindings: Keybindings{
			Global: keybinding{
				ActionQuit: {"ctrl+q"},
			},
			View: keybinding{
				ActionSelectPrevTab: {"h", "Left"},
				ActionSelectNextTab: {"l", "Right"},
				ActionClosePage:     {"ctrl+w"},
				ActionRedraw:        {"ctrl+l"},
				ActionFocusCmdLine:  {":"},
				ActionShowHelp:      {"?"},
			},
			Page: keybinding{
				// ActionPost: {"n"},
				ActionReloadPage: {"."},
			},
			Posts: keybinding{
				ActionScrollUp:     {"ctrl+j", "PageUp"},
				ActionScrollDown:   {"ctrl+k", "PageDown"},
				ActionCursorUp:     {"k", "Up"},
				ActionCursorDown:   {"j", "Down"},
				ActionCursorTop:    {"g", "Home"},
				ActionCursorBottom: {"G", "End"},
				ActionReaction:     {"f"},
				ActionUnreaction:   {"F"},
				ActionRepost:       {"t"},
				ActionUnrepost:     {"T"},
				ActionBookmark:     {"b"},
				ActionUnbookmark:   {"B"},
				ActionDelete:       {"D"},
				ActionOpenUserPage: {"i"},
				ActionReply:        {"r"},
				ActionOpenBrowser:  {"o"},
				ActionCopyUrl:      {"c"},
			},
			User: keybinding{
				ActionFollow:        {"w"},
				ActionUnfollow:      {"W"},
				ActionBlock:         {"x"},
				ActionUnblock:       {"X"},
				ActionMute:          {"u"},
				ActionUnmute:        {"U"},
				ActionOpenReactions: {"I"},
			},
		},
	}
}
