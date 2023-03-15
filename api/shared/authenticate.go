package shared

import (
	"fmt"
	"io"
)

const (
	// AuthCallbackAddr : コールバックを待つサーバーのアドレス
	AuthCallbackAddr = "localhost:3000"
	// AuthCallbackURL : コールバックURL
	AuthCallbackURL = "http://localhost:3000/callback"
)

// User : ユーザー情報
type User struct {
	UserID   string
	UserName string
	Token    string
}

// PrintAuthenticateURL : 認証用URLを出力
func PrintAuthenticateURL(w io.Writer, u string) {
	fmt.Fprintf(w, "Please access the following URL to approve the application\n\n%s\n", u)
}
