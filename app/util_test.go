package app

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestOpenExternalEditor(t *testing.T) {
	a := &App{}

	t.Run("エディタの起動コマンドが未指定の場合にエラーが返るか", func(t *testing.T) {
		err := a.openExternalEditor("")
		assert.EqualError(t, err, "please specify which editor to use")
	})
}

func TestGetHighlightId(t *testing.T) {
	tests := []struct {
		name string
		arg  []string
		want int
	}{
		{
			name: "1桁のIDを抽出",
			arg:  []string{"page_0"},
			want: 0,
		},
		{
			name: "2桁のIDを抽出",
			arg:  []string{"post_10"},
			want: 10,
		},
		{
			name: "3桁のIDを抽出",
			arg:  []string{"post_100"},
			want: 100,
		},
		{
			name: "nilが渡された",
			arg:  nil,
			want: -1,
		},
		{
			name: "IDがない",
			arg:  []string{"page_"},
			want: -1,
		},
		{
			name: "_がない",
			arg:  []string{"rinze"},
			want: -1,
		},
		{
			name: "解析できない形式",
			arg:  []string{"asahi_serizawa"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getHighlightId(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		s     []int
		f     func(int) bool
		want  int
		found bool
	}{
		{
			name: "条件を満たす",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e == 1
			},
			want:  0,
			found: true,
		},
		{
			name: "条件を満たさない",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e > 4
			},
			want:  -1,
			found: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, got := find(tt.s, tt.f)
			assert.Equal(t, tt.want, index)
			assert.Equal(t, tt.found, got)
		})
	}
}

func TestTruncate(t *testing.T) {
	runewidth.DefaultCondition.EastAsianWidth = false

	tests := []struct {
		name string
		s    string
		w    int
		want string
	}{
		{
			name: "そのままの文字列が返る",
			s:    "komiya_kaho",
			w:    20,
			want: "komiya_kaho",
		},
		{
			name: "丸められた文字列が返る",
			s:    "shirase_sakuyasan",
			w:    15,
			want: "shirase_sakuya…",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.s, tt.w)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "正しく分割できるか（英語）",
			s:    "komiya kaho",
			want: []string{"komiya", "kaho"},
		},
		{
			name: "正しく分割できるか（日本語）",
			s:    "小宮 果穂",
			want: []string{"小宮", "果穂"},
		},
		{
			name: "ダブルクオートで囲んだ部分が残るか",
			s:    `aketa mikoto "ikaruga luca" nanakusa nichika`,
			want: []string{"aketa", "mikoto", `ikaruga luca`, "nanakusa", "nichika"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := split(tt.s)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
