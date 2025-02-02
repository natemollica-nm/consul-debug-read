#!/usr/bin/env sh

set -e # Exit on error

ARCH=$(uname -m | grep -qE "aarch64|arm64" && echo "arm64" || echo "amd64")
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
VERSION=$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | awk -F'"' '/"tag_name"/ { print $4 }')
URL="https://github.com/natemollica-nm/consul-debug-read/releases/download/${VERSION}/consul-debug-read_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"

# Define colors
BOLD="\033[1m"; YELLOW="\033[93m"; RED="\033[91m"; BLUE="\033[96m"; CLEAR="\033[0m"

info() { echo "${BLUE}${BOLD}[INFO]${CLEAR} $*"; }
debug() { [ "$DEBUG" ] && echo "${YELLOW}${BOLD}[DEBUG]${CLEAR} $*"; }
warn() { echo "${YELLOW}${BOLD}[WARN]${CLEAR} $*" 1>&2; }
err() { echo "${RED}${BOLD}[ERROR]${CLEAR} $*" 1>&2; exit 1; }

# Prompt for sudo if required
require_sudo() {
  if [ "$(id -u)" -ne 0 ]; then
      info "Requesting sudo permissions..."
      sudo -v || err "Sudo required" && return 1
      return 0
  fi
  info "Running as root user..."
}

check_tool() {
    if ! command -v "$1" >/dev/null 2>&1; then
        warn "$1 not found. Install now? (y/n):"
        read -r resp && [ "$resp" = "y" ] && install_tool "$1" || err "Install $1 manually and re-run."
    else
        info "$1 detected: $(command -v "$1")"
    fi
}

install_tool() {
    case "$1" in
        "wget")
            info "Installing $1 using Homebrew..."
            command -v brew >/dev/null 2>&1 || /bin/sh -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
            brew install "$1" || err "Failed to install $1. Install manually and retry."
            ;;
        *) err "Unknown tool $1. Please install manually." ;;
    esac
}

cleanup() {
    debug "Cleaning up temporary files..."
    rm -f /tmp/consul-debug-read.tar.gz
    rm -rf /tmp/consul-debug-read
}
trap cleanup EXIT

install_latest_release() {
    info "Downloading consul-debug-read (${VERSION} | ${PLATFORM} | ${ARCH})..."
    require_sudo || {
      return 1
    }
    curl -fsSL "$URL" -o /tmp/consul-debug-read.tar.gz || wget -q "$URL" -O /tmp/consul-debug-read.tar.gz
    [ -f /tmp/consul-debug-read.tar.gz ] || err "Failed to download consul-debug-read."

    cd /tmp || exit
    tar -xzf consul-debug-read.tar.gz || err "Failed to extract consul-debug-read."
    [ -f consul-debug-read ] || err "Binary missing after extraction."

    sudo mv -f consul-debug-read /usr/local/bin/ || err "Failed to move binary to /usr/local/bin."
    chmod +x /usr/local/bin/consul-debug-read
    info "Installation complete: $(/usr/local/bin/consul-debug-read --version)"
}

setup_autocomplete() {
    # shellcheck disable=3044
    [ -f /usr/local/bin/consul-debug-read ] && complete -C /usr/local/bin/consul-debug-read consul-debug-read || true
}

main() {
    info "Checking prerequisites..."
    check_tool wget

    info "Installing consul-debug-read..."
    install_latest_release || { exit; }

    info "Setting up autocomplete support..."
    setup_autocomplete

    info "Consul-debug-read installed successfully!"
    echo
    echo
    info "Run 'consul-debug-read --help' to get started."
}

main "$@"