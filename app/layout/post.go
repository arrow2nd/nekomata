package layout

import (
	"bytes"
	"fmt"
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
	if p.Reference != nil {
		p = p.Reference
	}

	funcMap := template.FuncMap{
		"author": func() (string, error) {
			return l.createUserStr(i, p.Author)
		},
		"detail": func() (string, error) {
			return l.createpostDetail(p)
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

func (l *Layout) createpostDetail(p *sharedapi.Post) (string, error) {
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

	return createStyledText(l.Style.Tweet.Detail, buf.String()), nil
}
