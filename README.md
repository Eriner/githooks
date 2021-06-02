# githooks

Just some git hooks I use.

## Purpose

These tools exist almost entirely to mitigate human error. I make mistakes, and
these automations are safety-switches that future-me will thank present-me for.

There are occasions when I'll use `$ git --no-verify`, but I would like that to
be the exception.

Ideally when these tools declare a failure, it means things are **broken and need
to be fixed** - not that I should just override the verifications.

## Use

To build all tools, `$ ./build.sh`. (it just runs `go install ./...`)
The build script will print instructions for setting up your git templates dir.
From that point forward, new git repositories will automatically use these hooks.

To add to an existing project, run `$ usegithooks`, which is installed by `build.sh`.

The [hook](./hooks) scripts merely call the tools in this repository in sequence.

All tools are intended to be "smart" enough to know when they are and are not
relevant. They're all written in Go, and currently without any external dependencies.
These hooks can be installed easily on an air-gapped machine.

When invoked (when a tool deems itself "relevant" to the current commit) the tool
will either exit with an error, halting the git operation, or will print an
"all clear".

## Tools

Some of these tools may be broken out into separate sub-tools in the future,
ex: `hook-go` becomes `hook-go`, `hook-go-security`, `hook-go-test`, `hook-go-mod`,
etc.

### pre-commit

* [hook-interactive-secrets](./cmd/hook-interactive-secrets) checks for secrets,
and provides interactive override.
* [hook-go](./cmd/hook-go) lint, vet, and test Go code.
* [hook-json-check](./cmd/hook-json-check) ensures that all checked-in `*.json`
files parse as valid JSON.
* [hook-dead-symlinks](./cmd/hook-dead-symlinks) prevents dead symlinks from being
added to a repository.
* TODO: [hook-binary-files](./cmd/hook-binary-files) prevents binary files from being
added to a repository.
* TODO: [hook-big-files](./cmd/hook-big-files) prevents big files from being added
to a repository.
* TODO: [hook-dockerfile](./cmd/hook-dockerfile) lint Dockerfiles.
* TODO: [hook-k8s-manifests](./cmd/hook-k8s-manifests) lint k8s manifests.
* TODO: [hook-terraform](./cmd/hook-terraform) lint terraform plans.
* TODO: [hook-python](./cmd/hook-python) lint Python.
* TODO: [hook-img-optimize](./cmd/hook-img-optimize) losslessly-compress images.
* TODO: [hook-opa](./cmd/hook-opa) scan repositories with OPA.
* TODO: [hook-markdown](./cmd/hook-markdown) validate Markdown syntax.


### commit-msg

* [hook-emojify-commit-msg](./cmd/hook-emojify-commit-msg) adds emoji prefixes to
commit summaries.

### post-commit

### pre-push

* TODO: [hook-signoff-checker](./cmd/hook-signoff-checker) checks that commits are signed-off.
* TODO: [hook-gpg-sign-checker](./cmd/hook-gpg-sign-checker) checks that commits are
GPG signed.
