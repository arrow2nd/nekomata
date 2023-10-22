package exit

import (
	"fmt"
	"os"
)

// ExitCode : 終了コード
type ExitCode int

// GetInt : 数値を取得
func (e ExitCode) GetInt() int {
	return int(e)
}

const (
	// CodeOK : 正常
	CodeOK ExitCode = iota
	// CodeErr : エラー
	CodeErr
	// CodeErrInit : 初期化エラー
	CodeErrInit ExitCode = iota + 62
	// ExitCodeErrFileIO : 読み込みエラー
	CodeErrRead
	// CodeErrWrite : 書き込みエラー
	CodeErrWrite
	// CodeErrTerm : 端末関連のエラー
	CodeErrTerm
)

// OK : ログを出力して終了
func OK(s string) {
	fmt.Println(s)
	os.Exit(CodeOK.GetInt())
}

// Error : エラーを出力して終了
func Error(e string, c ExitCode) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", e)
	os.Exit(c.GetInt())
}
