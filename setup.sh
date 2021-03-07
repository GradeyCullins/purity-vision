#!/usr/bin/env bash

# Set the git hooks directory
if ! git config core.hooksPath .githooks; then
    echo "Failed to set githooks directory"
fi
