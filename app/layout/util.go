package layout

import (
	"fmt"
	"time"

	"github.com/mattn/go-runewidth"
)

func (l *Layout) convertDateString(createAt time.Time) string {
	format := ""

	// 今日の日付なら時刻のみを表示
	if isSameDate(createAt) {
		format = l.Appearancene.TimeFormat
	} else {
		format = fmt.Sprintf("%s %s", l.Appearancene.DateFormat, l.Appearancene.TimeFormat)
	}

	return createAt.Local().Format(format)
}

func isSameDate(t time.Time) bool {
	now := time.Now()
	location := now.Location()
	fixedTime := t.In(location)

	t1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	t2 := time.Date(fixedTime.Year(), fixedTime.Month(), fixedTime.Day(), 0, 0, 0, 0, location)

	return t1.Equal(t2)
}

// CreateStyledText : スタイルありのテキストを作成
func CreateStyledText(style, text, url string) string {
	return fmt.Sprintf("[%s:%s]%s[-:-:-:-]", style, url, text)
}

// Truncate : …で省略
func Truncate(s string, l int) string {
	return runewidth.Truncate(s, l, "…")
}
