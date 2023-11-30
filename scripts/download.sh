#!/usr/bin/env bash

set -e

ARCH=$( [[ "$(uname -m)" =~ aarch64|arm64 ]] && echo arm64 || echo amd64)
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
VERSION="$(curl -s https://api.github.com/repos/natemollica-nm/consul-debug-read/releases/latest | jq -r '.tag_name')"
URL=https://github.com/natemollica-nm/consul-debug-read/releases/download/"${VERSION}"/consul-debug-read_"${VERSION}"_"${PLATFORM}"_"${ARCH}".tar.gz

trap cleanup EXIT TERM SIGINT
function cleanup() {
  echo "cleaning up..."
  rm -f ./consul-debug-read.tar.gz >/dev/null 2>&1 || true
  rm -rf consul-debug-read >/dev/null 2>&1 || true
  echo "done"
  exit
}

function install_latest_release() {
  echo "downloading consul-debug-read (${VERSION} | ${PLATFORM} | ${ARCH})"
  rm -f ./consul-debug-read.tar.gz &>/dev/null || true
  ! test -f /usr/local/bin/consul-debug-read || { sudo rm -rf /usr/local/bin/consul-debug-read >/dev/null 2>&1 || true; }
  wget -q --show-progress --tries=3 --timeout=10 --retry-connrefused "${URL}" -O consul-debug-read.tar.gz >/dev/null 2>&1 || true
  tar -xf consul-debug-read.tar.gz >/dev/null 2>&1 || true
  sudo mv -f consul-debug-read /usr/local/bin/consul-debug-read
  command -v consul-debug-read >/dev/null 2>&1 || {
    echo "consul-debug-read installation failed! exiting."
    exit 1
  }
}

clear
install_latest_release
echo "consul-debug-read ${VERSION} installed => /usr/local/bin/consul-debug-read"