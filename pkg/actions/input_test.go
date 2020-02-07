// Copyright 2020 micnncim
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
