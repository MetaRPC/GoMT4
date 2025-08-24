# 🧾 ConfigExample (GoMT4)

**Goal:** show the minimal **config.json** used by this repo and how it’s consumed by the examples.

> Real files:
>
> * Config file: `examples/config/config.json`
> * Config loader: `examples/config/config.go` (reads JSON into a struct)
> * Used by: `examples/main.go`

---

## 📍 Location

```
examples/config/config.json
```

This file is read at startup by the example app.

---

## 🧩 Schema (what fields mean)

```json
{
  "Login": 501401178,
  "Password": "v8gctta",
  "Server": "RoboForex-Demo",
  "DefaultSymbol": "EURUSD"
}
```

* **Login** (number) — MT4 account login.
* **Password** (string) — investor or trade password. For demos, prefer **investor** (read‑only).
* **Server** (string) — exact broker server name (e.g., `RoboForex-Demo`).
* **DefaultSymbol** (string) — instrument to use by default (must match broker’s name; suffixes like `EURUSD.m` are different symbols).

---

## 🛠️ Edit & run

1. Open `examples/config/config.json` and fill your values.
2. Launch the example:

```powershell
cd GoMT4
go mod tidy
go run ./examples/main.go
```

If config is valid and MT4 is connected, you’ll see logs from the example service.

---

## 🔍 Validation (quick checks)

* Login is numeric, password non‑empty.
* Server name exactly matches MT4 (check in terminal login dialog).
* `DefaultSymbol` exists and is visible in Market Watch (`Show All`).

Minimal runtime probe (uses real calls):

```go
sum, err := account.AccountSummary(ctx)
if err != nil { return err }
_, err = account.Quote(ctx, cfg.DefaultSymbol)
if err != nil { return err }
```

---

## 🔐 Security notes

* `config.json` in this repo is meant for **local development**.
* Do **not** commit real trading credentials to a public repo. Use demo creds or keep the repo private.
* If you need secrets isolation, migrate to environment variables or an external secret store (optional; not required by this project).

---

## ⚠️ Common errors

* **`symbol not found`** → wrong/hidden symbol; check suffix (`EURUSD.m`) and Show All in Market Watch.
* **`invalid login/password/server`** → verify values against MT4 login dialog.
* **Timeouts on first run** → start MT4 manually once so it initializes data; then retry.

---

## 🔗 See also

* `EnsureSymbolVisible.md` — make sure `DefaultSymbol` is available.
* `AccountSummary.md` — quick health snapshot after config is loaded.
* `GetQuote.md` — confirm quotes for the default symbol.
