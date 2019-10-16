package lgtmapp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/micnncim/action-lgtm-reaction/pkg/lgtm"
)

const (
	randomURL      = "https://www.lgtm.app/g"
	imageURLFormat = "https://www.lgtm.app/p/%s"
)

type client struct {
	httpClient *http.Client
	log        *zap.Logger
}

var _ lgtm.Client = (*client)(nil)

func NewClient() (*client, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}, nil
}

func (c *client) GetRandom() (string, error) {
	log := c.log

	req, err := http.NewRequest(http.MethodGet, randomURL, nil)
	if err != nil {
		log.Error("unable to create new http request", zap.Error(err))
		return "", nil
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("unable to do http request", zap.Error(err))
		return "", nil
	}

	redirectedURL := resp.Request.URL.String()
	slugs := strings.Split(redirectedURL, "/")
	id := slugs[len(slugs)-1]
	imageURL := fmt.Sprintf(imageURLFormat, id)
	return lgtm.MarkdownStyle(imageURL), nil
}
