#!/bin/bash
set -e

echo "$TRAVIS_REPO_SLUG":"$TAG"
# build the docker container
echo "Building Docker container"
make build-image

if [ $? -eq 0 ]; then
	echo "Complete"
else
	echo "Build Failed"
	exit 1
fi
