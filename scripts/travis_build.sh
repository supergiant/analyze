#!/bin/bash
# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

#replace all slashes from string "task/S20-950" => "taskS20-950"
TRAVIS_BRANCH=${TRAVIS_BRANCH/\//}

export TAG=${TRAVIS_BRANCH:-unstable}


### Main
echo "Tag Name: ${TAG}"
# If a tag is pushed, tests are run, the docker container is built and pushed to
# dockerhub, and then a release is pushed to the releases page.
if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
	echo "release"
	# run linters
	make lint
	# run tests
	make test
	# Build Docker container
	make build
	# Push to Dockerhub
	make push
	# Push to releases page
	make push-release
# on an unstable branch, tests are run and the docker container is built and pushed.
elif [[ "$TRAVIS_BRANCH" == *release-* ]]; then
	echo "unstable branch"
	export TAG="${TAG}-unstable"
	echo "Tag Name: ${TAG}"
	# run linters
	make lint
	# run tests
	make test
	# Build docker container
	make build
	# Push to Dockerhub
	make docker-push
# if a push to master happens, tests are only run
elif [[ "$TRAVIS_BRANCH" == "master" ]]; then
	echo "master branch - test will only be run"
	echo "Tag Name: ${TAG}"
	# run linters
	make lint
	# run tests
	make test
	# Build docker container
	make build
else
# any other branch is considered a testing branch and will only run tests and build the container.
	echo "testing branch - run tests and docker build"
	export TAG="${TAG}-testing"
	echo "Tag Name: ${TAG}"
	# run linters
	make lint
	# run tests
	make test
	# Build docker container
	make build
fi