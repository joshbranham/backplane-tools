# CLAUDE.md - backplane-tools

## Overview

OpenShift backplane-tools: a CLI tool manager that installs, removes, and upgrades tools used to interact with OpenShift clusters. Built in Go using the Cobra CLI framework.

## Build Commands

```bash
make all        # Default: vet, fmt, mod, lint, test, build
make build      # Build binaries via goreleaser (linux/darwin, 386/amd64/arm64)
make test       # Run all tests
make lint       # Run golangci-lint (v1.55.0)
make fmt        # Format code with gofmt
make vet        # Run go vet
make mod        # Tidy go.mod
make coverage   # Generate code coverage via hack/codecov.sh
```

## Project Structure

```
cmd/                  # CLI subcommands (install, list, remove, upgrade)
pkg/tools/            # Tool implementations and registry
pkg/tools/base/       # Base types: Default, Github, Mirror
pkg/sources/          # Download sources (GitHub API, Cloud Storage, HTTP, OpenShift mirror)
pkg/utils/            # Utilities (checksum, unarchive, gpg, file ops)
hack/                 # CI scripts (codecov)
main.go               # Entry point (Cobra root command)
```

## Architecture

### Tool Interface (`pkg/tools/tools.go`)

All managed tools implement the `Tool` interface: `Name()`, `ExecutableName()`, `Install()`, `Configure()`, `Remove()`, `Installed()`, `InstalledVersion()`, `LatestVersion()`.

Tools are registered in the `AllTools()` function in `pkg/tools/tools.go`.

### Base Types (`pkg/tools/base/`)

- **Default** - Core functionality: directory management, symlinks, version tracking. Install path: `~/.local/bin/backplane/<tool>/<version>/`
- **Github** - Extends Default. Fetches releases from GitHub API with token auth support.
- **Mirror** - Extends Default. Downloads from mirror.openshift.com.

### Adding a New Tool

1. Create `pkg/tools/<toolname>/<toolname>.go`
2. Define a struct embedding a base type (e.g., `base.Github`)
3. Implement `New()` constructor and `Install()` method (download, verify checksum, extract, symlink)
4. Register in `pkg/tools/tools.go` `AllTools()` function

### Key Design Decisions

- Tools downloaded directly from original sources (no central package server)
- All versions retained locally; latest symlinked to `~/.local/bin/backplane/latest/`
- Downloads verified against checksums
- Parallel installation using WaitGroups
- OS/arch aliasing handles naming variations (amd64/x86_64, darwin/mac, etc.)

## Code Quality

- Linter config: `.golangci.yaml` (disable-all with explicit allowlist of ~60 linters)
- CI: OpenShift Prow (`.ci-operator.yaml`), not GitHub Actions
- Go 1.21+ required (toolchain 1.22.1)
- CGO_ENABLED=0, trimpath for reproducible builds
