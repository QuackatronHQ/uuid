# How to contribute

We definitely welcome patches and contribution to this project!

### Tips

Commits must be formatted according to the [Conventional Commits Specification](https://www.conventionalcommits.org).

Always try to include a test case! If it is not possible or not necessary,
please explain why in the pull request description.

Before sending a pull request, run:

```sh
gofmt -l .
go vet ./...
go test -race ./...
```

The module supports Go 1.19 and later. Do not use language features or
standard library APIs added after Go 1.19.

### Releasing

Commits that would precipitate a SemVer change, as described in the Conventional
Commits Specification, will trigger [`release-please`](https://github.com/google-github-actions/release-please-action)
to create a release candidate pull request. Once submitted, `release-please`
will create a release.

For tips on how to work with `release-please`, see its documentation.

### Legal requirements

In order to protect both you and ourselves, you will need to sign the
[Contributor License Agreement](https://cla.developers.google.com/clas).

You may have already signed it for other Google projects.
