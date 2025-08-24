# Setup & Environment (Windows, GoMT4 Project)

Audience: **beginner-friendly**. Comments in code are **English only**.

---

## 1) Prerequisites

* **Windows 10/11 recommended** (Windows 7/8 may work but not tested)
* **Git** ([https://git-scm.com/download/win](https://git-scm.com/download/win))
* **Go ≥ 1.21** ([https://go.dev/dl/](https://go.dev/dl/))
* **VS Code** + Go extension (by Google)

> Why: Go builds and runs the project; VS Code is the recommended IDE.

---

## 2) Where the pb files live

* The generated protobuf stubs are published as a **Go module**:

  ```
  git.mtapi.io/root/mrpc-proto/mt4/libraries/go
  ```
* Your code imports it like this (no `.git` suffix):

  ```go
  import pb "git.mtapi.io/root/mrpc-proto/mt4/libraries/go"
  ```
* Because this is a Go module and it is public, **you do not need to clone another repository**. Go will automatically fetch it on the first build.

---

## 3) How Go manages pb files

* Running `go mod tidy` will:

  * Download any missing modules (including the pb module).
  * Update `go.sum` with checksums.
  * It does **not overwrite your local code**; it only ensures dependencies in `go.mod` are present.
* If the pb module is already downloaded, `go mod tidy` just verifies it. If a new version is requested (`go get -u`), then Go will update it.

---

## 4) Working offline: `go mod vendor`

* `go mod vendor` copies all dependencies (including pb stubs) into a local `vendor/` folder.
* After running it, the project can be built offline because Go will use the `vendor/` folder instead of the internet.
* This is optional; useful for CI builds without internet access.

---

## 5) Configuration with `config.json`

The project reads account and server settings from a JSON config file:

```
examples/config/config.json
```

Example config:

```json
{
  "Login": 501401178,
  "Password": "v8gctta",
  "Server": "RoboForex-Demo",
  "DefaultSymbol": "EURUSD"
}
```

Adjust these values for your broker before running the project. Prefer using **investor (read‑only) password** for safety unless you need trading operations.

---

## 6) Project checkout & first run

```powershell
# Clone ONLY your GoMT4 repo
cd C:\Users\malin
git clone <YOUR-GoMT4-REPO-URL> GoMT4
cd GoMT4

# Install deps (fetches pb module too)
go mod tidy

# Run the example (entrypoint lives in examples/main.go)
go run ./examples/main.go
```

You should see logs and the first RPC interactions.

---

## 7) VS Code debug configuration

The file `.vscode/launch.json` tells VS Code **how to run and debug your program**.

* `program`: which Go file to launch (`examples/main.go`).
* `cwd`: current working directory when running.

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
  "gopls": {
    "ui.semanticTokens": true
  },
  "editor.formatOnSave": true
}
```

---

## 8) First run checklist

1. Start MT4 terminal once manually (to initialize its data folders).
2. Ensure firewall allows `RPC_LISTEN_ADDR` port (e.g. 50051) if you plan to connect from outside your machine. For local use on `127.0.0.1`, this is usually not required.
3. Run the example:

   ```powershell
   cd C:\Users\malin\GoMT4
   go run ./examples/main.go
   ```
4. You should see logs like `listening on 127.0.0.1:50051` and first RPC interactions.

---

## 9) Common pitfalls

* **Timeouts / no connection** → Check `config.json` values, firewall, MT4 terminal connectivity.
* **Symbol not found (`EURUSD`)** → Broker may add suffix (e.g. `EURUSD.m`). Ensure symbol is visible in MT4.
* **Volume or price rejected** → Always round with `Digits` and `LotStep` values.

