# Push 🚀

Never manually write a commit message again.

Longer explanation: this utility automates the process of committing and pushing code changes to your Git repository. It generates a commit message based on the changes using OpenAI's GPT model or falls back to "🚀" if no API key is provided.

by [@ben_makes_stuff](https://x.com/ben_makes_stuff)

## Prerequisites

- **Git**: Ensure you have Git installed and initialized in your project.
- **Go**: This utility is written in Go. Ensure you have Go installed.
- **OpenAI API Key**: If you want to use OpenAI for generating commit messages, set up an API key by creating a `.env` file in the root of your project and adding your key like this:
  ```
  OPENAI_API_KEY=your-openai-api-key
  ```

## Quick Start

1. Download the binary for your platform here: https://github.com/bakatz/push/releases
1. Then add it as an alias (replace darwin-arm64 with your computer's os and cpu architecture, and if desired replace `push` with something else like `ship`): `alias push=path/to/download/dir/push-darwin-arm64`
1. You're done! Just run `push` to start using it.

## Build

To build the utility, you need Golang installed (latest version from https://go.dev). Then simply run the following command:

```bash
make
```

This will compile the Go code and generate a bunch of cross-platform executables in the `build` directory.

Example binary path:
```bash
./build/push-darwin-arm64
```

## Usage

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

Feel free to customize and extend this utility!
