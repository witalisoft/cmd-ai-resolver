# Command AI Resolver (cmd-ai-resolver)

![demo](demo.gif)

## Description

`cmd-ai-resolver` is a command-line application that translate your ideas into shell commands using AI. It identifies
special AI instruction tags (e.g., `<AI>your ai prompt</AI>`) within your shell commands and sends the entire
command-line context along with the extracted prompt to an OpenAI LLM, and replaces the tag with the LLM-generated
shell command segment. Tool is configured to be used with the `VISUAL` environment variable, so you can use
an well-known keyboard shortcut (C-xC-e) to resolve your command-line requests with AI. When there are no AI tags found
in the command-line it will run the passthrough command, so you can still use it as a regular editor for
your command-line.

## Features

- based on the OpenAI LLM API with an option to specify a custom base URL `OPENAI_BASE_URL` (for proxies or self-hosted compatible endpoints)
- supports the `gpt-4.1-mini` model by default, configurable via the `OPENAI_BASE_MODEL` environment variable
- option to passthrough the command-line without AI tags to the editor to preserve the same experience as before with `VISUAL`
- debug logging support with `-d` flag to see the AI processing steps
- bash and zsh compatible

## Installation

Currently, the application is run from the source or by building a binary.

## Usage

1. Prepare environment variables for OpenAI API at least the API key:

```bash
export OPENAI_API_KEY="your_openai_api_key_here"
```

2. Prepare your shell by creating the wrapper script `cmd-ai-resolver` with pass-through functionality to your editor (e.g., `vim`):

```bash
cat $HOME/cmd-ai-resolver-wrapper.sh
#!/usr/bin/env bash

$HOME/cmd-ai-resolver --pass-through vim $1
```

3. Set the `VISUAL` environment variable to use the `cmd-ai-resolver-wrapper` script:

```
export VISUAL="$HOME/cmd-ai-resolver-wrapper.sh"
```

4. Test the setup by running resolving a command with AI (C-xC-e):

```bash
ls -l /some/directory | <AI>filter for text files and show only the last 5 entries</AI>
```

## TODO

-   **Multiple AI Tag Processing:** Currently, `cmd-ai-resolver` processes only the first `<AI>...</AI>` tag found in the input file. Future versions aim to support processing multiple AI tags.
