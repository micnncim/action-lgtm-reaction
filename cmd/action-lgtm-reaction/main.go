package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
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

	matchComment, err := matchTrigger(trigger, githubCommentBody)
	if err != nil {
		exit("invalid trigger: %v\n", err)
	}
	matchReview, err := matchTrigger(trigger, githubReviewBody)
	if err != nil {
		exit("invalid trigger: %v\n", err)
	}

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

// trigger is expected as JSON array like '["a", "b"]'.
func parseTrigger(trigger string) ([]string, error) {
	var a []string
	if err := json.Unmarshal([]byte(trigger), &a); err != nil {
		return nil, err
	}
	return a, nil
}

func matchTrigger(trigger, target string) (bool, error) {
	regexps, err := parseTrigger(trigger)
	if err != nil {
		return false, err
	}
	for _, s := range regexps {
		r := regexp.MustCompile(s)
		if r.MatchString(target) {
			return true, nil
		}
	}
	return false, nil
}
