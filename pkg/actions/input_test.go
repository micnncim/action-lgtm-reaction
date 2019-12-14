package actions

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetInput(t *testing.T) {
	var (
		envInputTrigger  = "INPUT_TRIGGER"
		envInputOverride = "INPUT_OVERRIDE"
		envInputSource   = "INPUT_SOURCE"
	)
	orgInputTrigger := os.Getenv(envInputTrigger)
	orgInputOverride := os.Getenv(envInputOverride)
	orgInputSource := os.Getenv(envInputSource)
	defer func() {
		os.Setenv(envInputTrigger, orgInputTrigger)
		os.Setenv(envInputOverride, orgInputOverride)
		os.Setenv(envInputSource, orgInputSource)
	}()

	tests := []struct {
		name string
		want Input
	}{
		{
			name: "ok",
			want: Input{
				Trigger:  "trigger",
				Override: true,
				Source:   "lgtmapp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(envInputTrigger, tt.want.Trigger)
			os.Setenv(envInputOverride, fmt.Sprint(tt.want.Override))
			os.Setenv(envInputSource, tt.want.Source)

			if got := GetInput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
