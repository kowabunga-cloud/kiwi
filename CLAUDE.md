# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About

Kiwi is the Kowabunga SD-WAN agent. It runs as a daemon on SD-WAN nodes and provides network services via RPC — primarily DNS management (with forwarding). It connects to a central Kowabunga platform over WebSocket and registers itself as an RPC service provider.

## Commands

```sh
# Build
go build ./...

# Run all tests
go test ./...

# Run a single test
go test ./internal/kiwi/ -run TestKiwiConfigParser

# Format
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...

# Tidy dependencies
go mod tidy

# Install pre-commit hooks (required before committing)
pre-commit install --install-hooks
```

## Architecture

The agent is structured around `github.com/kowabunga-cloud/common`, a shared library that provides the base agent framework, RPC server, WebSocket connectivity, and CLI parsing.

**Entry point**: `cmd/kiwi/main.go` → `kiwi.Daemonize()` (parses config, initializes logger, creates agent, runs event loop).

**Core types** (`internal/kiwi/`):
- `KiwiAgent` (`kiwi.go`) — embeds `*agents.KowabungaAgent` plus a `*DnsServer`. Sets `PostFlight` callback to shut down DNS on exit. Created via `NewKiwiAgent()`.
- `DnsServer` (`dns.go`) — in-memory UDP DNS server using `github.com/miekg/dns`. Stores A records as `"hostname.domain.": "ip1,ip2"`. Serves local records first; unknown queries are forwarded to recursors (default: 9.9.9.9, 149.112.112.112). Thread-safe via `sync.Mutex`.
- `Kiwi` (`services.go`) — RPC handler struct registered with the agent. Implements: `Capabilities`, `Reload`, `CreateDnsZone`, `DeleteDnsZone`, `CreateDnsRecord`, `UpdateDnsRecord`, `DeleteDnsRecord`. DNS zone methods are no-ops; `Reload` bulk-replaces all records.
- `KiwiAgentConfig` (`config.go`) — YAML config with `global` (from common library: ID, endpoint, apiKey, logLevel) and `dns` (port, recursors) sections. See `config/kiwi.yml` for a reference config.

**Data flow**: The platform sends RPC calls → `Kiwi` methods handle them → mutate `DnsServer.records` in-memory → DNS queries resolved locally or forwarded.

## Commit Convention

Unscoped [Conventional Commits](https://www.conventionalcommits.org/) are required (enforced by pre-commit). Types that trigger releases:
- `feat:` → minor version bump
- `feat!:` → major version bump
- `fix:`, `perf:`, `chore:` → patch bump
- `docs:` → no release

Releases are automated via semantic-release on the `master` branch.
