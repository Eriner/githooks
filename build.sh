#!/bin/bash

set -euf -o pipefail

go install ./...

echo "built all binaries!"
echo ""
echo "don't forget to add the hooks directory to your git config!"
echo "  $ git config --global core.hooksPath ${PWD}/hooks"
