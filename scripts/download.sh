#!/usr/bin/env bash
set -eEo pipefail

# Fail-safe: Ensure compatible bash version.
required_bash_version="4.0"
if (( BASH_VERSINFO[0] < ${required_bash_version%%.*} || (BASH_VERSINFO[0] == ${required_bash_version%%.*} && BASH_VERSINFO[1] < ${required_bash_version#*.}) )); then
  echo -e "\033[1;31mERROR:\033[0m Bash version ${required_bash_version} or higher is required to run this script. Detected version: ${BASH_VERSION}."
  echo "Please upgrade your bash installation to proceed."
  exit 1
fi

# Check for sudo privileges or re-execute with sudo
ensure_sudo() {
  if [[ $EUID -ne 0 ]]; then
    echo -e "\033[1;33mINFO:\033[0m This script requires elevated privileges. Prompting for sudo..."
    sudo bash "$0" "$@"
    exit 0
  fi
}

# Constants
readonly ARCH=$( [[ "$(uname -m)" =~ aarch64|arm64 ]] && echo arm64 || echo amd64 )
readonly PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
readonly URL_BASE="https://github.com/natemollica-nm/consul-debug-read/releases/download"
readonly VERSION="$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | jq -r '.tag_name')"
readonly DOWNLOAD_URL="${URL_BASE}/${VERSION}/consul-debug-read_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"
readonly INSTALL_DIR="/usr/local/bin"  # Default install directory
readonly DEBUG=${1:-0}

# Font Formatting for Logs
readonly COLOR_BOLD="\033[1m"
readonly COLOR_RESET="\033[0m"
readonly COLOR_YELLOW="\e[93m"
readonly COLOR_RED="\e[91m"
readonly COLOR_BLUE="\e[96m"

# Unified log function
log_message() {
  local level="$1"
  local color="$2"
  local message="$3"
  printf "\n%b[%s]%b %s\n" "${color}${COLOR_BOLD}" "${level}" "${COLOR_RESET}" "${message}"
}
info() { log_message "INFO" "$COLOR_BLUE" "$1"; }
debug() { [[ "$DEBUG" -ne 0 ]] && log_message "DEBUG" "$COLOR_YELLOW" "$1"; }
error() { log_message "ERROR" "$COLOR_RED" "$1" >&2; exit 1; }

# GOPATH Verification
verify_gopath_binary() {
  local gopath="${GOPATH:-$HOME/go}"
  local binary_path="$gopath/bin/consul-debug-read"

  if [[ -f "$binary_path" ]]; then
    info "Existing binary detected at $binary_path."
    read -p "Delete it? (y/n): " -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
      rm -f "$binary_path" || err "Failed to delete $binary_path"
      info "Binary removed."
    else
      info "Binary deletion skipped."
    fi
  else
    info "No existing binary found in $gopath/bin."
  fi
}

# Cleanup temporary files
cleanup() {
  debug "Cleaning up temporary files..."
  rm -f /tmp/consul-debug-read.tar.gz
  rm -rf /tmp/consul-debug-read
}
trap cleanup EXIT TERM SIGINT

# Download and install the latest release
install_latest_release() {
  info "Downloading consul-debug-read (${VERSION} for ${PLATFORM}/${ARCH})..."
  curl -sL -o /tmp/consul-debug-read.tar.gz "$DOWNLOAD_URL" || error "Download failed. Please check your connection and try again."

  info "Extracting and installing release..."
  tar -xzf /tmp/consul-debug-read.tar.gz -C /tmp || error "Extraction failed."

  local extracted_binary="/tmp/consul-debug-read"
  if [[ -f "$extracted_binary" ]]; then
    mv -f "$extracted_binary" "$INSTALL_DIR/consul-debug-read" || error "Failed to move binary to $INSTALL_DIR. Check your permissions."
    chmod +x "$INSTALL_DIR/consul-debug-read"
    info "Installation completed successfully! consul-debug-read (v$(consul-debug-read --version)) is now available in $(which consul-debug-read)."
  else
    error "Downloaded release does not include the expected binary. Please verify the release and try again."
  fi
}

# Main installation process
main() {
  verify_gopath_binary
  install_latest_release
}

# Run the sudo check and main script
ensure_sudo "$@"
main "$@"