# fxtx
WIP: TXN2 Fake transmission

## Development

### Test Release

```bash
goreleaser --skip-publish --rm-dist --skip-validate
```

### Release

```bash
GITHUB_TOKEN=$GITHUB_TOKEN goreleaser --rm-dist
```