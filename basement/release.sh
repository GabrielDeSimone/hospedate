#!/usr/bin/env bash

# Stop if a command fails
set -e

if [ $# -ne 2 ]; then
  echo "Please provide <yard-version> and <backyard-version>"
  exit 1
fi

YARD_VERSION=$1;
BACKYARD_VERSION=$2;

git fetch
if [ $(git rev-parse HEAD) != $(git rev-parse @{u}) ]; then
  echo "Please update local basement branch before continuing"
  exit 1
fi

PRODUCT_VERSION=$(grep -n "." VERSIONS | tail -n 1 | awk '{print $6}')
middle_number=$(echo $PRODUCT_VERSION | awk -F. '{print $2}')
middle_number=$((middle_number + 1))
NEW_PRODUCT_VERSION=$(echo $PRODUCT_VERSION | awk -F. '{print $1 "." '"$middle_number"' "." $3}')
echo "Releasing product version number: $NEW_PRODUCT_VERSION (last one was $PRODUCT_VERSION)"

printf " $YARD_VERSION | $BACKYARD_VERSION | $NEW_PRODUCT_VERSION\n" >> VERSIONS

echo "Saving changes..."
git add VERSIONS
git commit -m "Release $NEW_PRODUCT_VERSION"
git push

echo "New product version $NEW_PRODUCT_VERSION released."
