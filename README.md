# githooks

Just some git hooks I use.

## Use

To build all tools, `$ ./build.sh`. The build script will print instructions for
setting up your git templates dir. From that point forward, new git repositories
will automatically use these hooks.

To add to an existing project, run `$ usegithooks`, which is installed by `build.sh`.

All githooks are intended to be "smart" enough to know when they are and are not
relevant.

## Tools

### pre-commit

* [hook-interactive-secrets](./cmd/hook-interactive-secrets) checks for secrets,
and provides interactive override.
* [hook-go](./cmd/hook-go) lint, vet, and test Go code.
* [hook-json-check](./cmd/hook-json-check) ensures that all checked-in `*.json`
files parse as valid JSON.
* [hook-dead-symlinks](./cmd/hook-dead-symlinks) prevents dead symlinks from being
added to a repository.

### commit-msg

* [hook-emojify-commit-msg](./cmd/hook-emojify-commit-msg) adds emoji prefixes to
commit summaries.

### post-commit

### pre-push

* [hook-signoff-checker](./cmd/hook-signoff-checker) checks that commits are signed-off.
* [hook-gpg-sign-checker](./cmd/hook-gpg-sign-checker) checks that commits are
GPG signed.
