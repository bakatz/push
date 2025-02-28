package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const (
	maxCommitLength = 72
	defaultMessage  = "🚀"
)

type CommitMessageResponse struct {
	CommitMessage string `json:"commit_message" required:"true"`
}

func main() {
	_ = godotenv.Load() // Ignore error if .env file doesn't exist

	isInteractivePtr := flag.Bool("i", false, "interactive mode (edit commit message before pushing)")
	isDryRunPtr := flag.Bool("d", false, "dry-run mode (don't actually commit or push, just show what would have happened)")

	flag.Parse()

	isInteractive := false
	if isInteractivePtr != nil {
		isInteractive = *isInteractivePtr
	}

	isDryRun := false
	if isDryRunPtr != nil {
		isDryRun = *isDryRunPtr
	}

	if err := addAllChanges(); err != nil {
		fmt.Println("Could not add all changes due to an error:", err)
		return
	}

	if changes := getUnpushedChanges(); changes != "" {
		commitMessage := generateCommitMessage(changes)
		fmt.Printf("Generated commit message:\n%s\n", commitMessage)

		if err := commitChanges(commitMessage, isInteractive, isDryRun); err != nil {
			fmt.Println("Error committing changes:", err)
			return
		}
		fmt.Println("Changes committed successfully.")
	}

	if err := pushChanges(isDryRun); err != nil {
		fmt.Println("Error pushing changes:", err)
		return
	}

	fmt.Println("Changes pushed successfully.")
}

func addAllChanges() error {
	_, err := runCommandOrDryRun(exec.Command("git", "add", "."), false)
	return err
}

func getUnpushedChanges() string {
	cmd := exec.Command("git", "diff", "--staged")
	fmt.Printf("Running command: %s\n", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting uncommitted changes:", err)
		return ""
	}
	return string(output)
}

func generateCommitMessage(changes string) string {
	apiKey := os.Getenv("PUSH_OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("PUSH_OPENAI_API_KEY not set. Using default commit message.")
		return defaultMessage
	}

	description := getChangeDescription(changes, apiKey)
	return formatCommitMessage(description)
}

func getChangeDescription(changes, apiKey string) string {
	openaiConfig := openai.DefaultConfig(apiKey)
	if openaiBaseURL := os.Getenv("PUSH_OPENAI_BASE_URL"); openaiBaseURL != "" {
		openaiConfig.BaseURL = openaiBaseURL
		fmt.Printf("Custom OpenAI base URL detected. Using %s as the API Base URL.\n", openaiConfig.BaseURL)
	}

	model := openai.GPT4oMini
	if customModel := os.Getenv("PUSH_OPENAI_MODEL"); customModel != "" {
		model = customModel
		fmt.Printf("Custom OpenAI model detected. Using %s as the model.\n", model)
	}

	client := openai.NewClientWithConfig(openaiConfig)
	prompt := fmt.Sprintf("Analyze the following git diff and provide a concise description of the changes (%d characters or less):\n\n%s", maxCommitLength, changes)

	var schemaToUse CommitMessageResponse
	schema, err := jsonschema.GenerateSchemaForType(schemaToUse)
	if err != nil {
		fmt.Printf("Failed to generate schema for type, exiting: %v\n", err)
		os.Exit(-1)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that analyzes git diffs and generates concise commit messages.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "commit_message_schema",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v, using default message for this commit\n", err)
		return defaultMessage
	}

	if len(resp.Choices) > 0 {
		respObj := &CommitMessageResponse{}
		if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), respObj); err != nil || respObj.CommitMessage == "" {
			fmt.Println("ChatCompletion response could not be parsed as JSON or the generated message was empty, using default message for this commit")
			return defaultMessage
		}

		return respObj.CommitMessage
	} else {
		fmt.Println("ChatCompletion response was in an unexpected format, using default message for this commit")
		return defaultMessage
	}
}

func formatCommitMessage(description string) string {
	description = strings.TrimSpace(description)

	if len(description) > maxCommitLength {
		description = description[:maxCommitLength-3] + "..."
	}
	return description
}

func commitChanges(message string, isInteractive, isDryRun bool) error {
	if _, err := runCommandOrDryRun(exec.Command("git", "commit", "-m", message), isDryRun); err != nil {
		return err
	}

	if isInteractive {
		cmd := exec.Command("git", "commit", "--amend")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := runCommandOrDryRunNoOutput(cmd, isDryRun); err != nil {
			return err
		}
	}

	return nil
}

func pushChanges(isDryRun bool) error {
	_, err := runCommandOrDryRun(exec.Command("git", "push"), isDryRun)
	return err
}

func runCommandOrDryRun(cmd *exec.Cmd, isDryRun bool) (string, error) {
	if isDryRun {
		fmt.Printf("Dry-run mode enabled. Skipping command: %s\n", cmd.String())
		return "", nil
	} else {
		fmt.Printf("Running command: %s\n", cmd.String())
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Command output:\n%s\n", string(output))
			return string(output), err
		}
		fmt.Printf("Command successful. Output:\n%s\n", string(output))
		return string(output), nil
	}
}

func runCommandOrDryRunNoOutput(cmd *exec.Cmd, isDryRun bool) error {
	if isDryRun {
		fmt.Printf("Dry-run mode enabled. Skipping command: %s\n", cmd.String())
		return nil
	} else {
		fmt.Printf("Running command: %s\n", cmd.String())
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Println("Command successful.")
		return nil
	}
}
