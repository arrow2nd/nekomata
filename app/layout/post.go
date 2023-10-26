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
func CreatePostSeparator(sep string, w int) string {
	return ""
}

// Post : ポストのレイアウトを作成
func (l *Layout) Post(i int, p *sharedapi.Post) error {
	// リポストなら元ポストに置き換える
	if p.Reference != nil {
		l.printAnnotation("Reposted by", p.Author.DisplayName, "@"+p.Author.Username)
		p = p.Reference
	}

	// リプライ先のアノテーション
	if len(p.Mentions) > 0 {
		users := []string{}
		for _, u := range p.Mentions {
			users = append(users, u.DisplayName+" @"+u.Username)
		}
		l.printAnnotation("Reply to", strings.Join(users, ", "))
	}

	funcMap := template.FuncMap{
		"author": func() (string, error) {
			return l.createUserStr(i, p.Author)
		},
		"text": func() string {
			return l.createPostStr(p)
		},
		"detail": func() (string, error) {
			return l.createPostDetail(p)
		},
		"metrics": func() (string, error) {
			metrics := []string{}

			// リポスト数
			if p.RepostCount > 0 {
				text := l.createMetricsStr(&sharedapi.Reaction{
					Name:    l.Text.Repost,
					Count:   p.RepostCount,
					Reacted: p.Reposted,
				}, l.Style.Tweet.Repost, l.Style.Tweet.Reposted)

				metrics = append(metrics, text)
			}

			// リアクション
			for _, r := range p.Reactions {
				if text := l.createMetricsStr(&r, l.Style.Tweet.Like, l.Style.Tweet.Liked); text != "" {
					metrics = append(metrics, text)
				}
			}

			return strings.Join(metrics, " "), nil
		},
	}

	t, err := template.New("post").Funcs(funcMap).Parse(l.Template.Post)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, *p); err != nil {
		return err
	}

	// 不要な改行を除く
	post := strings.TrimSpace(buf.String())
	fmt.Fprintln(l.Writer, post)

	return nil
}

func (l *Layout) printAnnotation(t ...string) {
	text := strings.Join(t, " ")
	fmt.Fprintln(l.Writer, createStyledText(l.Style.Tweet.Annotation, text, ""))
}

func (l *Layout) createPostStr(p *sharedapi.Post) string {
	text := p.Text

	// メンションをハイライト
	styledMention := createStyledText(l.Style.Tweet.Mention, "$1@$2", "")
	text = regexp.MustCompile(`(^|[^\w@#$%&/])@(\w+)`).ReplaceAllString(text, styledMention)

	// ハッシュタグをハイライト
	for _, tag := range p.Tags {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)[#＃](%s\s|%s$)`, tag.Name, tag.Name))
		styledHashtag := createStyledText(l.Style.Tweet.HashTag, "#$1", tag.URL)
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

	return createStyledText(l.Style.Tweet.Detail, buf.String(), ""), nil
}

func (l *Layout) createMetricsStr(r *sharedapi.Reaction, normalStyle, revStyle string) string {
	if r.Count == 0 {
		return ""
	}

	style := normalStyle
	if r.Reacted {
		style = revStyle
	}

	text := fmt.Sprintf("%s %d", r.Name, r.Count)
	return createStyledText(style, text, "")
}
