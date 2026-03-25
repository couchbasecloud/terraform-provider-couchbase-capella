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
#  it will prompt you for one or more paths to prompt files.  if multiple prompt files are provided,
#  droids will run in parallel, each in its own git worktree.  you can view active worktrees with:
#
#    git worktree list
#
#  you can also review the session logs in ~/.factory/sessions/ to see the agent execution and for debugging.
#  the session log is a jsonl file.

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


echo "
You will be prompted for one or more paths to prompt files.
If multiple prompt files are provided, droids will run in parallel, each in its own git worktree.
View active worktrees with: git worktree list
"

PROMPT_FILES=()

while true; do
    while true; do
        read -rp "Enter the path to the prompt file (or 'cancel' to proceed): " PROMPT_FILE
        PROMPT_FILE="${PROMPT_FILE#"${PROMPT_FILE%%[![:space:]]*}"}"
        PROMPT_FILE="${PROMPT_FILE%"${PROMPT_FILE##*[![:space:]]}"}"
        if [[ "$PROMPT_FILE" == "cancel" ]]; then
            break
        fi
        # Expand ~ to $HOME since tilde expansion doesn't happen in variables
        PROMPT_FILE="${PROMPT_FILE/#\~/$HOME}"
        if [[ -z "$PROMPT_FILE" ]]; then
            echo "Error: no path entered. Please enter a path."
        elif [[ ! -f "$PROMPT_FILE" ]]; then
            echo "Error: file not found: $PROMPT_FILE. Please enter an existing file."
        elif [[ ! -r "$PROMPT_FILE" ]]; then
            echo "Error: no read permission: $PROMPT_FILE. Please give read permission to the file."
        else
            # Check for duplicate
            DUPLICATE=false
            for existing in "${PROMPT_FILES[@]+"${PROMPT_FILES[@]}"}"; do
                if [[ "$existing" == "$PROMPT_FILE" ]]; then
                    DUPLICATE=true
                    break
                fi
            done
            if [[ "$DUPLICATE" == true ]]; then
                echo "Warning: $PROMPT_FILE has already been added. Skipping."
            else
                PROMPT_FILES+=("$PROMPT_FILE")
                echo "Added: $PROMPT_FILE"
            fi
            break
        fi
    done

    # If user typed 'cancel', stop collecting
    if [[ "$PROMPT_FILE" == "cancel" ]]; then
        break
    fi

    read -rp "Add another prompt file? (y/n): " ADD_MORE
    if [[ "$ADD_MORE" != "y" && "$ADD_MORE" != "Y" ]]; then
        break
    fi
done

if [[ ${#PROMPT_FILES[@]} -eq 0 ]]; then
    echo "Warning: no prompt files were added. Exiting."
    exit 0
fi

PARENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "

Starting agent execution for ${#PROMPT_FILES[@]} prompt file(s).

See the session logs in ~/.factory/sessions/ to monitor the agent execution and for debugging.
The session log is a jsonl file.
"

PIDS=()
for PROMPT_FILE in "${PROMPT_FILES[@]}"; do
    echo "Launching droid for: $PROMPT_FILE"
    droid exec \
        --use-spec \
        --spec-model claude-opus-4-6 \
        --spec-reasoning-effort high \
        --model gemini-3-flash-preview \
        --reasoning-effort medium \
        --auto medium \
        --cwd "$PARENT_DIR" \
        --file "$PROMPT_FILE" > /dev/null &
    PIDS+=($!)
    sleep 3
done

echo "
Waiting for ${#PIDS[@]} droid job(s) to finish...
"

for i in "${!PIDS[@]}"; do
    if ! wait "${PIDS[$i]}"; then
        echo "Error: droid failed for: ${PROMPT_FILES[$i]} (PID ${PIDS[$i]})"
    fi
done


echo "

All agent executions completed.

Please manually verify the new resources by creating terraform script files and run them against a live environment.

Please review and run the acceptance tests in acceptance_tests/, and add any additional tests.
"