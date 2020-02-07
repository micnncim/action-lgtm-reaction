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

package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/micnncim/action-lgtm-reaction/pkg/lgtm"
)

const gifURLFormat = "https://media.giphy.com/media/%s/giphy.gif"

const apiBaseURLFormat = "https://api.giphy.com/v1/%s"

type Giphy struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Payload struct {
	Data []*Giphy `json:"data"`
}

type client struct {
	httpClient *http.Client
	apiKey     string
	log        *zap.Logger
}

var _ lgtm.Client = (*client)(nil)

func NewClient(apiKey string) (*client, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
		log:    log,
	}, nil
}

func (c *client) GetRandom() (string, error) {
	giphies, err := c.search("lgtm")
	if err != nil {
		return "", nil
	}
	if len(giphies) == 0 {
		return "", errors.New("no giphy")
	}
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(giphies))
	gifURL := fmt.Sprintf(gifURLFormat, giphies[index].ID)
	return lgtm.MarkdownStyle(gifURL), nil
}

func (c *client) search(q string) ([]*Giphy, error) {
	log := c.log.With(zap.String("q", q))

	apiURL := fmt.Sprintf(apiBaseURLFormat, "gifs/search")
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		log.Error("unable to make http request", zap.Error(err))
		return nil, err
	}
	query := req.URL.Query()
	query.Add("api_key", c.apiKey)
	query.Add("q", q)
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("unable to do http request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	payload := &Payload{}
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		log.Error("unable to decode json", zap.Error(err))
		return nil, err
	}

	return payload.Data, nil
}
