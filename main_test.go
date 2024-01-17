package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoalesceConfiguration_Key_CLIOverride(t *testing.T) {
	// Set up to restore the value later
	currentEnvVar := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", currentEnvVar)

	// Set up an override
	err := os.Setenv("OPENAI_API_KEY", "dummy_key")
	require.NoError(t, err)

	args := &argT{ApiKey: "overridden_key"}
	result, _ := coalesceConfiguration(args)

	assert.Equal(t, "overridden_key", result.ApiKey)
}

func TestCoalesceConfiguration_Model_CLIOverride(t *testing.T) {
	// Set up to restore the value later
	currentEnvVar := os.Getenv("CADRE_COMPLETION_MODEL")
	defer os.Setenv("CADRE_COMPLETION_MODEL", currentEnvVar)

	// Set up an override
	err := os.Setenv("CADRE_COMPLETION_MODEL", "gpt-4")
	require.NoError(t, err)

	args := &argT{CompletionModel: "overridden_model"}
	result, _ := coalesceConfiguration(args)

	assert.Equal(t, "overridden_model", result.CompletionModel)
}

func TestProcessPullRequest(t *testing.T) {
	// Create a mock GitHub client
	mockClient := &MockGithubClient{}

	// Create mock argT with URL
	args := &argT{
		URL: "https://github.com/user/repo/pull/123",
	}

	diffs, err := processPullRequest(args.URL, mockClient)

	// Check if there's no error returned
	assert.NoError(t, err)

	fmt.Printf("%v", diffs)
}
