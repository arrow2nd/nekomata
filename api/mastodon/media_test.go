package mastodon_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/arrow2nd/nekomata/api/mastodon"
	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/stretchr/testify/assert"
)

const testFilePath = "../../_testdata/strong_internet.png"

func TestUploadMedia(t *testing.T) {
	raw, _ := os.Open(testFilePath)
	defer raw.Close()

	fi, _ := raw.Stat()
	size := fi.Size()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, header, _ := r.FormFile("file")

		assert.Equal(t, filepath.Base(testFilePath), header.Filename)
		assert.Equal(t, size, header.Size)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id": "12345"}`)
	}))

	defer ts.Close()

	clientCred := &sharedapi.ClientCredential{
		Name:   "hoge",
		ID:     "fuga",
		Secret: "piyo",
	}

	m, err := mastodon.New(clientCred, &sharedapi.UserCredential{Server: ts.URL})
	assert.NoError(t, err)

	id, err := m.UploadMedia(filepath.Base(raw.Name()), raw)
	assert.NoError(t, err)

	assert.Equal(t, "12345", id)
}
