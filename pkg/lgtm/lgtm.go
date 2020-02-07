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

package lgtm

import (
	"fmt"
)

type Source int

const (
	SourceInvalid Source = iota
	SourceGiphy
	SourceLGTMApp
)

func (s Source) String() string {
	switch s {
	case SourceInvalid:
		return ""
	case SourceGiphy:
		return "giphy"
	case SourceLGTMApp:
		return "lgtmapp"
	}
	return ""
}

type Client interface {
	GetRandom() (string, error)
}

func MarkdownStyle(url string) string {
	return fmt.Sprintf("![](%s)", url)
}
