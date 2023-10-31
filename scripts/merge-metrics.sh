#!/usr/bin/env bash

set -e

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <debugPath>"
    exit 1
fi

debugPath="$1"

# Iterate through the items in the debugPath directory
for captureDir in "$debugPath"/*; do
    if [ -d "$captureDir" ]; then
        metrics_data="$(cat "$captureDir/metrics.json")"
        printf '%s\n' "$metrics_data" >> "$debugPath/metrics.json"
    fi
done
