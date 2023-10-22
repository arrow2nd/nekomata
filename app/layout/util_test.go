package layout

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSameDate(t *testing.T) {
	tests := []struct {
		name string
		arg  time.Time
		want bool
	}{
		{
			name: "現在の日時",
			arg:  time.Now(),
			want: true,
		},
		{
			name: "今日の日付",
			arg:  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
			want: true,
		},
		{
			name: "過去の日付",
			arg:  time.Date(2018, 4, 24, 0, 0, 0, 0, time.Local),
			want: false,
		},
		{
			name: "未来の日付",
			arg:  time.Now().Add(time.Hour * 24),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSameDate(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
