#!/bin/bash

#  this runs droid in headless mode.
#  it sources the .env file to inject the env var for the factory API key:
#
#  export FACTORY_API_KEY=fk-...
#
#  the .env file isn't required if you have the env var set in your shell, but this can be a convenient way to
#  manage env vars for droid without having to set them manually.
#
#  the reason droid does not source the .env file itself is that it might read the file and put the API key in the context window.
#  a future release may allow excluding files from droid
#
#  this will run droid in spec mode using claude opus 4.6, and implement the code using gemini 3 flash.
#  feel free to change these to your preferred models.
#  it will prompt you for the path to your prompt file.
#
#  stream-json lets you see the agents execution in real time on the terminal.  you can also review agent.log file
#  after the run for debugging.

set -euo pipefail

while true; do
    read -rp "Enter the path to the prompt file: " PROMPT_FILE
    if [[ -z "$PROMPT_FILE" ]]; then
        echo "Error: no path entered. Please try again."
    elif [[ ! -f "$PROMPT_FILE" ]]; then
        echo "Error: file not found: $PROMPT_FILE. Please try again."
    elif [[ ! -r "$PROMPT_FILE" ]]; then
        echo "Error: no read permission: $PROMPT_FILE. Please try again."
    else
        break
    fi
done

source .env

PARENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

droid exec \
    --output-format stream-json \
    --use-spec \
    --spec-model claude-opus-4-6 \
    --model gemini-3-flash-preview \
    --reasoning-effort high \
    --auto high \
    --cwd "$PARENT_DIR" \
    --file "$PROMPT_FILE" | tee agent.log
