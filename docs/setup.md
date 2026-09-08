# 💻 Setup & Environment (Windows, GoMT4 Project)

Audience: **beginner-friendly**. Comments in code are **English only**.

---

## 📦 1) Prerequisites

* **Windows 10/11 recommended** (Windows 7/8 may work but not tested)
* **Git** ([https://git-scm.com/download/win](https://git-scm.com/download/win))
* **Go ≥ 1.21** ([https://go.dev/dl/](https://go.dev/dl/))
* **VS Code** + Go extension (by Google)

👉 Why: Go builds and runs the project; VS Code is the recommended IDE.

---

## 📂 2) Where the pb files live

* The generated protobuf stubs are published as a **Go module**:

  ```
  github.com/MetaRPC/GoMT4
  ```
* Import without `.git` suffix:

  ```go
  import pb "github.com/MetaRPC/GoMT4"
  ```
* You do **not need to clone another repo**. Go fetches it automatically.

---

## ⚙️ 3) How Go manages pb files

* `go mod tidy` will:

  * Download missing modules (including pb).
  * Update `go.sum` with checksums.
  * It does **not overwrite your code**, just ensures deps exist.
* `go get -u` → updates to newer pb version.

---

## 📦 4) Working offline: `go mod vendor`

* `go mod vendor` → copies all deps into `vendor/`.
* Project can then build offline.
* Useful for CI/CD builds without internet.

---

## ⚙️ 5) Configuration with `config.json`

Config file path:

```
examples/config/config.json
```

Example:

```json
{
  "Login": 501401178,
  "Password": "v8gctta",
  "Server": "RoboForex-Demo",
  "DefaultSymbol": "EURUSD"
}
```

🔑 Use **investor password** for read-only unless you need trading.

---

## ▶️ 6) Project checkout & first run

```powershell
# Clone ONLY your GoMT4 repo
cd C:\Users\malin
git clone <YOUR-GoMT4-REPO-URL> GoMT4
cd GoMT4

# Install deps
go mod tidy

# Run example (entrypoint: examples/main.go)
go run ./examples/main.go
```

Expected: logs + first RPC interactions.

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
      "program": "${workspaceFolder}/examples/main.go",
      "cwd": "${workspaceFolder}"
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

1. Start MT4 terminal manually once (initializes data).

2. Firewall: allow port if connecting externally. For local `127.0.0.1`, not needed.

3. Run:

   ```powershell
   cd C:\Users\malin\GoMT4
   go run ./examples/main.go
   ```

4. You should see logs like `listening on mt4.mrpc.pro:443`.

---

## ⚠️ 9) Common pitfalls

* **Timeouts / no connection** → Check `config.json`, firewall, MT4 connectivity.
* **Symbol not found (`EURUSD`)** → Broker may add suffix (e.g., `EURUSD.m`). Ensure symbol is visible.
* **Volume or price rejected** → Always round with `Digits` and `LotStep` values.
