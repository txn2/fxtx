# fxtx
[![Go Report Card](https://goreportcard.com/badge/github.com/txn2/fxtx)](https://goreportcard.com/report/github.com/txn2/fxtx)
[![GoDoc](https://godoc.org/github.com/txn2/fxtx?status.svg)](https://godoc.org/github.com/txn2/fxtx)
[![Docker Container Image Size](https://shields.beevelop.com/docker/image/image-size/txn2/fxtx/latest.svg)](https://hub.docker.com/r/txn2/fxtx/)
[![Docker Container Layers](https://shields.beevelop.com/docker/image/layers/txn2/fxtx/latest.svg)](https://hub.docker.com/r/txn2/fxtx/)

WIP: TXN2 Fake transmission

## Install
### MacOS
`brew install txn2/tap/fxtx`

## Development

### Test Release

```bash
goreleaser --skip-publish --rm-dist --skip-validate
```

### Release

```bash
GITHUB_TOKEN=$GITHUB_TOKEN goreleaser --rm-dist
```