#!/usr/bin/env bash

set -eEo pipefail

ARCH=$( [[ "$(uname -m)" =~ aarch64|arm64 ]] && echo arm64 || echo amd64)
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
VERSION="$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | jq -r '.tag_name')"
URL=https://github.com/natemollica-nm/consul-debug-read/releases/download/"${VERSION}"/consul-debug-read_"${VERSION}"_"${PLATFORM}"_"${ARCH}".tar.gz

################ Font Formatting #################
##################################################
BOLD="\033[1m"; UNBOLD="\033[0m"; YELLOW="\e[93m";
RED="\e[91m"; BLUE="\e[96m"; END_COLOR="\e[0m"; INTENSE_YELLOW="\e[1;33m"
##################################################
DEBUG=${1:-0}
EXIT_CODE=0
# Return date/time at time of execution
function now { echo "$(date '+%d/%m/%Y-%H:%M:%S')"; }
# shellcheck disable=SC2145
function debug { test "$DEBUG" = 0 || { printf '\n%b%s' "$(now) ${YELLOW}""${BOLD}"[DEBUG]"${UNBOLD}""${END_COLOR} " " $@"; } }
# shellcheck disable=SC2145
function info { printf '\n%b%s' "$(now)  ${BLUE}""${BOLD}"[INFO]"${UNBOLD}""${END_COLOR} " " $@"; }
# shellcheck disable=SC2145
function warn { >&2 printf '\n%b%s' "$(now) ${INTENSE_YELLOW}""${BOLD}"[WARN]"${UNBOLD}""${END_COLOR} " " $@"; }
# shellcheck disable=SC2145
function err { >&2 printf '\n%b%s' "$(now) ${RED}""${BOLD}"[ERROR]"${UNBOLD}""${END_COLOR} " " $@"; EXIT_CODE=1 && echo && return "$EXIT_CODE"; }


# Define an associative array with tool names as keys and GitHub repository | software URLs as values
declare -A tools=(
    ["wget"]="https://formulae.brew.sh/formula/wget"
)

function sudo_prompt() {
  local attempt max_attempts auth
  attempt=0
  max_attempts=3

  while true; do
      # Prompt user for their sudo password with custom message
      info "Please enter sudo password for installation: "
      read -s -r auth && echo "$auth" | sudo -S -v -p ''
      if [ $? -eq 0 ]; then
          info "sudo authentication successful"
          # Place your script commands that require sudo here
          break
      else
          if [ $attempt -ge $max_attempts ]; then
              err "sudo authentication failed. Maximum attempts reached!"
              exit 1
          else
              warn "sudo authentication failed - please re-attempt sudo authentication password:"
              ((attempt++))
          fi
      fi
  done
}

function go_path_verification() {
  GOPATH="${GOPATH:-$HOME/go}"
  # Check if GOPATH is set
  if ! [ -d "$GOPATH" ]; then
    info "GOPATH directory not present, continuing"
    return 0 # Return as we cannot remove from unknown location
  fi

  # Construct the path to the binary
  local binaryPath="$GOPATH/bin/consul-debug-read" response

  # Check if the binary exists
  if [ -f "$binaryPath" ]; then
    # Prompt user for deletion
    info "previous binary found in \$GOPATH/bin"
    printf '\n\n%s\n' "    *==> Having '$GOPATH/bin/consul-debug-read' present can introduce conflicts when trying to run consul-debug-read."
    printf '%s' "    *==> Do you want to delete it? (y/n): "
    read -r response </dev/stdin

    case $response in
      [Yy]* )
        rm -f "$binaryPath" || { return 1; } # Delete the binary
        info "deleted consul-debug-read binary from $GOPATH/bin"
        return 0
        ;;
      * )
        info "skipping previous consul-debug-read binary deletion from $GOPATH/bin"
        return 0
        ;;
    esac
  else
    info "Previous binary at $GOPATH/bin not found, continuing"
  fi
  return 0
}

function check_and_install_brew() {
  # Check that brew is installed
  if ! command -v brew >/dev/null 2>&1; then
    # not installed
    read -p "Brew not detected. Install now? (y/n)" -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
    fi

    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
}
trap cleanup EXIT TERM SIGINT
function cleanup() {
  debug "install consul-debug-read: running cleaning up..."
  rm -f /tmp/consul-debug-read.tar.gz || {
   debug "install consul-debug-read: failed to remove /tmp/consul-debug-read.tar.gz"
  }
  rm -rf /tmp/consul-debug-read >/dev/null || {
   debug "install consul-debug-read: failed to remove /tmp/consul-debug-read"
  }
  debug "install consul-debug-read cleanup complete!"
  exit
}

