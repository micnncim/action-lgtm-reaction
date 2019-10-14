package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/micnncim/action-lgtm-reaction/pkg/giphy"
	"github.com/micnncim/action-lgtm-reaction/pkg/github"
)

func main() {
	trigger := os.Getenv("INPUT_TRIGGER")
	givenComment := os.Getenv("GITHUB_COMMENT_BODY")
	if strings.ToUpper(trigger) != strings.ToUpper(givenComment) {
		fmt.Fprintf(os.Stderr, "no match issue comment\n")
		return
	}

	apiKey := os.Getenv("GIPHY_API_KEY")
	giphyClient, err := giphy.NewClient(apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create giphy client: %v\n", err)
		os.Exit(1)
	}
	giphies, err := giphyClient.Search("lgtm")
	if err != nil {
		os.Exit(1)
	}
	if len(giphies) == 0 {
		fmt.Fprintf(os.Stderr, "no giphy contents found\n")
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	githubClient, err := github.NewClient(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create github client: %v\n", err)
		os.Exit(1)
	}

	repository := os.Getenv("GITHUB_REPOSITORY")
	slugs := strings.Split(repository, "/")
	if len(slugs) != 2 {
		fmt.Fprintf(os.Stderr, "invalid repository: %v\n", repository)
		os.Exit(1)
	}
	owner, repo := slugs[0], slugs[1]
	issueNumber := os.Getenv("GITHUB_ISSUE_NUMBER")
	number, err := strconv.Atoi(issueNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to convert string to int in issue number\n")
		os.Exit(1)
	}
	ctx := context.Background()
	comment := giphies[0].GIFURLInMarkdownStyle()
	if err := githubClient.CreateIssueComment(ctx, owner, repo, number, comment); err != nil {
		os.Exit(1)
	}
}
