package layout

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/arrow2nd/nekomata/api/sharedapi"
)

func (l *Layout) createUserStr(i int, u *sharedapi.Account) (string, error) {
	name := truncate(u.DisplayName, l.Width)

	// DisplayNameが無い場合はUsernameを使う
	if name == "" {
		name = truncate(u.Username, l.Width)
	}

	// カーソル選択用のタグを埋め込む
	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, CreatePostHighlightTag(i), name)
	}

	funcMap := template.FuncMap{
		"displayName": func() string {
			return createStyledText(l.Style.User.Name, name, "")
		},
		"username": func() string {
			return createStyledText(l.Style.User.UserName, truncate("@"+u.Username, l.Width/2), "")
		},
		"badges": func() string {
			badges := []string{}

			if u.Verified {
				badges = append(
					badges,
					createStyledText(l.Style.User.Verified, l.Icon.Verified, ""),
				)
			}

			if u.Private {
				badges = append(
					badges,
					createStyledText(l.Style.User.Private, l.Icon.Private, ""),
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
