# Auto Commit and Push Utility

This utility automates the process of committing and pushing code changes to your Git repository. It generates a commit message based on the changes using OpenAI's GPT model or falls back to a default message if no API key is provided.

## Prerequisites

- **Git**: Ensure you have Git installed and initialized in your project.
- **Go**: This utility is written in Go. Ensure you have Go installed.
- **OpenAI API Key**: If you want to use OpenAI for generating commit messages, set up an API key by creating a `.env` file in the root of your project and adding your key like this:
  ```
  OPENAI_API_KEY=your-openai-api-key
  ```

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/your-username/your-repo.git
cd your-repo
```

## Build

To build the utility, simply run the following command:

```bash
make
```

This will compile the Go code and generate an executable.

## Usage

To run the utility, execute the following command:

```bash
./your-executable
```

The utility will perform the following steps:
1. Add all changes using `git add .`.
2. Generate a commit message based on the changes using OpenAI's GPT model (if available) or use the default message `🚀`.
3. Commit the changes with the generated message.
4. Push the changes to the remote repository.

## Example Output

```bash
Running command: git add .
Running command: git diff --staged
Generated commit message:
Refactored function to improve readability

Running command: git commit -m "Refactored function to improve readability"
Commit successful. Output:
[main 1f4e9d4] Refactored function to improve readability
 1 file changed, 3 insertions(+), 1 deletion(-)

Running command: git push
Push successful. Output:
To https://github.com/your-username/your-repo.git
   1f4e9d4..a1b2c3d  main -> main
Changes committed and pushed successfully.
```

## Troubleshooting

- **OpenAI API Key Not Set**: If the API key is not set, the utility will use the default commit message `🚀`.
- **Git Errors**: If there are issues with committing or pushing, ensure your Git setup is correct, and that you have the necessary permissions.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to customize and extend this utility as per your needs!