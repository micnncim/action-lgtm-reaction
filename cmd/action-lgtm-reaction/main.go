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
	githubToken   = os.Getenv("GITHUB_TOKEN")
	giphyAPIKey   = os.Getenv("GIPHY_API_KEY")
	githubContext = os.Getenv("GITHUB")
	trigger       = os.Getenv("INPUT_TRIGGER")
	override      = os.Getenv("INPUT_OVERRIDE")
)

type GitHubContext struct {
	Repository string `json:"repository"`
	Event      struct {
		Issue struct {
			Number int `json:"number"`
		}
		Comment struct {
			ID   int    `json:"id"`
			Body string `json:"body"`
		} `json:"comment"`
		PullRequest struct {
			Number int `json:"number"`
		} `json:"pull_request"`
		Review struct {
			ID   int    `json:"id"`
			Body string `json:"body"`
		} `json:"review"`
	} `json:"event"`
}

func main() {
	ghContext, err := parseGitHubContext(githubContext)
	if err != nil {
		exit("unable to parse github context: %v\n", err)
	}

	needOverride, err := strconv.ParseBool(override)
	if err != nil {
		exit("unable to parse string to bool in override flag: %v\n", err)
	}

	matchComment, err := matchTrigger(trigger, ghContext.Event.Comment.Body)
	if err != nil {
		exit("invalid trigger: %v\n", err)
	}
	matchReview, err := matchTrigger(trigger, ghContext.Event.Review.Body)
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

	slugs := strings.Split(ghContext.Repository, "/")
	if len(slugs) != 2 {
		exit("invalid githubRepository: %v\n", ghContext.Repository)
	}
	owner, repo := slugs[0], slugs[1]

	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(giphies))
	comment := giphies[index].GIFURLInMarkdownStyle()

	ctx := context.Background()

	if needUpdateComment {
		if err := githubClient.UpdateIssueComment(ctx, owner, repo, ghContext.Event.Comment.ID, comment); err != nil {
			exit("unable to update issue comment: %v\n", err)
		}
		return
	}

	if needCreateComment {
		if err := githubClient.CreateIssueComment(ctx, owner, repo, ghContext.Event.Issue.Number, comment); err != nil {
			exit("unable to create issue comment: %v\n", err)
		}
		return
	}

	if needUpdateReview {
		if err := githubClient.UpdateReview(ctx, owner, repo, ghContext.Event.PullRequest.Number, ghContext.Event.Review.ID, comment); err != nil {
			exit("unable to update review: %v\n", err)
		}
		return
	}
}

func exit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func parseGitHubContext(s string) (*GitHubContext, error) {
	gc := &GitHubContext{}
	if err := json.Unmarshal([]byte(s), gc); err != nil {
		return nil, err
	}
	return gc, nil
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
