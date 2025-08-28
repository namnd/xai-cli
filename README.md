# xai

A simple CLI wrapper for [xAI models](https://docs.x.ai/docs/overview) (e.g., Grok-3-mini, Grok-4) to boost coding productivity with features like basic code analysis, prompting, and chatting. At the moment, it's mainly built for the [xai.nvim](https://github.com/namnd/xai.nvim) Neovim plugin.

## Features

- **Analyze**: Scan files or directories and summarize codebase structure.
- **Chat**: Interactive conversation.
- **Prompt**: Send custom prompts to xAI for responses (used by xai.nvim)
- **Chat history**: View and resume past chats via a local SQLite database (use fzf for fuzzy search)

## Installation

1. Ensure you have Go installed (v1.19+).
2. Clone the repository:
   ```
   git clone https://github.com/namnd/xai-cli
   cd xai-cli
   go build -o ~/.local/bin/xai # or whatever folder in your $PATH 
   ```

## Setup

- Obtain an xAI API key from [x.ai](https://docs.x.ai/).
- Run `xai setup` to store it in `~/.xai/config`.
- It will also initialize a sqlit3 db file `~/.xai/chat.db`

## Usage

Run the CLI with subcommands. First, set up your API key:
```
xai setup
```

Examples:
- Analyze a directory: `./xai-cli analyze /path/to/code`
- Start a chat: `./xai-cli chat`
- View chat history: `./xai-cli chat history`
- Run a prompt: `./xai-cli prompt \"Explain Go pointers\"`

