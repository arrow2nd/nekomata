package config

const PreferencesVersion = 1

// Feature : 機能
type Feature struct {
	// MainUser : メインで使用するユーザー
	MainUser string `toml:"main_user"`
	// LoadPostCount : 1度に読み込む投稿数
	LoadPostCount int `toml:"load_post_count"`
	// MaxPostCount : 投稿の最大蓄積数
	MaxPostCount int `toml:"max_post_count"`
	// UseExternalEditor : 投稿内容に編集に外部エディタを使用するか
	UseExternalEditor bool `toml:"use_external_editor"`
	// IsLocaleCJK : ロケールがCJKか
	IsLocaleCJK bool `toml:"is_locale_cjk"`
	// StartupCmds : 起動時に実行するコマンド
	StartupCmds []string `toml:"startup_cmds"`
}

// Appearance : 外観
type Appearance struct {
	// StyleFilePath : 配色テーマファイルのパス
	StyleFilePath string `toml:"style_file"`
	// DateFormat : 日付のフォーマット
	DateFormat string `toml:"date_fmt"`
	// TimeFormat : 時刻のフォーマット
	TimeFormat string `toml:"time_fmt"`
	// UserBIOMaxRow : ユーザBIOの最大表示行数
	UserBIOMaxRow int `toml:"user_bio_max_row"`
	// UserProfilePaddingX : ユーザプロフィールの左右パディング
	UserProfilePaddingX int `toml:"user_profile_padding_x"`
	// UserDetailSeparator : ユーザ詳細のセパレータ
	UserDetailSeparator string `toml:"user_detail_separator"`
	// HidePostSeparator : 投稿間のセパレータを非表示
	HidePostSeparator bool `toml:"hide_post_separator"`
	// PostSeparator : 投稿のセパレータ
	PostSeparator string `toml:"post_separator"`
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
}

// Text : 表示テキスト
type Text struct {
	// Bookmarked : ブックマーク済み
	Bookmarked string `toml:"bookmarked"`
	// Repost : リポストの単位
	Repost string `toml:"repost"`
	// Loading : 読み込み中
	Loading string `toml:"loading"`
	// NoPosts : ポスト無し
	NoPosts string `toml:"no_posts"`
	// PostTextAreaHint : テキストエリアのヒント
	PostTextAreaHint string `toml:"post_textarea_hint"`
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
	// TabBookmark : ブックマークタブ
	TabBookmark string `toml:"tab_bookmark"`
	// TabUser : ユーザタブ
	TabUser string `toml:"tab_user"`
	// TabSearch : 検索タブ
	TabSearch string `toml:"tab_search"`
	// TabLikes : いいねリストタブ
	TabLikes string `toml:"tab_likes"`
	// TabAnnouncement : アナウンスタブ
	TabAnnouncement string `toml:"announcement_home"`
	// TabDocs : ドキュメントタブ
	TabDocs string `toml:"tab_docs"`
}

// Icon : アイコン
type Icon struct {
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
}

// Preferences : 環境設定
type Preferences struct {
	Version     int             `toml:"version"`
	Feature     Feature         `toml:"feature"`
	Confirm     map[string]bool `toml:"comfirm"`
	Appearance  Appearance      `toml:"appearance"`
	Template    Template        `toml:"template"`
	Text        Text            `toml:"text"`
	Icon        Icon            `toml:"icon"`
	Keybindings Keybindings     `toml:"keybinding"`
}

// defaultPreferences : デフォルト設定
func defaultPreferences() *Preferences {
	return &Preferences{
		Version: PreferencesVersion,
		Feature: Feature{
			MainUser:          "",
			LoadPostCount:     25,
			MaxPostCount:      250,
			UseExternalEditor: false,
			IsLocaleCJK:       true,
			StartupCmds: []string{
				"home",
			},
		},
		Confirm: map[string]bool{
			"reaction":        true,
			"remove reaction": true,
			"repost":          true,
			"remove repost":   true,
			"delete":          true,
			"follow":          true,
			"unfollow":        true,
			"block":           true,
			"unblock":         true,
			"mute":            true,
			"unmute":          true,
			"post":            true,
			"quit":            true,
		},
		Appearance: Appearance{
			StyleFilePath:       "style_default.toml",
			DateFormat:          "2006/01/02",
			TimeFormat:          "15:04:05",
			UserBIOMaxRow:       3,
			UserProfilePaddingX: 4,
			UserDetailSeparator: " | ",
			HidePostSeparator:   false,
			PostSeparator:       "─",
			GraphChar:           "\u2588",
			GraphMaxWidth:       30,
			TabSeparator:        "|",
			TabMaxWidth:         20,
		},
		Template: Template{
			Post:           "{{ author }}\n{{ text }}\n{{ detail }}\n{{ metrics }}",
			PostAnnotation: "{text} {author_name} {author_username}",
			PostDetail:     "{{ createdAt }}{{ if .Via }} | via {{ .Via }}{{ end }}",
			PostPoll:       "",
			PostPollGraph:  "",
			PostPollDetail: "",
			User:           "{{ displayName }} {{ username }} {{ badges }}",
			UserDetail:     "",
		},
		Text: Text{
			Bookmarked:       "Bookmarked",
			Repost:           "Repost",
			Loading:          "Loading...",
			NoPosts:          "No posts ฅ^-ω-^ฅ",
			PostTextAreaHint: "Meow",
			TabHome:          "Home",
			TabGlobal:        "Global",
			TabLocal:         "Local",
			TabList:          "List: {name}",
			TabMention:       "Mention",
			TabBookmark:      "Bookmark",
			TabUser:          "User: @{name}",
			TabSearch:        "Search: {query}",
			TabLikes:         "Likes: @{name}",
			TabAnnouncement:  "Announcement",
			TabDocs:          "Docs: {name}",
		},
		Icon: Icon{
			Geo:      "📍",
			Link:     "🔗",
			Pinned:   "📌",
			Verified: "✅",
			Private:  "🔒",
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
				ActionReloadPage: {"."},
			},
			Posts: keybinding{
				ActionScrollUp:           {"ctrl+j", "PageUp"},
				ActionScrollDown:         {"ctrl+k", "PageDown"},
				ActionCursorUp:           {"k", "Up"},
				ActionCursorDown:         {"j", "Down"},
				ActionCursorTop:          {"g", "Home"},
				ActionCursorBottom:       {"G", "End"},
				ActionPostReaction:       {"f"},
				ActionPostRemoveReaction: {"F"},
				ActionPostRepost:         {"t"},
				ActionPostRemoveRepost:   {"T"},
				ActionPostDelete:         {"D"},
				ActionUserFollow:         {"w"},
				ActionUserUnfollow:       {"W"},
				ActionUserBlock:          {"x"},
				ActionUserUnblock:        {"X"},
				ActionUserMute:           {"u"},
				ActionUserUnmute:         {"U"},
				ActionOpenUserPage:       {"i"},
				ActionOpenUserLikes:      {"I"},
				ActionPost:               {"n"},
				ActionReply:              {"r"},
				ActionOpenBrowser:        {"o"},
				ActionCopyUrl:            {"c"},
			},
		},
	}
}
