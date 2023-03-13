package shared

import (
	"fmt"
	"io"
)

type AuthResponse struct {
	UserID   string
	UserName string
	Token    string
}

func PrintAuthURL(w io.Writer, u string) {
	fmt.Fprintf(w, "Please access the following URL to approve the application\n\n%s\n", u)
}
