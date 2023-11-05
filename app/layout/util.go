package layout

import (
	"fmt"
	"os"
	"time"

	"github.com/arrow2nd/nekomata/app/exit"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// convertDateString : 日付文字列に変換
func (l *Layout) convertDateString(createAt time.Time) string {
	format := ""

	// 今日の日付なら時刻のみを表示
	if isSameDate(createAt) {
		format = l.Appearance.FormatTime
	} else {
		format = fmt.Sprintf("%s %s", l.Appearance.FormatDate, l.Appearance.FormatTime)
	}

	return createAt.Local().Format(format)
}

// GetWindowWidth : 表示領域の幅を取得
func GetWindowWidth() int {
	fd := int(os.Stdout.Fd())

	w, _, err := term.GetSize(fd)
	if err != nil {
		exit.Error(err.Error(), exit.CodeErrTerm)
	}

	return w - 2
}

// CreateStyledText : スタイルありのテキストを作成
func CreateStyledText(style, text, url string) string {
	return fmt.Sprintf("[%s:%s]%s[-:-:-:-]", style, url, text)
}

// Truncate : …で省略
func Truncate(s string, l int) string {
	return runewidth.Truncate(s, l, "…")
}

// isSameDate : 日付が同じかどうか
func isSameDate(t time.Time) bool {
	now := time.Now()
	location := now.Location()
	fixedTime := t.In(location)

	t1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	t2 := time.Date(fixedTime.Year(), fixedTime.Month(), fixedTime.Day(), 0, 0, 0, 0, location)

	return t1.Equal(t2)
}
