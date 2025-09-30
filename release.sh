#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 0.2.0"
  exit 1
fi

VERSION="$1"
TAG="v$VERSION"

# Sanity check
if [[ ! "$VERSION" =~ ^[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Error: version must look like 1.2.3"
  exit 1
fi

# Make sure we're clean
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: you have uncommitted changes."
  exit 1
fi

# Create annotated tag
git tag -a "$TAG" -m "Release $TAG"

# Push tag to GitHub
git push origin "$TAG"

echo "✅ Created and pushed $TAG"
