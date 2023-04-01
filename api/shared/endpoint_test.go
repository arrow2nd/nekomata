package shared_test

import (
	"net/url"
	"testing"

	"github.com/arrow2nd/nekomata/api/shared"
	"github.com/stretchr/testify/assert"
)

func TestEnpoint(t *testing.T) {
	host := "https://example.com"

	tests := []struct {
		name string
		e    shared.Endpoint
		p    url.Values
		want string
	}{
		{
			name: "単純な結合のみ",
			e:    "/api/hoge",
			p:    nil,
			want: "/api/hoge",
		},
		{
			name: "パスパラメータあり",
			e:    "/api/hoge/:id",
			p:    url.Values{":id": []string{"fuga"}},
			want: "/api/hoge/fuga",
		},
		{
			name: "クエリパラメータあり",
			e:    "/api/hoge",
			p:    nil,
			want: "/api/hoge?a=fuga&b=piyo",
		},
		{
			name: "パスパラメータ + クエリパラメータ",
			e:    "/api/hoge/:id",
			p:    url.Values{":id": []string{"fuga"}},
			want: "/api/hoge/fuga?a=piyo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.e.URL(host, tc.p)
			assert.Equal(t, host+tc.want, got)
		})
	}
}
