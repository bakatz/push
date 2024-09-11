package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

const (
	maxCommitLength = 72
	defaultMessage  = "🚀"
)

func main() {
	_ = godotenv.Load() // Ignore error if .env file doesn't exist

	isInteractive := false
	isInteractivePtr := flag.Bool("i", false, "interactive mode (edit commit message before pushing)")
	flag.Parse()

	if isInteractivePtr != nil {
		isInteractive = *isInteractivePtr
	}

	if err := addAllChanges(); err != nil {
		fmt.Println("Could not add all changes due to an error:", err)
		return
	}

	if changes := getUnpushedChanges(); changes != "" {
		commitMessage := generateCommitMessage(changes)
		fmt.Printf("Generated commit message:\n%s\n", commitMessage)

		if err := commitChanges(commitMessage, isInteractive); err != nil {
			fmt.Println("Error committing changes:", err)
			return
		}
		fmt.Println("Changes committed successfully.")
	}

	if err := pushChanges(); err != nil {
		fmt.Println("Error pushing changes:", err)
		return
	}

	fmt.Println("Changes pushed successfully.")
}

func addAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	fmt.Printf("Running command: %s\n", cmd.String())
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error adding all changes:", err)
		return err
	}
	return nil
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
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set. Using default commit message.")
		return defaultMessage
	}

	description := getChangeDescription(changes, apiKey)
	return formatCommitMessage(description)
}

func getChangeDescription(changes, apiKey string) string {
	client := openai.NewClient(apiKey)
	prompt := fmt.Sprintf("Analyze the following git diff and provide a concise description of the changes (72 characters or less):\n\n%s", changes)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
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
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v, using default message for this commit\n", err)
		return defaultMessage
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}

	return defaultMessage
}

func formatCommitMessage(description string) string {
	description = strings.TrimSpace(description)
	if len(description) > maxCommitLength {
		description = description[:maxCommitLength-3] + "..."
	}
	return description
}

func commitChanges(message string, isInteractive bool) error {
	cmd := exec.Command("git", "commit", "-m", message)
	fmt.Printf("Running command: %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command output:\n%s\n", string(output))
		return err
	}
	fmt.Printf("Commit successful. Output:\n%s\n", string(output))

	if isInteractive {
		cmd := exec.Command("git", "commit", "--amend")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Running command: %s\n", cmd.String())
		if err := cmd.Run(); err != nil {
			return err
		}

		fmt.Println("Commit amend successful.")
	}

	return nil
}

func pushChanges() error {
	cmd := exec.Command("git", "push")
	fmt.Printf("Running command: %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command output:\n%s\n", string(output))
		return err
	}
	fmt.Printf("Push successful. Output:\n%s\n", string(output))
	return nil
}
