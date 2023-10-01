package app

import (
	"fmt"
	"os"
)

// exitCode : 終了コード
type exitCode int

// GetInt : 数値を取得
func (e exitCode) GetInt() int {
	return int(e)
}

const (
	// exitCodeOK : 正常
	exitCodeOK exitCode = iota
	// exitCodeErr : エラー
	exitCodeErr
	// exitCodeErrInit : 初期化エラー
	exitCodeErrInit exitCode = iota + 62
	// ExitCodeErrFileIO : 読み込みエラー
	exitCodeErrRead
	// exitCodeErrWrite : 書き込みエラー
	exitCodeErrWrite
	// exitCodeErrTerm : 端末関連のエラー
	exitCodeErrTerm
)

// exit : ログを出力して終了
func exit(s string) {
	fmt.Println(s)
	os.Exit(exitCodeOK.GetInt())
}

// exitError : エラーを出力して終了
func exitError(e string, c exitCode) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", e)
	os.Exit(c.GetInt())
}
