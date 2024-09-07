# Push 🚀

Never write a commit message again.

Longer explanation: this utility automates the process of committing and pushing code changes to your Git repository. It generates a commit message based on the changes using OpenAI's `gpt-4o-mini` model or falls back to "🚀" if no API key is provided.

by [@ben_makes_stuff](https://x.com/ben_makes_stuff)

[![Buy me a coffee](https://img.buymeacoffee.com/button-api/?text=Buy%20me%20a%20coffee&emoji=☕&slug=ben_makes_stuff&button_colour=FFDD00&font_colour=000000&font_family=Lato&outline_colour=000000&coffee_colour=ffffff)](https://www.buymeacoffee.com/ben_makes_stuff)


## Requirements
- **Git**: Ensure you have Git installed and initialized in your project.

## Quick Start

1. Download the binary for your platform here: https://github.com/bakatz/push/releases
1. Then add it as an alias (replace darwin-arm64 with your computer's os and cpu architecture, and if desired replace `push` with something else like `ship`): `alias push=path/to/download/dir/push-darwin-arm64`
1. You're done! Just run `push` to start using it.

Optionally: if you want to use OpenAI for generating commit messages, make sure you set your environment variable like this before running the above steps:
  ```
  OPENAI_API_KEY=your-openai-api-key
  ```

## Dev Requirements

- **Go**: This utility is written in Go. Ensure you have Go installed.

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