# Function to install a tool
function install_tool() {
    local tool_name="$1"
    local tool_url="${tools[$tool_name]}"

    check_and_install_brew
    if [ -n "$tool_url" ]; then
        info "installing $tool_name..."
        if [[ "$tool_url" =~ .*brew.* ]]; then
          if ! brew install "$tool_name"; then
            err "failed to install $tool_name using homebrew | please install manually and reattempt scripted installation"
          fi
        else
          warn "missing $tool_name on local machine ($HOSTNAME) | please install manually and reattempt scripted installation"
          exit 
        fi
    else
        err "cannot install $tool_name. please install it manually."
    fi
}

# Function to check if a tool is installed
function check_tool_installed() {
    local tool_name="$1"

    if command -v "$tool_name" >/dev/null 2>&1; then
        info "$tool_name verified installed => $(which "$tool_name")"
    else
        warn "$tool_name not found on local machine ($HOSTNAME)."
        read -p "install $tool_name now? (y/n): " -n 1 -r
        if [[ "$REPLY" =~ ^[Yy]$ ]]; then
            install_tool "$tool_name"
            if command -v "$tool_name" >/dev/null 2>&1; then
              info "$tool_name verified installed => $(which "$tool_name")"
            else
              err "$tool_name installation failed. please install manually and reattempt scripted installation"
            fi
        fi
    fi
}

function install_autocomplete() {
  debug "running complete -C /usr/local/bin/consul-debug-read consul-debug-read"
  complete -C /usr/local/bin/consul-debug-read consul-debug-read || {
    warn "failed running autocompletion shell for consul-debug-read => 'complete -C /usr/local/bin/consul-debug-read consul-debug-read'"
  }
}

function install_latest_release() {
  info "downloading consul-debug-read (${VERSION} | ${PLATFORM} | ${ARCH})"
  rm -f /tmp/consul-debug-read.tar.gz || {
   err "install consul-debug-read: failed to remove /tmp/consul-debug-read.tar.gz"
  }
  ! test -f /usr/local/bin/consul-debug-read || {
    debug "install consul-debug-read: removing file named consul-debug-read from /usr/local/bin"
    sudo rm -rf /usr/local/bin/consul-debug-read >/dev/null 2>&1 || true;
  }
  go_path_verification || {
    err "failed to remove $GOPATH/bin/consul-debug-read"
  }
  ! test -d /usr/local/bin/consul-debug-read || {
    debug "install consul-debug-read: removing dir named consul-debug-read from /usr/local/bin"
    sudo rm -rf /usr/local/bin/consul-debug-read >/dev/null 2>&1 || true;
  }
  debug "install consul-debug-read: pull binary from ${URL}"
  /bin/bash -c "$(curl -fsSL "${URL}" -o /tmp/consul-debug-read.tar.gz)" || {
    warn "install consul-debug-read: curl attempt to download failed, switching to wget"
    wget -q --show-progress --tries=3 --timeout=10 --retry-connrefused "${URL}" -O /tmp/consul-debug-read.tar.gz >/dev/null 2>&1 || {
      err "install consul-debug-read: failed to download binary from ${URL}"
    }
  }
  debug "install consul-debug-read: running tarball extraction from /tmp"
  tar -xf /tmp/consul-debug-read.tar.gz -C /tmp 2>&1 || {
    err "install consul-debug-read: failed to untar /tmp/consul-debug-read.tar.gz"
  }
  debug "install consul-debug-read: confirming extraction"
  test -f /tmp/consul-debug-read || {
    err "install consul-debug-read: consul-debug-read binary not found at /tmp/consul-debug-read"
  }
  debug "install consul-debug-read: moving consul-debug-read to /usr/local/bin"
  sudo mv -f /tmp/consul-debug-read /usr/local/bin/consul-debug-read 2>&1 || {
    err "install consul-debug-read: failed to move binary to /usr/local/bin/"
  }
  debug "install consul-debug-read: running 'command -v consul-debug-read'"
  command -v consul-debug-read >/dev/null 2>&1 || {
    err "install consul-debug-read: installation failed! exiting."
  }
}
clear
info "starting scripted installation of consul-debug-read cli tool"
info "verifying prerequisites installed for downloading consul-debug-read"
# Loop through the tools and check if they are installed
for tool in "${!tools[@]}"; do
    check_tool_installed "$tool"
done
info "prerequisite verification complete, continuing with installation"
sudo_prompt && install_latest_release
info "consul-debug-read version: installed successfully!"
printf '\n\n%s\n%s\n%s\n%s\n%s\n' \
  "       consul-debug-read: debug bundle cli-tool" \
  "    ===============================================" \
  "    version:            $(/usr/local/bin/consul-debug-read --version)" \
  "    install location:   /usr/local/bin/consul-debug-read" \
  "    autocomplete cmd:   complete -C /usr/local/bin/consul-debug-read consul-debug-read"