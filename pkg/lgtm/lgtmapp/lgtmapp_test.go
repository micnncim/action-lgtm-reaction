package lgtmapp

import (
	"testing"
)

func Test_stripImageURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				url: "https://www.lgtm.app/i/4F5vFPNW3",
			},
			want: "https://www.lgtm.app/p/4F5vFPNW3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOriginalImageURL(tt.args.url); got != tt.want {
				t.Errorf("getRandomImageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
