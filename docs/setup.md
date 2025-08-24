# Setup & Environment (Windows, GoMT4 Project)

Audience: **beginner-friendly**. Comments in code are **English only**. This setup is tailored to your repo layout (entrypoint: `examples/main.go`) and to the fact that **pb files are fetched as a Go module from the network** — no extra repo cloning.

---

## 1) Prerequisites

* **Windows 10/11 recommended** (Windows 7/8 may work but not tested)
* **Git** ([https://git-scm.com/download/win](https://git-scm.com/download/win))
* **Go ≥ 1.21** ([https://go.dev/dl/](https://go.dev/dl/))
* **VS Code** + Go extension (by Google)

> Why: Go builds and runs the project; VS Code is the recommended IDE.

---

## 2) Where the pb files live (and how they are connected)

* The generated protobuf stubs are published as a **Go module** at:

  ```
  git.mtapi.io/root/mrpc-proto/mt4/libraries/go
  ```
* Your code imports it like this (no `.git` suffix):

  ```go
  import pb "git.mtapi.io/root/mrpc-proto/mt4/libraries/go"
  ```
* Because this is a Go module, **you do not need to clone another repository**. Go will download the module automatically on the first build.

---

## 3) Authentication for private modules (one‑time setup)

If the module is private, tell Go where it may fetch private code and make sure Git credentials are configured.

```powershell
# Allow private modules under this host
go env -w GOPRIVATE=git.mtapi.io

# Make sure your Git credentials (SSH or HTTPS token) work for git.mtapi.io
# For HTTPS, run once and cache credentials:
# git config --global credential.helper manager-core
```

> After this, `go build` / `go run` will resolve `pb` from the network automatically.

---

## 4) Project checkout & first run

```powershell
# Clone ONLY your GoMT4 repo
cd C:\Users\malin
git clone <YOUR-GoMT4-REPO-URL> GoMT4
cd GoMT4

# Install deps (this will fetch the pb module too)
go mod tidy

# Run the example (entrypoint lives in examples/main.go)
go run ./examples/main.go
```

You should see logs and the first RPC interactions.

---

## 5) How to update the pb module (when API changes)

You control updates via standard Go tooling. Pick one of the following:

* **Move to the latest compatible version**

  ```powershell
  go get -u git.mtapi.io/root/mrpc-proto/mt4/libraries/go@latest
  go mod tidy
  ```

* **Pin to an explicit version** (recommended for production)

  ```powershell
  go get git.mtapi.io/root/mrpc-proto/mt4/libraries/go@vX.Y.Z
  go mod tidy
  ```

  In `go.mod` you will see a line like:

  ```
  require git.mtapi.io/root/mrpc-proto/mt4/libraries/go vX.Y.Z
  ```

* **Rollback** to a known good version

  ```powershell
  go get git.mtapi.io/root/mrpc-proto/mt4/libraries/go@vPrevious
  go mod tidy
  ```

> No manual copying of `.proto` or `.pb.go` is needed. The module is the single source of truth.

---

## 6) Optional: work offline (vendoring)

If you need fully offline builds after the first fetch:

```powershell
go mod vendor
```

This will copy dependencies (including `pb`) into `./vendor/`. CI can then build without external network access.

---

## 7) VS Code configuration (repo‑specific)

Create `.vscode/` in the repo root.

**.vscode/launch.json**

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Go: Run main example",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/examples/main.go",
      "envFile": "${workspaceFolder}/.env",
      "cwd": "${workspaceFolder}"
    }
  ]
}
```

**.vscode/settings.json**

```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "gopls": { "ui.semanticTokens": true },
  "editor.formatOnSave": true
}
```

---

## 8) Environment file (example)

Create `.env` in the repo root. Commit only `.env.example` and **ignore** `.env`.

**.env.example**

```dotenv
MT4_TERMINAL_PATH=C:\\Program Files (x86)\\MetaTrader 4\\terminal64.exe
RPC_LISTEN_ADDR=127.0.0.1:50051
MT4_LOGIN=12345678
MT4_PASSWORD=your-investor-or-trade-password
MT4_SERVER=YourBroker-Server
LOG_LEVEL=info
SYMBOLS=EURUSD,GBPUSD,USDJPY
```

**.gitignore** addition

```
.env
```

---

## 9) Troubleshooting (pb module)

* **`protoc-gen-go: not found`** — not required for consumers; you are not generating code locally.
* **`module git.mtapi.io/... not found`** — check `GOPRIVATE` and your Git credentials for `git.mtapi.io`.
* **`no matching versions`** — specify a valid tag (e.g., `@vX.Y.Z`) or use `@latest`.
* **build breaks after update** — pin back to the last known good version in `go.mod`.

---

## 10) Next steps

* Continue to [Architecture & Data Flow](architecture.md).
* Add [Troubleshooting & FAQ](troubleshooting.md) with common MT4 specifics (symbol suffixes, Digits/LotStep rounding).
* Use the Cookbook for practical code recipes.
