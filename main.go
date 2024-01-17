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

const (
	DefaultFilePerms   = 0o644
	DefaultOpenAIModel = "gpt-3.5-turbo-instruct"
)

type argT struct {
	URL             string `cli:"*url" usage:"The GitHub pull request URL"`
	ApiKey          string `cli:"key" env:"OPENAI_API_KEY"  usage:"Your OpenAI API key. Leave this blank to use environment variable OPENAI_API_KEY"`
	CompletionModel string `cli:"model" env:"CADRE_COMPLETION_MODEL" usage:"The OpenAI API model to use for code reviews."`
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

		printGreeting()

		mergedArgs, err := coalesceConfiguration(argv)
		if err != nil {
			return fmt.Errorf("couldn't figure out the configuration. %s", err)
		}

		if mergedArgs.ApiKey == "" {
			return fmt.Errorf(
				"no API key provided, either pass it with the `--key` flag " +
					"or set the OPENAI_API_KEY environment variable")
		}

		fmt.Printf("📡 Getting pull request from GitHub...\n")
		parsedDiffFiles, err := processPullRequest(mergedArgs.URL, &GithubDiffClient{})
		if err != nil {
			return err
		}

		fmt.Printf("\n⌛  Processing %d diff files.  This may take a while...\n\n", len(parsedDiffFiles))
		reviews, err := getCodeReviews(
			parsedDiffFiles,
			argv.CompletionModel,
			mergedArgs.ApiKey,
			&OpenAICompletionService{},
		)
		if err != nil {
			return err
		}

		saveReviews(reviews)
		fmt.Println("Done! 🏁")

		return nil
	}))
}

func saveReviews(reviews []ReviewedDiff) {
	for _, review := range reviews {

		if review.Error != nil {
			fmt.Printf("⚠️ couldn't get the review for %s:  %s\n",
				path.Base(review.Diff.FilePathNew),
				review.Error,
			)

			continue
		}

		filename := path.Base(review.Diff.FilePathNew) + ".md"

		// Ensure the directory exists
		dir := path.Dir(filename)

		// If it doesn't exist, create it
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, DefaultFilePerms); err != nil {
				fmt.Printf("⚠️ couldn't create directory for %s: %s\n", dir, err)

				continue
			}
		}

		// Save the review to disk
		err := saveReviewToFile(filename, review.Review)
		if err != nil {
			fmt.Printf("⚠️ couldn't save the review for %s:  %s\n",
				filename,
				err,
			)

			continue
		}

		fmt.Printf("💾 Saved review to %s\n", filename)
	}
}

func getCodeReviews(diffs []*gh.GitDiff, model, apiKey string, svc CompletionServiceInterface) ([]ReviewedDiff, error) {
	reviewChan := make(chan ReviewedDiff)
	var reviews []ReviewedDiff

	for _, diff := range diffs {
		go func(d *gh.GitDiff) {
			fmt.Printf("🤖 Getting code review for %s\n", path.Base(d.FilePathNew))

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

func saveReviewToFile(filename, reviewContent string) error {
	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("file %s already exists, not overwriting", filename)
	}

	// Write the review content to the file
	err := os.WriteFile(filename, []byte(reviewContent), DefaultFilePerms)
	if err != nil {
		return fmt.Errorf("failed to write review to file: %s", err)
	}

	return nil
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

	// If no model is provided, use the default model
	if cliArgs.CompletionModel == "" {
		cliArgs.CompletionModel = DefaultOpenAIModel
	}

	return cliArgs, nil
}

func printGreeting() {
	fmt.Println(`
 _____   ___ ____________ _____ 
/  __ \ / _ \|  _  \ ___ \  ___|
| /  \// /_\ \ | | | |_/ / |__  
| |    |  _  | | | |    /|  __| 
| \__/\| | | | |/ /| |\ \| |___ 
 \____/\_| |_/___/ \_| \_\____/ 

	`)
}
