#!/usr/bin/env bash

set -e # Exit immediately on error

export GITHUB_TOKEN=$CONSUL_DEBUG_GH_TOKEN

VERSION_FILE="internal/read/version.go"
MAIN_BRANCH="main"

# Check if version file exists
if [ ! -f "$VERSION_FILE" ]; then
  echo "Version file not found: $VERSION_FILE"
  exit 1
fi

# Step 1: Detect version from internal/read/version.go
CURRENT_VERSION=$(ggrep -oP 'Version\s*=\s*"\K[^"]+' "$VERSION_FILE")
if [ -z "$CURRENT_VERSION" ]; then
  echo "Unable to extract version from $VERSION_FILE"
  exit 1
fi

echo "Current version detected: $CURRENT_VERSION"

# Step 2: Check for the latest release on GitHub
LATEST_RELEASE=$(gh release list --limit 1 | awk '{ print $1 }')

if [ -z "$LATEST_RELEASE" ]; then
  echo "No GitHub releases found."
else
  echo "Latest release on GitHub: $LATEST_RELEASE"
fi

# Compare current version with the latest release
if [ "$CURRENT_VERSION" == "$LATEST_RELEASE" ]; then
  echo "The current version is already the latest release."
else
  echo "The current version is different from the latest release."
fi

# Step 3: Prompt user to confirm or update version
read -rp "Do you want to update the version? (y/n): " UPDATE_VERSION
if [[ "$UPDATE_VERSION" =~ ^[Yy]$ ]]; then
  read -rp "Enter new version (current: $CURRENT_VERSION): " NEW_VERSION
  if [ -z "$NEW_VERSION" ]; then
    echo "Invalid input. Exiting."
    exit 1
  fi

  echo "Updating version to: $NEW_VERSION"

  # Update the version in the version.go file
  sed -i.bak "s/Version\s*=\s*\"$CURRENT_VERSION\"/Version = \"$NEW_VERSION\"/" "$VERSION_FILE"

  # Confirm the change
  echo "Version updated to: $NEW_VERSION"
  grep 'Version' "$VERSION_FILE"
  read -rp "Commit and proceed with $NEW_VERSION? (y/n): " CONFIRM
  if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "Aborting."
    mv "$VERSION_FILE.bak" "$VERSION_FILE" # Restore backup
    exit 1
  fi

  # Clean up backup file
  rm "$VERSION_FILE.bak"

else
  echo "Keeping current version: $CURRENT_VERSION"
  NEW_VERSION=$CURRENT_VERSION
fi

# Step 4: Git add, commit, and push changes to main
git checkout "$MAIN_BRANCH"
git pull origin "$MAIN_BRANCH"

echo "Adding and committing version update to Git..."
git add "$VERSION_FILE"
git commit -m "Update to version $NEW_VERSION"
git push origin "$MAIN_BRANCH"

# Step 5: Create a git tag and push
echo "Creating git tag v$NEW_VERSION..."
git tag -a "v$NEW_VERSION" -m "Release version v$NEW_VERSION"
git push origin "v$NEW_VERSION"

# Step 6: Run goreleaser
echo "Running goreleaser release..."
if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN is not set. Please set it and retry: export GITHUB_TOKEN=\$CONSUL_DEBUG_GH_TOKEN"
  exit 1
fi

goreleaser release --clean
if [ $? -eq 0 ]; then
  echo "Release completed successfully!"
else
  echo "Release failed."
  exit 1
fi