package layout

import (
	"io"

	"github.com/arrow2nd/nekomata/config"
)

type Layout struct {
	Writer     io.Writer
	Template   *config.Template
	Appearance *config.Appearance
	Text       *config.Text
	Style      *config.Style
}
