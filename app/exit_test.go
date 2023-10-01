package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	t.Run("終了コードが取得できる", func(t *testing.T) {
		code := exitCodeOK.GetInt()
		want := 0

		assert.Equal(t, want, code)
	})
}
