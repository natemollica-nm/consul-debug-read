#!/usr/bin/env bash

export TEARDOWN="${1:-false}"

eval "$(cat scripts/formatting.env)"

SCRIPT_NAME=$(basename "$0")

exit_code=0
err() { >&2 printf '%s %b%s %s\e[0m\n' "$(now)" "   ${RED}[ERROR]${RESET} ${LIGHT_YELLOW}${SCRIPT_NAME}${RESET}: ${DIM}" "$@"; exit_code=1; }
# Use trap to call cleanup when the script exits or errors out
trap 'cleanup' EXIT TERM ERR
# Cleanup resources
cleanup() { exit "$exit_code"; }
now(){ date '+%d/%m/%Y-%H:%M:%S'; }
warn() { >&2 printf '%s %b%s %s\e[0m\n' "$(now)" "    ${INTENSE_YELLOW}[WARN]${RESET} ${LIGHT_YELLOW}${SCRIPT_NAME}${RESET}: ${DIM}" "$@"; }
log() { printf '%s %b%s %s\e[0m\n' "$(now)" "    ${LIGHT_CYAN}[INFO]${RESET} ${LIGHT_YELLOW}${SCRIPT_NAME}${RESET}: ${DIM}" "$@"; }

if [ "$TEARDOWN" != true ]; then
  log "aws credentials invalid, refreshing doormat aws cli credentials ..."
  # Generate new credentials using your CLI tool and capture the output
  NEW_CREDS="$(doormat aws export --account "$(doormat aws list | tail -n1 | awk '{print $1}' | awk -F '/' '{print $2}' | sed -e 's/-developer//g')")"
  if [ -z "$NEW_CREDS" ]; then
    err "null credentials returned, failed to export via doormat-cli..."
    exit 1
  fi
else
  NEW_CREDS="export AWS_ACCESS_KEY_ID=null && export AWS_SECRET_ACCESS_KEY=null && export AWS_SESSION_TOKEN=null && export AWS_SESSION_EXPIRATION=null"
fi

# Backup the existing .env file
log "backing up .env => .env.backup"
cp .env .env.backup
echo -e "\n" >> .env
# Split the output into individual export commands and process each
(
  IFS='&&' read -ra CMD_ARRAY <<<"$NEW_CREDS"
  for cmd in "${CMD_ARRAY[@]}"; do
    # Trim leading and trailing whitespace
    cmd=$(echo "$cmd" | xargs)

    # Extract the variable name and value
    if [[ "$cmd" =~ ^export\ ([A-Z_]+)=(.*)$ ]]; then
        VAR_NAME=${BASH_REMATCH[1]}
        VAR_VALUE=${BASH_REMATCH[2]}

        # Escape & for use in sed replacement string
        # shellcheck disable=2001
        VAR_VALUE_ESCAPED=$(echo "$VAR_VALUE" | sed 's/[&]/\\&/g')

        # Check if the variable already exists
        if grep -q "^export $VAR_NAME=" .env && [ "$TEARDOWN" != true ]; then
            log "updating $VAR_NAME"
            # Variable exists, use sed to replace its value
            sed -i.bak "s|^export $VAR_NAME=.*|export $VAR_NAME=\"$VAR_VALUE_ESCAPED\"|" .env
        elif [ "$TEARDOWN" = true ]; then
            log "removing $VAR_NAME"
            # Variable exists, use sed to replace its value
            sed -i.bak "s|^export $VAR_NAME=.*||" .env
        else
            log "adding $VAR_NAME"
            # Variable doesn't exist, append it to the file
            echo "export $VAR_NAME=\"$VAR_VALUE_ESCAPED\"" >> .env
        fi
        rm -f .env.bak  # Remove the backup file created by sed, if any
    fi
  done
)
if [ "$TEARDOWN" != true ]; then
  log "credentials updated in .env"
  eval "$(cat .env)"
  # Reload the .env file into the current session
  if aws sts get-caller-identity --no-cli-pager >/dev/null 2>&1; then
    log "successfully updated doormat aws cli credentials => .env"
    rm .env.backup # Remove the manual backup
  else
    err "failed to update credentials (aws sts call failed...)"
    exit 1
  fi
else
  rm .env.backup # Remove the manual backup
  sed -i '' -e :a -e '/^\n*$/{$d;N;ba' -e '}' .env ## Trim the extra newlines from the end of the file
  log "credentials removed from .env"
fi
log "complete!"