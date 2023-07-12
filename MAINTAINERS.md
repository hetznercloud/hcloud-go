# Maintainers Guide

This guide is intended for _maintainers_.

## Release Branches

All development targets the `main` branch. New releases for the 2.x series are cut from this.

For the 1.x series which will be supported until at least September 1 2023 we also have a `release-1.x` branch where we try to backport all bug fixes.

Backports are done by [tibdex/backport](https://github.com/tibdex/backport). Apply the label `backport release-1.x` to the PRs and once they are merged a new PR is opened.
