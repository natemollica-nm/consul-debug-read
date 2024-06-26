#!/usr/bin/env bash

set -e

HASHICORP_REPO="${1:-consul}"
GH_PAT_TOKEN="${2:-"$GH_PAT"}"

if [ -z "$GH_PAT_TOKEN" ]; then
  echo "Null GH_TOKEN - Set your GH PAT to the environment var named \"GH_PATH\"!"
  exit 1
fi

GITHUB_API=https://api.github.com

curl -sSLf \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $GH_PAT_TOKEN" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "$GITHUB_API"/repos/hashicorp/"${HASHICORP_REPO}"/releases | jq -r '.[] | {Release: .name, Release_Date: .published_at} | to_entries[] | "\(.key): \(.value)"'