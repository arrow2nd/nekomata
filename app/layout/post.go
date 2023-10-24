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
	}

	t, err := template.New("post").Funcs(funcMap).Parse(l.Template.Post)
	if err != nil {
		return err
	}

	if err := t.Execute(l.Writer, *p); err != nil {
		return err
	}

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
