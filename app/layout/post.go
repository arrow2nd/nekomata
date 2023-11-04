package layout

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

// CreatePostHighlightTag : ポストのハイライトタグを作成
func CreatePostHighlightTag(id int) string {
	return fmt.Sprintf("post_%d", id)
}

// CreatePostSeparator : セパレータを作成
func (l *Layout) CreatePostSeparator(sep string, w int) string {
	return CreateStyledText(l.Style.Tweet.Separator, strings.Repeat(sep, l.Width), "")
}

// CreatePost : 投稿のレイアウトを作成
func (l *Layout) CreatePost(i int, p *sharedapi.Post) (string, error) {
	layout := ""

	// ピン止めツイート
	// TODO: 後で対応する
	// if i == 0 && t.pinned != nil {
	// 	annotation += fmt.Sprintf("[gray:-:-]%s Pinned Tweet[-:-:-]", icon.Pinned)
	// }

	// リポストなら元ポストに置き換える
	if p.Reference != nil {
		layout += l.createAnnotation("Reposted by", p.Author.DisplayName, "@"+p.Author.Username)
		p = p.Reference
	}

	// リプライ先のアノテーション
	if len(p.Mentions) > 0 {
		users := []string{}
		for _, u := range p.Mentions {
			users = append(users, u.DisplayName+" @"+u.Username)
		}
		layout += l.createAnnotation("Reply to", strings.Join(users, ", "))
	}

	funcMap := template.FuncMap{
		"author": func() (string, error) {
			return l.createUser(i, p.Author)
		},
		"text": func() string {
			return l.createPostText(p)
		},
		"detail": func() (string, error) {
			return l.createPostDetail(p)
		},
		"metrics": func() (string, error) {
			metrics := []string{}

			// ブックマーク済み
			if p.Bookmarked {
				metrics = append(metrics, CreateStyledText(l.Style.Tweet.Bookmarked, l.Text.Bookmarked, ""))
			}

			// リポスト数
			if p.RepostCount > 0 {
				text := l.createPostMetrics(&sharedapi.Reaction{
					Name:    l.Text.Repost,
					Count:   p.RepostCount,
					Reacted: p.Reposted,
				}, l.Style.Tweet.Repost, l.Style.Tweet.Reposted)

				metrics = append(metrics, text)
			}

			// リアクション
			for _, r := range p.Reactions {
				if text := l.createPostMetrics(&r, l.Style.Tweet.Like, l.Style.Tweet.Liked); text != "" {
					metrics = append(metrics, text)
				}
			}

			return strings.Join(metrics, " "), nil
		},
	}

	t, err := template.New("post").Funcs(funcMap).Parse(l.Template.Post)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, *p); err != nil {
		return "", err
	}

	// 不要な改行を除く
	return strings.TrimSpace(buf.String()), nil
}

func (l *Layout) createAnnotation(t ...string) string {
	text := strings.Join(t, " ")
	return CreateStyledText(l.Style.Tweet.Annotation, text, "") + "\n"
}

func (l *Layout) createPostText(p *sharedapi.Post) string {
	text := p.Text

	// メンションをハイライト
	styledMention := CreateStyledText(l.Style.Tweet.Mention, "$1@$2", "")
	text = regexp.MustCompile(`(^|[^\w@#$%&/])@(\w+)`).ReplaceAllString(text, styledMention)

	// ハッシュタグをハイライト
	for _, tag := range p.Tags {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)[#＃](%s\s|%s$)`, tag.Name, tag.Name))
		styledHashtag := CreateStyledText(l.Style.Tweet.HashTag, "#$1", tag.URL)
		text = re.ReplaceAllString(text, styledHashtag)
	}

	return text
}

func (l *Layout) createPostDetail(p *sharedapi.Post) (string, error) {
	funcMap := template.FuncMap{
		"createdAt": func() string {
			return l.convertDateString(p.CreatedAt)
		},
	}

	t, err := template.New("PostDetail").Funcs(funcMap).Parse(l.Template.PostDetail)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, p); err != nil {
		return "", err
	}

	return CreateStyledText(l.Style.Tweet.Detail, buf.String(), ""), nil
}

func (l *Layout) createPostMetrics(r *sharedapi.Reaction, normalStyle, revStyle string) string {
	if r.Count == 0 {
		return ""
	}

	style := normalStyle
	if r.Reacted {
		style = revStyle
	}

	text := fmt.Sprintf("%s %d", r.Name, r.Count)
	return CreateStyledText(style, text, "")
}
