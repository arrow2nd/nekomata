package app

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/arrow2nd/nekomata/app/layout"
	"github.com/mattn/go-runewidth"
)

// openExternalEditor : 外部エディタを開く
func (a *App) openExternalEditor(editor string, args ...string) error {
	if editor == "" {
		return errors.New("please specify which editor to use")
	}

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var err error

	if global.isCLI {
		err = cmd.Run()
	} else {
		a.app.Suspend(func() {
			err = cmd.Run()
		})
		a.app.Draw()
	}

	if err != nil {
		return fmt.Errorf("failed to open editor (%s) : %w", editor, err)
	}

	return nil
}

// getStringDisplayRow : 文字列の表示行数を取得
func getStringDisplayRow(s string, w int) int {
	row := 0

	// タグを削除
	tagDeleted := replaceAll(s, `\[(.+?:.+?:.+?|".*?")\]`, "")

	for _, s := range strings.Split(tagDeleted, "\n") {
		r := int(math.Ceil(float64(runewidth.StringWidth(s)) / float64(w)))
		row += r
	}

	return row
}

// getHighlightId : ハイライト一覧からIDを取得（見つからない場合 -1 が返る）
func getHighlightId(ids []string) int {
	if ids == nil {
		return -1
	}

	i := strings.Index(ids[0], "_")
	if i == -1 || i+1 >= len(ids[0]) {
		return -1
	}

	id, err := strconv.Atoi(ids[0][i+1:])
	if err != nil {
		return -1
	}

	return id
}

// find : スライス内から任意の条件を満たす値を探す
func find[T any](s []T, f func(T) bool) (int, bool) {
	for i, v := range s {
		if f(v) {
			return i, true
		}
	}

	return -1, false
}

// truncate : 文字列を指定幅で丸める
func truncate(s string, w int) string {
	return runewidth.Truncate(s, w, "…")
}

// split : 文字列をスペースで分割（ダブルクオートで囲まれた部分は残す）
func split(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	return r.Read()
}

// replaceAll : 正規表現にマッチした文字列を一斉置換
func replaceAll(str, reg, rep string) string {
	replace := regexp.MustCompile(reg)
	return replace.ReplaceAllString(str, rep)
}

// createStatusMessage : ラベル付きステータスメッセージを作成
func createStatusMessage(label, status string) string {
	width, _ := layout.GetWindowSize()
	status = strings.ReplaceAll(status, "\n", " ")

	return truncate(fmt.Sprintf("[%s] %s", label, status), width)
}
