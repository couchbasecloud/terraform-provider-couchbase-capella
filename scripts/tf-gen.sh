#!/bin/bash

#  this runs droid in headless mode.  ensure the script has execute permission.
#
#  it sources the .env file to inject the env var for the factory API key:
#
#  export FACTORY_API_KEY=fk-...
#
#  the .env file isn't required if you have the env var set in your shell, but this can be a convenient way to
#  manage env vars for droid without having to set them manually.
#
#  this will run droid in spec mode using claude opus 4.6, and implement the code using gemini 3 flash.
#  feel free to change these to your preferred models.
#
#  it will prompt you for the path to your prompt file.
#
#  you can also review the session logs in ~/.factory/sessions/ to see the agent execution and for debugging.

set -euo pipefail

source .env

if [[ -z "${FACTORY_API_KEY}" ]]; then
  echo "Error: FACTORY_API_KEY is not set. Please set it in the shell or in the .env file."
  exit 1
fi

DEPENDENCIES=(droid go)
MISSING=()
for dep in "${DEPENDENCIES[@]}"; do
    if ! command -v "$dep" &>/dev/null; then
        MISSING+=("$dep")
    fi
done

if [[ ${#MISSING[@]} -gt 0 ]]; then
    echo "Error: the following dependencies are not installed or not in PATH:"
    for dep in "${MISSING[@]}"; do
        echo "  - $dep"
    done
    exit 1
fi


while true; do
    read -rp "Enter the path to the prompt file: " PROMPT_FILE
    # Expand ~ to $HOME since tilde expansion doesn't happen in variables
    PROMPT_FILE="${PROMPT_FILE/#\~/$HOME}"
    if [[ -z "$PROMPT_FILE" ]]; then
        echo "Error: no path entered. Please enter a path."
    elif [[ ! -f "$PROMPT_FILE" ]]; then
        echo "Error: file not found: $PROMPT_FILE. Please enter an existing file."
    elif [[ ! -r "$PROMPT_FILE" ]]; then
        echo "Error: no read permission: $PROMPT_FILE. Please give read permission to the file."
    else
        break
    fi
done

PARENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "

Starting agent execution with droid...

See the session logs in ~/.factory/sessions/ to monitor the agent execution and for debugging.
"

droid exec \
    --use-spec \
    --spec-model claude-opus-4-6 \
    --spec-reasoning-effort high \
    --model gemini-3-flash-preview \
    --reasoning-effort medium \
    --auto medium \
    --cwd "$PARENT_DIR" \
    --file "$PROMPT_FILE"

echo "

Agent execution completed.

Please manually verify the new resource by creating a terraform script file and
run it against a live environment.

Please review the acceptance tests in acceptance_tests/ and add any additional tests as needed.
"