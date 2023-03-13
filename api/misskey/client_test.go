package misskey

import (
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := &shared.Config{
		Host: "example.com",
	}

	want := "https://example.com/api"
	client := New(c)

	assert.Equal(t, want, client.baseURL, "指定したホスト名を元にベースURLが作成できているか")
}
