# poc-shell-mini-proxy

> **Proof‑of‑concept (PoC)** HTTP ↔ interactive shell gateway. 100 % experimental, **not intended for production use**.

I spun this up to learn some Go basics and to explore a tooling pattern I might eventually proxy into containers that have no direct shell access via API (for example, plain Kubernetes pods or locked‑down cloud runtimes where an exec endpoint is disabled). I am **not a Go developer** – so expect rough edges, corner‑cutting and breaking changes while I experiment.

## What it does

- Spawns one long‑lived `/bin/sh` per **UUID**
- Maps a very small **HTTP API** onto that shell
- Streams combined `stdout/stderr` back to the caller using **Server‑Sent Events (SSE)**

| Method   | Path      | Purpose                           |
| -------- | --------- | --------------------------------- |
| `PUT`    | `/{uuid}` | Create a shell session            |
| `POST`   | `/{uuid}` | Run a command inside that session |
| `GET`    | `/{uuid}` | Follow live output via SSE        |
| `DELETE` | `/{uuid}` | Terminate the session & clean up  |

⚠️ **Security notice:** this literally executes arbitrary shell code. Run only on trusted hosts or inside a very tight sandbox.

## Prerequisites

- **Go 1.24+**

```bash
brew install go            # Go toolchain
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
brew install hurl          # API integration tests
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `$PATH`.

## Setup & Development

```bash
# Format + imports
goimports -w .
go fmt ./...

# Static analysis
go vet ./...
staticcheck ./...

# Run with hot-reload
go run .

# Build release binary
go build -o shell-proxy .
```

Server listens on **`localhost:8080`** by default.

## Automated Tests (Hurl)

```bash
uuid=$(uuidgen)
unknown=00000000-0000-0000-0000-deadbeefdead

hurl --variable uuid=$uuid tests/happy_flow.hurl
hurl --variable uuid=$uuid tests/duplicate_session.hurl
hurl --variable uuid=$uuid tests/deleted_session.hurl
hurl --variable unknown=$unknown tests/unknown_session.hurl
```
