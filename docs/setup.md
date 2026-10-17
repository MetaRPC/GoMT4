# 💻 Setup & Environment (GoMT4 Project)

Audience: **beginner-friendly**. Comments in code are **English only**.

---

## 📦 1) Prerequisites

* **Windows, Linux or macOS** — the client is plain Go; nothing MetaTrader-specific runs on your machine
* **Git** ([https://git-scm.com/downloads](https://git-scm.com/downloads))
* **Go ≥ 1.23** ([https://go.dev/dl/](https://go.dev/dl/)) — the modules declare `go 1.23.6`
* **VS Code** + Go extension (by Google), or any Go IDE
* An **MT4 account** (login, password, broker server name) and a **MetaRPC API key** from [https://mrpc.pro/signup](https://mrpc.pro/signup) (without a key the calls run as a trial)

👉 Why: Go builds and runs the project. The MT4 terminal itself is run by MetaRPC in the cloud: your program connects to `mt4.mrpc.pro:443` (gRPC over TLS), and MetaRPC starts and manages an MT4 terminal logged in to your account. You do not install MetaTrader, open ports, or run any server.

---

## 📂 2) Where the code lives

The repository has two Go modules (tied together by `go.work` at the root):

| Folder | Module | What it is |
| --- | --- | --- |
| `package/` | `github.com/MetaRPC/GoMT4/package` | Generated protobuf messages, enums and gRPC clients |
| `examples/` | `github.com/MetaRPC/GoMT4/examples` | The client (`examples/mt4`: `MT4Account`, `MT4Service`), the runnable `main.go` and its `config/` |

* The generated stubs can be fetched by Go on their own:

  ```bash
  go get github.com/MetaRPC/GoMT4/package
  ```

  ```go
  import pb "github.com/MetaRPC/GoMT4/package"
  ```
* The client (`MT4Account` / `MT4Service`) is part of the examples module and is used **from a clone** of the repository:

  ```go
  import "github.com/MetaRPC/GoMT4/examples/mt4"
  ```

  To use it from your own module, clone the repo and point a `replace` directive at it (adjust the paths):

  ```
  require (
      github.com/MetaRPC/GoMT4/examples v0.0.0
      github.com/MetaRPC/GoMT4/package v0.0.0
  )

  replace (
      github.com/MetaRPC/GoMT4/examples => ../GoMT4/examples
      github.com/MetaRPC/GoMT4/package => ../GoMT4/package
  )
  ```

---

## ⚙️ 3) How Go manages dependencies

* `go mod tidy` will:

  * Download missing modules (gRPC, protobuf, uuid).
  * Update `go.sum` with checksums.
  * It does **not overwrite your code**, just ensures deps exist.
* `go get -u github.com/MetaRPC/GoMT4/package@latest` → updates the generated stubs (in your own module).

---

## 📦 4) Working offline: `go mod vendor`

* `go mod vendor` (in your own module) → copies all deps into `vendor/`; inside this repository, which is a Go workspace (`go.work`), the equivalent is `go work vendor` at the root.
* The project can then build offline.
* Useful for CI/CD builds without internet. Calls still need network access to `mt4.mrpc.pro:443` at run time.

---

## ⚙️ 5) Configuration with `config.json`

Config file path:

```
examples/config/config.json
```

Example:

```json
{
  "Login": 12345678,
  "Password": "<your-mt4-password>",
  "Server": "MetaQuotes-Demo",
  "DefaultSymbol": "EURUSD",
  "ApiKey": "<your-mrpc-api-key>"
}
```

* `Server` is the broker server name as shown in MetaTrader (the example connects with `ConnectByServerName`).
* `"Login": 0` makes the example open a fresh demo account on `Server` and use it.
* `"ApiKey": "TRIAL"` (or empty) uses the `MRPC_API_KEY` environment variable if it is set, otherwise the trial.

🔑 An **investor password** logs in read-only: account data, quotes and history work, trading calls are refused.

---

## ▶️ 6) Project checkout & first run

```powershell
git clone https://github.com/MetaRPC/GoMT4.git
cd GoMT4/examples

# Install deps
go mod download

# Run the example (entrypoint: examples/main.go, reads config/config.json from the current folder)
go run .
```

Expected: `✅ Connected to MT4 server`, the symbol parameters, then the output of each example call.

---

## 🐞 7) VS Code debug configuration

File: `.vscode/launch.json`

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Go: Run main example",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/examples",
      "cwd": "${workspaceFolder}/examples"
    }
  ]
}
```

File: `.vscode/settings.json`

```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "gopls": {
    "ui.semanticTokens": true
  },
  "editor.formatOnSave": true
}
```

---

## ✅ 8) First run checklist

1. `examples/config/config.json` has your login, password, server name and API key (or `"Login": 0` for a demo account).

2. Outbound HTTPS (port 443) to `mt4.mrpc.pro` is allowed by your network.

3. Run:

   ```powershell
   cd GoMT4/examples
   go run .
   ```

4. You should see `✅ Connected to MT4 server`. The first connect can take a while: MetaRPC starts a terminal for the account (the example waits up to 30 seconds).

---

## ⚠️ 9) Common pitfalls

* **Timeouts / no connection** → Check `config.json` (login, password, server name), the API key, and that `mt4.mrpc.pro:443` is reachable.
* **Symbol not found (`EURUSD`)** → Broker may add suffix (e.g., `EURUSD.m`). Check the names with `ShowAllSymbols`.
* **Volume or price rejected** → Round prices with `Digits` and volumes with `VolumeMin` / `VolumeStep` from `SymbolParams`.
