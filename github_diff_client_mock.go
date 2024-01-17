package main

import (
	"context"

	"github.com/google/go-github/v57/github"
	diff "github.com/kmesiab/go-github-diff"
)

type GithubDiffClientInterface interface {
	ParsePullRequestURL(pullRequestURL string) (*diff.PullRequestURL, error)
	GetPullRequest(ctx context.Context, pr *diff.PullRequestURL, client *github.Client) (string, error)
	ParseGitDiff(diff string, ignoreList []string) []*diff.GitDiff
}

type MockGithubClient struct{}

// ParsePullRequestURL is a pass-through to the
// real client.  It doesn't need to be mocked.
func (c *MockGithubClient) ParsePullRequestURL(url string) (*diff.PullRequestURL, error) {
	return diff.ParsePullRequestURL(url)
}

func (c *MockGithubClient) GetPullRequest(
	_ context.Context,
	_ *diff.PullRequestURL,
	_ *github.Client,
) (string, error) {
	return `
diff --git a/file1.txt b/file1.txt
index abcdef1..1234567 100644
--- a/file1.txt
+++ b/file1.txt
@@ -1,4 +1,4 @@
 This is the old content
-This is a line added in the PR
+This is a line added in the PR - Modified
 This is the rest of the content
`, nil
}

func (c *MockGithubClient) ParseGitDiff(_ string, _ []string) []*diff.GitDiff {
	// Create a mock GitDiff slice for testing
	mockDiffs := []*diff.GitDiff{
		{
			FilePathOld: "file1.txt",
			FilePathNew: "file1.txt",
			Index:       "abcdef1..1234567",
			DiffContents: `
This is the old content
-This is a line added in the PR
+This is a line added in the PR - Modified
This is the rest of the content
`,
		},
		{
			FilePathOld: "file2.txt",
			FilePathNew: "file2.txt",
			Index:       "1234567..abcdef1",
			DiffContents: `
Another file diff
+Added a line
-Removed a line
`,
		},
	}

	return mockDiffs
}
