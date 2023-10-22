package layout

import (
	"io"

	"github.com/arrow2nd/nekomata/config"
)

type Layout struct {
	Writer       io.Writer
	Width        int
	Template     *config.Template
	Appearancene *config.Appearancene
	Icon         *config.Icon
	Style        *config.Style
}
