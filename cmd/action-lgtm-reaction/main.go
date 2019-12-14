package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/micnncim/action-lgtm-reaction/pkg/github"
	"github.com/micnncim/action-lgtm-reaction/pkg/lgtm"
	"github.com/micnncim/action-lgtm-reaction/pkg/lgtm/giphy"
	"github.com/micnncim/action-lgtm-reaction/pkg/lgtm/lgtmapp"
)

var (
	githubToken = os.Getenv("GITHUB_TOKEN")
	giphyAPIKey = os.Getenv("GIPHY_API_KEY")
)

var (
	trigger  = os.Getenv("INPUT_TRIGGER")
	override = os.Getenv("INPUT_OVERRIDE")
	source   = os.Getenv("INPUT_SOURCE")
)

type GitHubEvent struct {
	Comment struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	} `json:"comment"`
	Issue struct {
		Number int `json:"number"`
	} `json:"issue"`
	PullRequest struct {
		Number int `json:"number"`
	} `json:"pull_request"`
	Review struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	} `json:"review"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	e, err := getGitHubEvent()
	if err != nil {
		return err
	}

	needCreateComment, needUpdateComment, needUpdateReview, err := checkActionNeeded(e)
	if err != nil {
		return err
	}
	if !needCreateComment && !needUpdateComment && !needUpdateReview {
		fmt.Fprintf(os.Stderr, "no need to do any action\n")
		return nil
	}

	owner, repo, err := getGitHubRepo()
	if err != nil {
		return err
	}

	lc, err := createLGTMClient()
	if err != nil {
		return err
	}
	lgtmComment, err := lc.GetRandom()
	if err != nil {
		return err
	}

	ctx := context.Background()

	gc, err := github.NewClient(githubToken)
	if err != nil {
		return err
	}

	switch {
	case needUpdateComment:
		if err := gc.UpdateIssueComment(ctx, owner, repo, e.Comment.ID, lgtmComment); err != nil {
			return err
		}
		return nil

	case needCreateComment:
		var number int
		switch {
		case e.Issue.Number != 0:
			number = e.Issue.Number
		case e.PullRequest.Number != 0:
			number = e.PullRequest.Number
		default:
			return errors.New("issue number or pull request number don't exist")
		}
		if err := gc.CreateIssueComment(ctx, owner, repo, number, lgtmComment); err != nil {
			return err
		}
		return nil

	case needUpdateReview:
		if err := gc.UpdateReview(ctx, owner, repo, e.PullRequest.Number, e.Review.ID, lgtmComment); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func createLGTMClient() (c lgtm.Client, err error) {
	switch source {
	case lgtm.SourceGiphy.String():
		c, err = giphy.NewClient(giphyAPIKey)
		return
	case lgtm.SourceLGTMApp.String():
		c, err = lgtmapp.NewClient()
		return
	default:
		err = errors.New("not support source")
		return
	}
}

func checkActionNeeded(e *GitHubEvent) (needCreateComment, needUpdateComment, needUpdateReview bool, err error) {
	needOverride, err := strconv.ParseBool(override)
	if err != nil {
		return
	}

	matchComment, err := matchTrigger(trigger, e.Comment.Body)
	if err != nil {
		return
	}
	matchReview, err := matchTrigger(trigger, e.Review.Body)
	if err != nil {
		return
	}

	needCreateComment = (matchComment || matchReview) && !needOverride
	needUpdateComment = matchComment && needOverride
	needUpdateReview = matchReview && needOverride

	return
}

func getGitHubEvent() (*GitHubEvent, error) {
	p := os.Getenv("GITHUB_EVENT_PATH")
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var e GitHubEvent
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func getGitHubRepo() (owner, repo string, err error) {
	r := os.Getenv("GITHUB_REPOSITORY")
	s := strings.Split(r, "/")
	if len(s) != 2 {
		err = fmt.Errorf("invalid githubRepository: %v\n", r)
		return
	}
	owner, repo = s[0], s[1]
	return
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
