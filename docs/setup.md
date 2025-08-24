# Setup & Environment (Windows, GoMT4 Project)

Audience: **beginner‑friendly**. Comments in code are **English only**. This guide matches your repo layout (entrypoint: `examples/main.go`) and the fact that **pb files are fetched as a public Go module from the network** — no extra repo cloning.

---

## 1) Prerequisites

* **Windows 10/11 recommended** (Windows 7/8 may work but not tested)
* **Git** — [https://git-scm.com/download/win](https://git-scm.com/download/win)
* **Go ≥ 1.21** — [https://go.dev/dl/](https://go.dev/dl/)
* **VS Code** + Go extension (by Google)

> Why: Go builds/runs the project; VS Code gives you debug + ergonomics.

---

## 2) Where the pb files come from (and how they are connected)

* Generated protobuf stubs live in a **public Go module**:

  ```
  git.mtapi.io/root/mrpc-proto/mt4/libraries/go
  ```
* Your code imports it like this (note: **no `.git` suffix**):

  ```go
  import pb "git.mtapi.io/root/mrpc-proto/mt4/libraries/go"
  ```
* Because it is a Go module, **you do not clone another repository**. Go will download the module automatically during the first build.

**Update policy:** You manage the pb version via normal Go commands (see §5). No manual copying of `.proto` or `.pb.go` is needed.

---

## 3) Project checkout & first run

```powershell
# Clone ONLY your GoMT4 repo
cd C:\Users\malin
git clone <YOUR-GoMT4-REPO-URL> GoMT4
cd GoMT4

# Resolve & download dependencies (this will fetch the pb module too)
go mod tidy

# Run the example (entrypoint lives in examples/main.go)
go run ./examples/main.go
```

You should see logs and the first RPC interactions.

---

## 4) What exactly does `go mod tidy` do?

* Reads your `go.mod` and computes the **minimal set of required modules** for your imports.
* **Downloads** any missing modules to your local module cache (not into your repo).
* **Adds/removes lines** in `go.mod` as needed and updates `go.sum` (checksums).
* It **does not overwrite your source files** and **does not create duplicates** of your code. It only manages `go.mod`/`go.sum` and your local cache.

> If the pb dependency is already in cache and the version hasn’t changed, `go mod tidy` will not re‑download or overwrite anything.

---

## 5) How to update (or pin) the pb module

Choose one path:

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

  Your `go.mod` will contain:

  ```
  require git.mtapi.io/root/mrpc-proto/mt4/libraries/go vX.Y.Z
  ```

* **Rollback** if something broke after an update

  ```powershell
  go get git.mtapi.io/root/mrpc-proto/mt4/libraries/go@vPrevious
  go mod tidy
  ```

---

## 6) What is `go mod vendor` and how to work offline

* `go mod vendor` **copies your resolved dependencies** (including pb) into a local `./vendor/` folder inside the repo.
* With Go **1.14+**, if `vendor/` exists *and* the main module’s `go` directive is ≥ 1.14, the `go` tool **prefers the vendor folder by default**. Otherwise run with `-mod=vendor`.
* After vendoring (and a first successful fetch), you can **build offline** because the compiler reads packages from `./vendor/` instead of the internet.

**Use cases:** reproducible builds in CI without external network; air‑gapped environments.

**Command:**

```powershell
go mod vendor
# (optional) enforce vendor mode explicitly
go build -mod=vendor ./...
```

---

## 7) VS Code debug configuration — what is this block?

This JSON is a VS Code file named **`.vscode/launch.json`**. It tells VS Code *how to run your app in debug mode*.

* `program`: which Go entrypoint to run (here: `examples/main.go`).
* `envFile`: path to a file whose variables will be injected into the process environment at launch (see §8).
* `cwd`: working directory for the process.

Create `.vscode/launch.json` with:

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

> After this, pressing **F5** in VS Code will run your app with a debugger attached.

---

## 8) What is `.env` and why do we need it?

* `.env` is a simple **key=value** text file with configuration values. Example:

  ```dotenv
  MT4_TERMINAL_PATH=C:\\Program Files (x86)\\MetaTrader 4\\terminal64.exe
  RPC_LISTEN_ADDR=127.0.0.1:50051
  MT4_LOGIN=12345678
  MT4_PASSWORD=your-investor-or-trade-password
  MT4_SERVER=YourBroker-Server
  LOG_LEVEL=info
  SYMBOLS=EURUSD,GBPUSD,USDJPY
  ```
* In our setup, **VS Code reads this file** (because we set `envFile` in `launch.json`) and **injects these variables** into the process environment when starting your app.
* That means: you don’t need to write a `.env` parser in code — VS Code already passes the variables to your app under the debugger.
* Commit only the **example** (`.env.example`) without secrets and add the real `.env` to `.gitignore`.

**.gitignore addition**

```
.env
```

---

## 9) First run checklist (quick)

1. Start MT4 terminal once manually.
2. Make sure firewall allows your port (e.g. 50051):

   ```powershell
   New-NetFirewallRule -DisplayName "GoMT4 gRPC" -Direction Inbound -Protocol TCP -LocalPort 50051 -Action Allow
   ```
3. Run the example:

   ```powershell
   cd C:\Users\malin\GoMT4
   go run ./examples/main.go
   ```

---

## 10) Common pitfalls

* **`module git.mtapi.io/... not found`** → temporary network issue; try again or pin a tag `@vX.Y.Z`.
* **`no matching versions`** → specify a valid tag or use `@latest`.
* **Timeouts / no connection** → check `.env` values, firewall, MT4 terminal connectivity.
* **Symbol not found (`EURUSD`)** → broker may use suffix (e.g. `EURUSD.m`); ensure the symbol is visible in MT4.
* **Volume/price rejected** → always round using symbol `Digits` and `LotStep`.

