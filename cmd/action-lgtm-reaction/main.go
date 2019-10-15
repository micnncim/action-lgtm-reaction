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
	githubToken             = os.Getenv("GITHUB_TOKEN")
	giphyAPIKey             = os.Getenv("GIPHY_API_KEY")
	githubRepository        = os.Getenv("GITHUB_REPOSITORY")
	githubIssueNumber       = os.Getenv("GITHUB_ISSUE_NUMBER")
	githubCommentBody       = os.Getenv("GITHUB_COMMENT_BODY")
	githubCommentID         = os.Getenv("GITHUB_COMMENT_ID")
	githubPullRequestNumber = os.Getenv("GITHUB_PULL_REQUEST_NUMBER")
	githubReviewBody        = os.Getenv("GITHUB_REVIEW_BODY")
	githubReviewID          = os.Getenv("GITHUB_REVIEW_ID")
	trigger                 = os.Getenv("INPUT_TRIGGER")
	override                = os.Getenv("INPUT_OVERRIDE")
)

func main() {
	needOverride, err := strconv.ParseBool(override)
	if err != nil {
		exit("unable to parse string to bool in override flag: %v\n", err)
	}

	matchComment := strings.ToUpper(trigger) == strings.ToUpper(githubCommentBody)
	matchReview := strings.ToUpper(trigger) == strings.ToUpper(githubReviewBody)

	needCreateComment := (matchComment || matchReview) && !needOverride
	needUpdateComment := matchComment && needOverride
	needUpdateReview := matchReview && needOverride

	if !needCreateComment && !needUpdateComment && !needUpdateReview {
		fmt.Fprintf(os.Stderr, "no need to do action\n")
		return
	}

	giphyClient, err := giphy.NewClient(giphyAPIKey)
	if err != nil {
		exit("unable to create giphy client: %v\n", err)
	}
	giphies, err := giphyClient.Search("lgtm")
	if err != nil {
		exit("unable to search giphy :%v\n", err)
	}
	if len(giphies) == 0 {
		exit("no giphy contents found\n")
	}

	githubClient, err := github.NewClient(githubToken)
	if err != nil {
		exit("unable to create github client: %v\n", err)
	}

	slugs := strings.Split(githubRepository, "/")
	if len(slugs) != 2 {
		exit("invalid githubRepository: %v\n", githubRepository)
	}
	owner, repo := slugs[0], slugs[1]

	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(giphies))
	comment := giphies[index].GIFURLInMarkdownStyle()

	ctx := context.Background()

	if needUpdateComment {
		commentID, err := strconv.ParseInt(githubCommentID, 10, 64)
		if err != nil {
			exit("unable to convert string to int in issue number: %v\n", err)
		}
		if err := githubClient.UpdateIssueComment(ctx, owner, repo, int(commentID), comment); err != nil {
			exit("unable to update issue comment: %v\n", err)
		}
		return
	}

	if needCreateComment {
		number, err := strconv.Atoi(githubIssueNumber)
		if err != nil {
			exit("unable to convert string to int in issue number: %v\n", err)
		}
		if err := githubClient.CreateIssueComment(ctx, owner, repo, number, comment); err != nil {
			exit("unable to create issue comment: %v\n", err)
		}
		return
	}

	if needUpdateReview {
		number, err := strconv.Atoi(githubPullRequestNumber)
		if err != nil {
			exit("unable to convert string to int in issue number: %v\n", err)
		}
		reviewID, err := strconv.Atoi(githubReviewID)
		if err != nil {
			exit("unable to convert string to int in review id: %v\n", err)
		}
		if err := githubClient.UpdateReview(ctx, owner, repo, number, reviewID, comment); err != nil {
			exit("unable to update review: %v\n", err)
		}
		return
	}
}

func exit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
