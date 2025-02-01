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

echo "Current version detected: v$CURRENT_VERSION"

# Step 2: Check for the latest release on GitHub
LATEST_RELEASE=$(gh release list --limit 1 | awk '{ print $1 }')

if [ -z "$LATEST_RELEASE" ]; then
  echo "No GitHub releases found."
else
  echo "Latest release on GitHub: $LATEST_RELEASE"
fi

# Compare current version with the latest release
if [ "v$CURRENT_VERSION" == "$LATEST_RELEASE" ]; then
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

# Switch to the main branch and pull the latest changes
if ! git branch --show-current | grep -q "$MAIN_BRANCH"; then
    if git checkout "$MAIN_BRANCH"; then
      echo "Switched to branch '$MAIN_BRANCH' successfully."
    else
      echo "Failed to switch to branch '$MAIN_BRANCH'"
      exit 1
    fi
fi

# Add and commit version update only if there are changes
if ! git diff --quiet "$VERSION_FILE"; then
    echo "Adding and committing version update to Git..."
    git add "$VERSION_FILE"
    git commit -m "Update to version $NEW_VERSION"
    git push origin "$MAIN_BRANCH"
else
    echo "No changes detected in $VERSION_FILE. Skipping commit."
fi

# Step 5: Create a git tag and push
echo "Creating git tag v$NEW_VERSION..."
if git tag -l | grep -q "v$NEW_VERSION"; then
    read -rp "v$NEW_VERSION tagged branch already present. Delete and recreate? (y/n): " RECREATE_TAGGED_BR
    if [[ "${RECREATE_TAGGED_BR}" =~ ^[Yy]$ ]]; then
        git tag -d "v${NEW_VERSION}"
        git push --delete origin "v${NEW_VERSION}"
        gh release delete "v${NEW_VERSION}" --yes
    else
        echo "Fix tagged branch conflict and rerun $(basename "${0}"), exiting..."
        exit 1
    fi
fi
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