package giphy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const gifURLFormat = "https://media.giphy.com/media/%s/giphy.gif"

const apiBaseURLFormat = "https://api.giphy.com/v1/%s"

type Giphy struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	URL    string `json:"url"`
	Title  string `json:"title"`
	GIFURL string `json:"gifurl"`
}

func (g *Giphy) GIFURLInMarkdownStyle() string {
	return fmt.Sprintf("![](%s)", g.GIFURL)
}

type Payload struct {
	Data []*Giphy `json:"data"`
}

type Client struct {
	httpClient *http.Client
	apiKey     string
	log        *zap.Logger
}

func NewClient(apiKey string) (*Client, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
		log:    log,
	}, nil
}

func (c *Client) Search(q string) ([]*Giphy, error) {
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

	data := payload.Data
	for _, g := range data {
		g.GIFURL = fmt.Sprintf(gifURLFormat, g.ID)
	}
	return data, nil
}
