package github

import (
	"context"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/micnncim/action-lgtm-reaction/pkg/pointer"
)

type Client struct {
	githubClient *github.Client
	log          *zap.Logger
}

func NewClient(token string) (*Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Client{
		githubClient: github.NewClient(tc),
		log:          log,
	}, nil
}

func (c *Client) CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) error {
	log := c.log.With(
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.Int("number", number),
		zap.String("body", body),
	)

	_, _, err := c.githubClient.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
		Body: pointer.String(body),
	})
	if err != nil {
		log.Error("unable to create issue", zap.Error(err))
		return err
	}
	return nil
}
