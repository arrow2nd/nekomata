package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// RecieveAuthenticateCode : サーバーを建てて認証コードを受け取る
func RecieveAuthenticateCode(queryName string, validator func(string) bool) (string, error) {
	mux := http.NewServeMux()
	recieved := make(chan string, 1)

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get(queryName)

		if validator(q) {
			recieved <- q
			w.Write([]byte("Authentication complete! You may close this page."))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		recieved <- ""
	})

	// サーバーを建ててリダイレクトを待機
	serve := http.Server{
		Addr:    AuthCallbackAddr,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- serve.ListenAndServe()
	}()

	code := <-recieved

	if err := serve.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("shutdown server error: %w", err)
	}

	if err := <-serverErr; err != http.ErrServerClosed {
		return "", fmt.Errorf("listen server error: %w", err)
	}

	if code == "" {
		return "", fmt.Errorf("failed to recieve %s", queryName)
	}

	return code, nil
}

// PrintAuthenticateURL : 認証用URLを出力
func PrintAuthenticateURL(w io.Writer, u string) {
	fmt.Fprintf(w, "Please access the following URL to approve the application\n\n%s\n", u)
}
