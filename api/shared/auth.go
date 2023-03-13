package shared

import (
	"fmt"
	"io"
)

// AuthResponse : アプリケーション認証のレスポンス
type AuthResponse struct {
	UserID   string
	UserName string
	Token    string
}

// PrintAuthURL : 認証用URLを出力
func PrintAuthURL(w io.Writer, u string) {
	fmt.Fprintf(w, "Please access the following URL to approve the application\n\n%s\n", u)
}
