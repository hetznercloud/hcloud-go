# Maintainers Guide

This guide is intended for _maintainers_.

## Release Branches

All development targets the `main` branch. New releases for the latest major version series are cut from this branch.

For the older major versions, we also have `release-<major>.x` branches where we try to backport all bug fixes.

Backports are done by [tibdex/backport](https://github.com/tibdex/backport). Apply the label `backport release-<major>.x` to the PRs and once they are merged a new PR is opened.
