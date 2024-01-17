package main

import (
	"context"

	"github.com/google/go-github/v57/github"
	diff "github.com/kmesiab/go-github-diff"
)

// GithubDiffClient is a wrapper around the go-github-diff
// client, such that it implements the GithubDiffClientInterface
type GithubDiffClient struct{}

func (c *GithubDiffClient) ParsePullRequestURL(url string) (*diff.PullRequestURL, error) {
	return diff.ParsePullRequestURL(url)
}

func (c *GithubDiffClient) GetPullRequest(
	ctx context.Context,
	pullRequestURL *diff.PullRequestURL,
	ghClient *github.Client,
) (string, error) {
	return diff.GetPullRequest(ctx, pullRequestURL, ghClient)
}

func (c *GithubDiffClient) ParseGitDiff(diffString string, ignoreFiles []string) []*diff.GitDiff {
	return diff.ParseGitDiff(diffString, ignoreFiles)
}
