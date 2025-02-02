#!/bin/sh
set -e

# Constants
ARCH=$(uname -m | grep -qiE "aarch64|arm64" && echo "arm64" || echo "amd64")
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
URL_BASE="https://github.com/natemollica-nm/consul-debug-read/releases/download"
VERSION=$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | awk -F'"' '/"tag_name"/ { print $4 }')
DOWNLOAD_URL="$URL_BASE/$VERSION/consul-debug-read_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"
INSTALL_DIR="/usr/local/bin"  # Default install directory
DEBUG="${1:-0}"

# Font Formatting for Logs
COLOR_BOLD="\033[1m"
COLOR_RESET="\033[0m"
COLOR_YELLOW="\033[93m"
COLOR_RED="\033[91m"
COLOR_BLUE="\033[96m"

# Unified log function
log_message() {
  level="$1"
  color="$2"
  message="$3"
  printf "\n%s[%s]%s %s\n" "${color}${COLOR_BOLD}" "$level" "${COLOR_RESET}" "$message"
}
info() { log_message "INFO" "$COLOR_BLUE" "$1"; }
debug() { [ "$DEBUG" -ne 0 ] && log_message "DEBUG" "$COLOR_YELLOW" "$1"; }
error() { log_message "ERROR" "$COLOR_RED" "$1"; exit 1; }

# Cleanup temporary files
trap cleanup EXIT TERM INT
cleanup() {
  debug "Cleaning up temporary files..."
  rm -f /tmp/consul-debug-read.tar.gz
  rm -rf /tmp/consul-debug-read
}

# GOPATH Verification
verify_gopath_binary() {
  gopath="${GOPATH:-$HOME/go}"
  binary_path="$gopath/bin/consul-debug-read"

  if [ -f "$binary_path" ]; then
    info "Existing binary detected at $binary_path."
    printf "Delete it? (y/n): "
    read -r response
    if echo "$response" | grep -qi "^y"; then
      rm -f "$binary_path" || error "Failed to delete $binary_path"
      info "Binary removed."
    else
      info "Binary deletion skipped."
    fi
  else
    info "No existing binary found in $gopath/bin."
  fi
}

# Cleanup temporary files (defined above as well for trap)
cleanup() {
  debug "Cleaning up temporary files..."
  rm -f /tmp/consul-debug-read.tar.gz
  rm -rf /tmp/consul-debug-read
}

# Download and install the latest release
install_latest_release() {
  info "Downloading consul-debug-read ($VERSION for $PLATFORM/$ARCH)..."
  curl -sL -o /tmp/consul-debug-read.tar.gz "$DOWNLOAD_URL" || error "Download failed. Please check your connection and try again."

  info "Extracting and installing release..."
  tar -xzf /tmp/consul-debug-read.tar.gz -C /tmp || error "Extraction failed."

  extracted_binary="/tmp/consul-debug-read"
  if [ -f "$extracted_binary" ]; then
    mv -f "$extracted_binary" "$INSTALL_DIR/consul-debug-read" || error "Failed to move binary to $INSTALL_DIR. Check your permissions."
    chmod +x "$INSTALL_DIR/consul-debug-read"
    info "Installation completed successfully! consul-debug-read (v$("$INSTALL_DIR/consul-debug-read" --version)) is now available in $(command -v consul-debug-read)."
  else
    error "Downloaded release does not include the expected binary. Please verify the release and try again."
  fi
}

# Main installation process
main() {
  verify_gopath_binary
  install_latest_release
}

# Call main script
main "$@"