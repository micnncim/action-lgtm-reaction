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
