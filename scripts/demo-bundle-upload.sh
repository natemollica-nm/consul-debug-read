#!/usr/bin/env bash

set -e

# Configuration
REMOTE_USER="ubuntu"
REMOTE_HOST="consul-support-demo.nathan-mollica.sbx.hashidemos.io"
REMOTE_PATH="/var/www/html/debug_bundles/"
SSH_KEY_PATH=".offsite/ssh/aws-nate-ohio.pem"

# Function to upload file using rsync
upload_file() {
    local file="$1"
    local remote_dir="${2#/}"

    echo "Uploading '$file' to $REMOTE_HOST *==> $remote_dir/"
    rsync -avz -e "ssh -i $SSH_KEY_PATH -o ServerAliveInterval=60 -o ServerAliveCountMax=30 -o Compression=yes" "$file" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH/$remote_dir/"
}

# Check if the argument is passed
if [ -z "$1" ] || [ -z "$2" ] || [[ "$1" =~ -h|--help ]]; then
    echo "Usage: $0 <file-or-directory-path> <remote-path>"
    exit 1
fi

# Get the absolute path of the provided argument
INPUT_PATH=$(readlink -f "$1")
REMOTE_DIR="$2"

# Check if the provided path is a directory
if [ -d "$INPUT_PATH" ]; then
    # Find all .tar.gz files in the directory and upload each
    find "$INPUT_PATH" -type f -name "*.tar.gz" | while read -r file; do
        echo "Uploading $file..."
        upload_file "$file" "$REMOTE_DIR"
    done
elif [ -f "$INPUT_PATH" ] && [[ "$INPUT_PATH" == *.tar.gz ]]; then
    # If it's a single .tar.gz file, upload it
    echo "Uploading $INPUT_PATH..."
    upload_file "$INPUT_PATH" "$REMOTE_DIR"
else
    echo "The provided path is not a .tar.gz file or directory containing .tar.gz files."
    exit 1
fi

echo "Upload completed."
