#!/bin/bash

go install ./...

echo "built all binaries!"
echo ""
echo "don't forget to set/symlink your git hooks template dir!"
echo "  $ mkdir ~/.git-templates"
echo "  $ ln -s ${PWD}/hooks ~/.git-templates/hooks"
echo "  $ git config --global init.templatedir ~/.git-templates"
echo ""
echo "after you've created and configured the directory as above, you can add to existing repos with:"
echo "  $ usegithooks"
