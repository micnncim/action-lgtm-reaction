package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/micnncim/action-lgtm-reaction/pkg/giphy"
	"github.com/micnncim/action-lgtm-reaction/pkg/github"
)

var (
	githubToken       = os.Getenv("GITHUB_TOKEN")
	giphyAPIKey       = os.Getenv("GIPHY_API_KEY")
	githubRepository  = os.Getenv("GITHUB_REPOSITORY")
	githubCommentBody = os.Getenv("GITHUB_COMMENT_BODY")
	githubIssueNumber = os.Getenv("GITHUB_ISSUE_NUMBER")
	trigger           = os.Getenv("INPUT_TRIGGER")
)

func main() {
	if strings.ToUpper(trigger) != strings.ToUpper(githubCommentBody) {
		fmt.Fprintf(os.Stderr, "no match issue comment\n")
		return
	}

	giphyClient, err := giphy.NewClient(giphyAPIKey)
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

	githubClient, err := github.NewClient(githubToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create github client: %v\n", err)
		os.Exit(1)
	}

	slugs := strings.Split(githubRepository, "/")
	if len(slugs) != 2 {
		fmt.Fprintf(os.Stderr, "invalid githubRepository: %v\n", githubRepository)
		os.Exit(1)
	}
	owner, repo := slugs[0], slugs[1]
	number, err := strconv.Atoi(githubIssueNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to convert string to int in issue number\n")
		os.Exit(1)
	}
	ctx := context.Background()

	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(giphies))
	comment := giphies[index].GIFURLInMarkdownStyle()

	if err := githubClient.CreateIssueComment(ctx, owner, repo, number, comment); err != nil {
		os.Exit(1)
	}
}
