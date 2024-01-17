package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	goenv "github.com/Netflix/go-env"
	"github.com/google/go-github/v57/github"
	gh "github.com/kmesiab/go-github-diff"
	"github.com/mkideal/cli"
)

type argT struct {
	URL    string `cli:"*url" usage:"The GitHub pull request URL"`
	ApiKey string `cli:"key" env:"OPENAI_API_KEY"  usage:"Your OpenAI API key. Leave this blank to use environment variable OPENAI_API_KEY"`
}

type ReviewedDiff struct {
	Diff   gh.GitDiff `json:"diff"`
	Review string     `json:"review"`
	Model  string     `json:"model"`
	Error  error      `json:"error"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		mergedArgs, err := coalesceConfiguration(argv)
		if err != nil {
			return fmt.Errorf("couldn't figure out the configuration. %s", err)
		}

		if mergedArgs.ApiKey == "" {
			return fmt.Errorf(
				"no API key provided, either pass it with the `--key` flag " +
					"or set the OPENAI_API_KEY environment variable")
		}

		parsedDiffFiles, err := processPullRequest(mergedArgs.URL, &GithubDiffClient{})
		if err != nil {
			return err
		}

		fmt.Printf("Processing %d diff files.  This may take a while...\n", len(parsedDiffFiles))

		reviews, err := getCodeReviews(parsedDiffFiles, "gpt-4", mergedArgs.ApiKey, &OpenAICompletionService{})
		if err != nil {
			return err
		}

		for _, review := range reviews {

			if review.Error != nil {
				fmt.Printf("ERROR: couldn't get the review for %s:  %s\n",
					path.Base(review.Diff.FilePathNew),
					review.Error,
				)

				continue
			}

			filename := "reviews/" + path.Base(review.Diff.FilePathNew) + ".md"
			err := os.WriteFile(filename, []byte(review.Review), 0o644)

			fmt.Printf("Saved review to %s\n", filename)

			if err != nil {
				fmt.Printf("couldn't save the review for %s:  %s",
					filename,
					err,
				)

				continue
			}
		}

		return nil
	}))
}

func getCodeReviews(diffs []*gh.GitDiff, model, apiKey string, svc CompletionServiceInterface) ([]ReviewedDiff, error) {
	reviewChan := make(chan ReviewedDiff)
	var reviews []ReviewedDiff

	for _, diff := range diffs {
		go func(d *gh.GitDiff) {
			fmt.Printf("Processing %s\n", path.Base(d.FilePathNew))

			review, err := svc.GetCompletion(d.DiffContents, model, apiKey)

			rDiff := ReviewedDiff{
				Error:  err,
				Diff:   *d,
				Review: review,
				Model:  model,
			}

			reviewChan <- rDiff
		}(diff)
	}

	for range diffs {
		review := <-reviewChan
		reviews = append(reviews, review)
	}

	close(reviewChan)

	return reviews, nil
}

func processPullRequest(prURL string, ghClient GithubDiffClientInterface) ([]*gh.GitDiff, error) {
	pullRequestUrl, err := ghClient.ParsePullRequestURL(prURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client := github.NewClient(nil)
	diffString, err := ghClient.GetPullRequest(ctx, pullRequestUrl, client)
	if err != nil {
		return nil, err
	}

	ignoreList := []string{
		".github",
		".gitignore",
		".travis.yml",
		"LICENSE",
		".md",
		".mod",
		".sum",
	}

	parsedDiff := ghClient.ParseGitDiff(diffString, ignoreList)

	return parsedDiff, nil
}

func coalesceConfiguration(cliArgs *argT) (*argT, error) {
	envArgs := &argT{}

	// Unmarshal environment variables into the envArgs struct
	_, err := goenv.UnmarshalFromEnviron(envArgs)
	if err != nil {
		return nil, err
	}

	// Default to the command line overriding
	// the environment variables
	if cliArgs.ApiKey == "" {
		cliArgs.ApiKey = envArgs.ApiKey
	}

	return cliArgs, nil
}