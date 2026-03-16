#!/bin/bash

# this runs droid in headless mode.
#  it sources the .env file to inject the env var for the factory API key:
#
#  export FACTORY_API_KEY=fk-...
#
#  the reason droid does not source the .env file itself is that it might read the file and put the API key in the context window.
#  a future release may allow excluding files from droid
#
#  this will run droid in spec mode using claude opus 4.6, and implement the code using gemini 3 flash.
#  feel free to change these to your preferred models.
#  stream-json lets you see the agents execution in real time.
#  it reads a prompt from ~/prompt.  change this to point to your prompt file.

set -euo pipefail

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
    --file ~/prompt | tee agent.log