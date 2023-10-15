package config

const PreferencesVersion = 1

// Feature : 機能
type Feature struct {
	// MainUser : メインで使用するユーザー
	MainUser string `toml:"main_user"`
	// LoadTweetsLimit : 1度に読み込むツイート数
	LoadTweetsLimit int `toml:"load_tweets_limit"`
	// AccmulateTweetsLimit : ツイートの最大蓄積数
	AccmulateTweetsLimit int `toml:"accmulate_tweets_limit"`
	// UseExternalEditor : ツイート編集に外部エディタを使用するか
	UseExternalEditor bool `toml:"use_external_editor"`
	// IsLocaleCJK : ロケールがCJKか
	IsLocaleCJK bool `toml:"is_locale_cjk"`
	// StartupCmds : 起動時に実行するコマンド
	StartupCmds []string `toml:"startup_cmds"`
}

// Appearancene : 外観
type Appearancene struct {
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
	// HideTweetSeparator : ツイート間のセパレータを非表示
	HideTweetSeparator bool `toml:"hide_tweet_separator"`
	// HideQuoteTweetSeparator : 引用ツイートのセパレータを非表示
	HideQuoteTweetSeparator bool `toml:"hide_quote_tweet_separator"`
	// TweetSeparator : ツイートのセパレータ
	TweetSeparator string `toml:"tweet_separator"`
	// QuoteTweetSeparator : 引用ツイートのセパレータ
	QuoteTweetSeparator string `toml:"quote_tweet_separator"`
	// GraphChar : 投票グラフの表示に使用する文字
	GraphChar string `toml:"graph_char"`
	// GraphMaxWidth : 投票グラフの最大表示幅
	GraphMaxWidth int `toml:"graph_max_width"`
	// TabSeparator : タブのセパレータ
	TabSeparator string `toml:"tab_separator"`
	// TabMaxWidth : タブの最大表示幅
	TabMaxWidth int `toml:"tab_max_width"`
}

// Layout : 表示レイアウト
type Layout struct {
	// Tweet : ツイート
	Tweet string `toml:"tweet"`
	// TweetAnotation : ツイートアノテーション
	TweetAnotation string `toml:"tweet_anotation"`
	// TweetDetail : ツイート詳細
	TweetDetail string `toml:"tweet_detail"`
	// TweetPoll : 投票
	TweetPoll string `toml:"tweet_poll"`
	// TweetPollGraph : 投票グラフ
	TweetPollGraph string `toml:"tweet_poll_graph"`
	// TweetPollDetail : 投票詳細
	TweetPollDetail string `toml:"tweet_poll_detail"`
	// User : ユーザプロフィール
	User string `toml:"user"`
	// UserInfo : ユーザ情報
	UserInfo string `toml:"user_info"`
}

// Text : 表示テキスト
type Text struct {
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
	// TweetView : ツイートビューのキーバインド
	TweetView keybinding `toml:"tweet"`
}

// Preferences : 環境設定
type Preferences struct {
	Version     int             `toml:"version"`
	Feature     Feature         `toml:"feature"`
	Confirm     map[string]bool `toml:"comfirm"`
	Appearance  Appearancene    `toml:"appearance"`
	Layout      Layout          `toml:"layout"`
	Text        Text            `toml:"text"`
	Icon        Icon            `toml:"icon"`
	Keybindings Keybindings     `toml:"keybinding"`
}

// defaultPreferences : デフォルト設定
func defaultPreferences() *Preferences {
	return &Preferences{
		Version: PreferencesVersion,
		Feature: Feature{
			MainUser:             "",
			LoadTweetsLimit:      25,
			AccmulateTweetsLimit: 250,
			UseExternalEditor:    false,
			IsLocaleCJK:          true,
			StartupCmds: []string{
				"home",
			},
		},
		Confirm: map[string]bool{
			"like":      true,
			"unlike":    true,
			"retweet":   true,
			"unretweet": true,
			"delete":    true,
			"follow":    true,
			"unfollow":  true,
			"block":     true,
			"unblock":   true,
			"mute":      true,
			"unmute":    true,
			"tweet":     true,
			"quit":      true,
		},
		Appearance: Appearancene{
			StyleFilePath:           "style_default.toml",
			DateFormat:              "2006/01/02",
			TimeFormat:              "15:04:05",
			UserBIOMaxRow:           3,
			UserProfilePaddingX:     4,
			UserDetailSeparator:     " | ",
			HideTweetSeparator:      false,
			HideQuoteTweetSeparator: false,
			TweetSeparator:          "─",
			QuoteTweetSeparator:     "-",
			GraphChar:               "\u2588",
			GraphMaxWidth:           30,
			TabSeparator:            "|",
			TabMaxWidth:             20,
		},
		Layout: Layout{
			Tweet:           "{annotation}\n{user_info}\n{text}\n{poll}\n{detail}",
			TweetAnotation:  "{text} {author_name} {author_username}",
			TweetDetail:     "{created_at} | via {via}\n{metrics}",
			TweetPoll:       "{graph}\n{detail}",
			TweetPollGraph:  "{label}\n{graph} {per} {votes}",
			TweetPollDetail: "{status} | {all_votes} votes | ends on {end_date}",
			User:            "{user_info}\n{bio}\n{user_detail}",
			UserInfo:        "{name} {username} {badge}",
		},
		Text: Text{
			Repost:           "RP",
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
			TweetView: keybinding{
				ActionScrollUp:       {"ctrl+j", "PageUp"},
				ActionScrollDown:     {"ctrl+k", "PageDown"},
				ActionCursorUp:       {"k", "Up"},
				ActionCursorDown:     {"j", "Down"},
				ActionCursorTop:      {"g", "Home"},
				ActionCursorBottom:   {"G", "End"},
				ActionTweetLike:      {"f"},
				ActionTweetUnlike:    {"F"},
				ActionTweetRetweet:   {"t"},
				ActionTweetUnretweet: {"T"},
				ActionTweetDelete:    {"D"},
				ActionUserFollow:     {"w"},
				ActionUserUnfollow:   {"W"},
				ActionUserBlock:      {"x"},
				ActionUserUnblock:    {"X"},
				ActionUserMute:       {"u"},
				ActionUserUnmute:     {"U"},
				ActionOpenUserPage:   {"i"},
				ActionOpenUserLikes:  {"I"},
				ActionTweet:          {"n"},
				ActionQuote:          {"q"},
				ActionReply:          {"r"},
				ActionOpenBrowser:    {"o"},
				ActionCopyUrl:        {"c"},
			},
		},
	}
}
