#!/usr/bin/env bash

set -eEo pipefail

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

ARCH=$( [[ "$(uname -m)" =~ aarch64|arm64 ]] && echo arm64 || echo amd64)
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
VERSION="$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | jq -r '.tag_name')"
URL=https://github.com/natemollica-nm/consul-debug-read/releases/download/"${VERSION}"/consul-debug-read_"${VERSION}"_"${PLATFORM}"_"${ARCH}".tar.gz

trap cleanup EXIT TERM SIGINT
function cleanup() {
  debug "install consul-debug-read: running cleaning up..."
  rm -f /tmp/consul-debug-read.tar.gz || {
   debug "install consul-debug-read: failed to remove /tmp/consul-debug-read.tar.gz"
  }
  rm -rf /tmp/consul-debug-read >/dev/null || {
   debug "install consul-debug-read: failed to remove /tmp/consul-debug-read"
  }
  info "install consul-debug-read: done"
  exit
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
  ! test -d /usr/local/bin/consul-debug-read || {
    debug "install consul-debug-read: removing dir named consul-debug-read from /usr/local/bin"
    sudo rm -rf /usr/local/bin/consul-debug-read >/dev/null 2>&1 || true;
  }
  debug "install consul-debug-read: pull binary from ${URL}"
  wget -q --show-progress --tries=3 --timeout=10 --retry-connrefused "${URL}" -O /tmp/consul-debug-read.tar.gz >/dev/null 2>&1 || {
    err "install consul-debug-read: failed to download binary from ${URL}"
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
install_latest_release
info "consul-debug-read ${VERSION} installed => /usr/local/bin/consul-debug-read"