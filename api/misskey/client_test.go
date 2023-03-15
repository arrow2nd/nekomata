package misskey

import (
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := &shared.ClientOpts{
		Server: "https://example.com",
	}

	want := "https://example.com"
	client := New(c)

	assert.Equal(t, want, client.opts.Server)
}
