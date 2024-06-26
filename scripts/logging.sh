#!/usr/bin/env bash

export exit_code=0
trap 'cleanup' TERM ERR # Use trap to call cleanup when the script exits or errors out
cleanup() { exit "$exit_code"; }

## Logging Functions
now(){ date '+%d/%m/%Y-%H:%M:%S'; }
err() { >&2 printf '%s %b%s %s\e[0m\n' "$(now)" "${RED}[ERROR]${RESET} ${DIM}" "$@"; exit_code=1; }
warn() { >&2 printf '%s %b%s %s\e[0m\n' "$(now)" "${INTENSE_YELLOW}[WARN]${RESET} ${DIM}" "$@"; }
log() { printf '%s %b%s %s\e[0m\n' "$(now)" "${LIGHT_CYAN}[INFO]${RESET} ${DIM}" "$@"; }
prompt() { printf '%s %b%s %s\e[0m' "$(now)" "${INTENSE_YELLOW}[USER]${RESET} ${DIM}${BLINK}" "$@"; }

# Helper function to extract the value from the command line argument
# It handles both "--key value" and "--key=value" formats
extract_value() {
    local arg="$1"
    local next_arg="$2"

    if [[ "$arg" == *"="* ]]; then
        echo "${arg#*=}"  # Returns value after '='
    else
        echo "$next_arg"  # Returns the next argument
    fi
}

# Define the function to handle print messages with advanced formatting
print_msg() {
    local msg

    # shellcheck disable=2183
    printf '%s %b%s %s\e[0m\n' "$(now)" "${LIGHT_CYAN}[INFO]${RESET} ${DIM}" "$1"
    shift # Remove the initial message from the parameters

    # Loop through the remaining arguments to print additional messages
    for msg in "$@"; do
        # shellcheck disable=2183
        printf '%*s%b%s %s\e[0m\n' 27 '' "${LIGHT_GREEN}*==>${RESET} ${DIM}" "$msg"
    done
}