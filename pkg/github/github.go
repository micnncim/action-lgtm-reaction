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

package github

import (
	"context"

	"github.com/google/go-github/v28/github"
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
		log.Error("unable to create issue comment", zap.Error(err))
		return err
	}
	return nil
}

func (c *Client) UpdateIssueComment(ctx context.Context, owner, repo string, commentID int, body string) error {
	log := c.log.With(
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.Int("commentID", commentID),
		zap.String("body", body),
	)

	_, _, err := c.githubClient.Issues.EditComment(ctx, owner, repo, int64(commentID), &github.IssueComment{
		Body: pointer.String(body),
	})
	if err != nil {
		log.Error("unable to update issue comment", zap.Error(err))
		return err
	}
	return nil
}

func (c *Client) UpdateReview(ctx context.Context, owner, repo string, number, reviewID int, body string) error {
	log := c.log.With(
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.Int("number", number),
		zap.Int("reviewID", reviewID),
		zap.String("body", body),
	)

	_, _, err := c.githubClient.PullRequests.UpdateReview(ctx, owner, repo, number, int64(reviewID), body)
	if err != nil {
		log.Error("unable to update review", zap.Error(err))
		return err
	}
	return nil
}
