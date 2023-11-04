package config

import "github.com/gdamore/tcell/v2"

type color string

// ToColor : tcell.Colorに変換
func (c color) ToColor() tcell.Color {
	return tcell.GetColor(string(c))
}

// AppStyle : アプリ全体のスタイル
type AppStyle struct {
	BackgroundColor color  `toml:"background_color"`
	BorderColor     color  `toml:"border_color"`
	TextColor       color  `toml:"text_color"`
	SubTextColor    color  `toml:"sub_text_color"`
	EmphasisText    string `toml:"emphasis_text"`
}

// TabStyle : タブのスタイル
type TabStyle struct {
	Text            string `toml:"text"`
	BackgroundColor color  `toml:"background_color"`
}

// AutocompleteStyle : 補完候補のスタイル
type AutocompleteStyle struct {
	TextColor               color `toml:"text_color"`
	BackgroundColor         color `toml:"background_color"`
	SelectedBackgroundColor color `toml:"selected_background_color"`
}

// StatusBarStyle : ステータスバーのスタイル
type StatusBarStyle struct {
	Text            string `toml:"text"`
	BackgroundColor color  `toml:"background_color"`
}

// PostStyle : 投稿のスタイル
type PostStyle struct {
	Annotation string `toml:"annotation"`
	Detail     string `toml:"detail"`
	Like       string `toml:"like"`
	Liked      string `toml:"liked"`
	Repost     string `toml:"repost"`
	Reposted   string `toml:"reposted"`
	Bookmarked string `toml:"bookmarked"`
	HashTag    string `toml:"hashtag"`
	Mention    string `toml:"mention"`
	PollGraph  string `toml:"poll_graph"`
	PollDetail string `toml:"poll_detail"`
	Separator  string `toml:"separator"`
}

// UserStyle : ユーザのスタイル
type UserStyle struct {
	Name     string `toml:"name"`
	UserName string `toml:"user_name"`
	Detail   string `toml:"detail"`
	Verified string `toml:"verified"`
	Private  string `toml:"private"`
}

// MetricsStyle : ユーザメトリクスのスタイル
type MetricsStyle struct {
	PostsText                string `toml:"posts_text"`
	PostsBackgroundColor     color  `toml:"posts_background_color"`
	FollowingText            string `toml:"following_text"`
	FollowingBackgroundColor color  `toml:"following_background_color"`
	FollowersText            string `toml:"followers_text"`
	FollowersBackgroundColor color  `toml:"followers_background_color"`
}

// Style : スタイル定義
type Style struct {
	App          AppStyle          `toml:"app"`
	Tab          TabStyle          `toml:"tab"`
	Autocomplate AutocompleteStyle `toml:"autocomplete"`
	StatusBar    StatusBarStyle    `toml:"statusbar"`
	Post         PostStyle         `toml:"post"`
	User         UserStyle         `toml:"user"`
	Metrics      MetricsStyle      `toml:"metrics"`
}

func defaultStyle() *Style {
	return &Style{
		App: AppStyle{
			BackgroundColor: "#000000",
			BorderColor:     "#ffffff",
			TextColor:       "#f9f9f9",
			SubTextColor:    "#979797",
			EmphasisText:    "maroon:-:bi",
		},
		Tab: TabStyle{
			Text:            "white:-:-",
			BackgroundColor: "#000000",
		},
		Autocomplate: AutocompleteStyle{
			TextColor:               "#000000",
			BackgroundColor:         "#808080",
			SelectedBackgroundColor: "#C0C0C0",
		},
		StatusBar: StatusBarStyle{
			Text:            "black:-:-",
			BackgroundColor: "#ffffff",
		},
		Post: PostStyle{
			Annotation: "teal:-:-",
			Detail:     "gray:-:-",
			Like:       "pink:-:-",
			Liked:      "pink:-:r",
			Repost:     "lime:-:-",
			Reposted:   "lime:-:r",
			Bookmarked: "red:-:r",
			HashTag:    "aqua:-:-",
			Mention:    "aqua:-:-",
			PollGraph:  "aqua:-:-",
			PollDetail: "gray:-:-",
			Separator:  "gray:-:-",
		},
		User: UserStyle{
			Name:     "white:-:b",
			UserName: "gray:-:i",
			Detail:   "gray:-:-",
			Verified: "blue:-:-",
			Private:  "gray:-:-",
		},
		Metrics: MetricsStyle{
			PostsText:                "black:-:-",
			PostsBackgroundColor:     "#a094c7",
			FollowingText:            "black:-:-",
			FollowingBackgroundColor: "#84a0c6",
			FollowersText:            "black:-:-",
			FollowersBackgroundColor: "#89b8c2",
		},
	}
}
