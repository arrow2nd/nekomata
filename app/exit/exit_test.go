package exit_test

import (
	"testing"

	"github.com/arrow2nd/nekomata/app/exit"
	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	t.Run("終了コードが取得できる", func(t *testing.T) {
		code := exit.CodeOK.GetInt()
		want := 0

		assert.Equal(t, want, code)
	})
}
