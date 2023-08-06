[![CodeQL](https://github.com/blck-snwmn/gh-prls/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/blck-snwmn/gh-prls/actions/workflows/github-code-scanning/codeql)
[![release](https://github.com/blck-snwmn/gh-prls/actions/workflows/release.yml/badge.svg)](https://github.com/blck-snwmn/gh-prls/actions/workflows/release.yml)

A gh extension to list PRs

## Create
```bash
$ gh extension create --precompiled=go gh-prls
```

## Build
```bash
go build
```

## Install
local
```bash
$ gh extension install .
```

from repository
```bash
$ gh extension install https://github.com/blck-snwmn/gh-prls
```

## Run
```bash
$ gh prls
```

## Update
```bash
$ git tag -a v1.1.0 -m 'update'
```
