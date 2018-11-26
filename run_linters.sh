#!/bin/bash
set -e

echo "Running linters"

make lint

# Check for errors
if [ $? -eq 0 ]; then
	echo "Linting Passed"
else
	echo "Linting Failed"
	exit 1
fi
