package layout

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

func (l *Layout) createUser(i int, u *sharedapi.Account) (string, error) {
	width, _ := GetWindowSize()
	name := Truncate(u.DisplayName, width)

	// DisplayNameが無い場合はUsernameを使う
	if name == "" {
		name = Truncate(u.Username, width)
	}

	// カーソル選択用のタグを埋め込む
	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, CreatePostHighlightTag(i), name)
	}

	funcMap := template.FuncMap{
		"displayName": func() string {
			return CreateStyledText(l.Style.User.Name, name, "")
		},
		"username": func() string {
			return CreateStyledText(l.Style.User.UserName, Truncate("@"+u.Username, width/2), "")
		},
		"badges": func() string {
			badges := []string{}

			if u.Verified {
				badges = append(
					badges,
					CreateStyledText(l.Style.User.Verified, l.Text.Verified, ""),
				)
			}

			if u.Private {
				badges = append(
					badges,
					CreateStyledText(l.Style.User.Private, l.Text.Private, ""),
				)
			}

			return strings.Join(badges, " ")
		},
	}

	t, err := template.New("User").Funcs(funcMap).Parse(l.Template.User)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, *u); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
