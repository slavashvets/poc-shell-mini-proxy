# shell-proxy

A **minimal HTTP ↔ interactive shell** gateway written in idiomatic Go.  
Spawns one long-lived `/bin/sh` per **UUID** and streams stdout/stderr to the client via **Server-Sent Events (SSE)**.

## Features

- **Interactive shells** scoped by UUID
- **SSE**-based real-time streaming of command output
- **Simple HTTP API**:

  | Method   | Path      | Description                                 |
  | -------- | --------- | ------------------------------------------- |
  | `PUT`    | `/{uuid}` | Create a new shell session                  |
  | `POST`   | `/{uuid}` | Execute a command inside that shell         |
  | `GET`    | `/{uuid}` | Stream stdout/stderr via SSE                |
  | `DELETE` | `/{uuid}` | Terminate the shell and destroy the session |

## Prerequisites

- **Go 1.24+** (darwin/arm64 & x86_64)
- **Homebrew** (macOS) – for easy installs

```bash
brew install go            # Go toolchain
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
brew install hurl          # API integration tests
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `$PATH`.

## Setup & Development

```bash
git clone <repo-url> shell-proxy
cd shell-proxy

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

---

## Quick Manual Test (cURL)

```bash
uuid=$(uuidgen)

# 1. Create session
curl -X PUT http://localhost:8080/$uuid

# 2. Execute command
curl -X POST http://localhost:8080/$uuid \
     -d 'echo hello && echo done && exit'

# 3. Stream output (SSE)
curl -N http://localhost:8080/$uuid
# → data:hello
# → data:done

# 4. Delete session
curl -X DELETE http://localhost:8080/$uuid
```

## Automated Tests (Hurl)

```bash
uuid=$(uuidgen)
unkown=00000000-0000-0000-0000-deadbeefdead 

hurl --variable uuid=$uuid tests/happy_flow.hurl
hurl --variable uuid=$uuid tests/duplicate_session.hurl
hurl --variable uuid=$uuid tests/deleted_session.hurl
hurl --variable unknown=$unkown tests/unknown_session.hurl
```

## Security Notice

This gateway executes **arbitrary shell code**.
Use only on trusted hosts or protect with appropriate sandboxing/ACLs.
