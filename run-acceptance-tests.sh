#!/bin/bash
set -e

# Load variables from terraform.tfvars
if [ ! -f terraform.tfvars ]; then
  echo "terraform.tfvars not found! Please create it in the project root."
  exit 1
fi

# Export variables from terraform.tfvars (robust to spaces and quotes)
while IFS='=' read -r key value; do
  # Remove whitespace and comments
  key=$(echo "$key" | xargs)
  value=$(echo "$value" | sed 's/["\r\n]*$//;s/^[ \t]*//;s/[ \t]*$//')
  # Remove inline comments
  value=$(echo "$value" | cut -d'#' -f1 | xargs)
  # Remove surrounding quotes
  value=$(echo "$value" | sed 's/^"//;s/"$//')
  if [[ -n "$key" && -n "$value" ]]; then
    export TF_VAR_${key}="${value}"
  fi
done < <(grep -v '^#' terraform.tfvars | grep '=')

echo "Exported variables:"
env | grep TF_VAR_

echo "Running acceptance tests..."
make testacc
