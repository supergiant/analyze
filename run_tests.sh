#!/bin/bash
set -e

echo "Running tests"

make test

# Check for errors
if [ $? -eq 0 ]; then
	echo "Tests Passed"
else
	echo "Tests Failed"
	exit 1
fi
