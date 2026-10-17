# 🧾 ConfigExample (GoMT4)

**Goal:** show the minimal **config.json** used by this repo and how it’s consumed by the examples.

> Real files:
>
> * Config file: `examples/config/config.json`
> * Config loader: `examples/config/config.go` (`LoadConfig` reads JSON into `config.Config`)
> * Used by: `examples/main.go`

---

## 📍 Location

```
examples/config/config.json
```

This file is read at startup by the example app (`config.LoadConfig("config/config.json")`, relative to the `examples` folder you run it from).

---

## 🧩 Schema (what fields mean)

```json
{
  "Login": 12345678,
  "Password": "<your-mt4-password>",
  "Server": "MetaQuotes-Demo",
  "DefaultSymbol": "EURUSD",
  "ApiKey": "<your-mrpc-api-key>"
}
```

* **Login** (number) — MT4 account login. `0` makes `LoadConfig` open a new demo account on `Server` (through `https://mt4.mrpc.pro/DemoAccount/Open`) and fill in `Login`, `Password` and `Server`.
* **Password** (string) — investor or trade password. With the **investor** password you can read but not trade.
* **Server** (string) — exact broker server name as shown in MetaTrader (e.g., `MetaQuotes-Demo`); passed to `ConnectByServerName`.
* **DefaultSymbol** (string) — instrument to use by default (must match broker’s name; suffixes like `EURUSD.m` are different symbols).
* **ApiKey** (string) — your MetaRPC API key. Empty or `"TRIAL"` → `LoadConfig` takes `MRPC_API_KEY` from the environment if set, otherwise the trial is used. `main.go` also accepts `--api-key=...` on the command line.

---

## 🛠️ Edit & run

1. Open `examples/config/config.json` and fill your values.
2. Launch the example:

```powershell
cd GoMT4/examples
go mod download
go run .
```

If the config is valid, you’ll see `✅ Connected to MT4 server` and the output of the example service.

---

## 🔍 Validation (quick checks)

* Login is numeric, password non‑empty (or `Login` is `0` for an auto demo account).
* Server name exactly matches MetaTrader (as in the terminal's login dialog).
* `DefaultSymbol` exists on the broker (`account.Symbols(ctx)` lists the names).

Minimal runtime probe (uses real calls):

```go
cfg, err := config.LoadConfig("config/config.json")
if err != nil { return err }

account, err := mt4.NewMT4Account(uint64(cfg.Login), cfg.Password, "", uuid.Nil, cfg.ApiKey)
if err != nil { return err }
defer account.Disconnect()

if err := account.ConnectByServerName(ctx, cfg.Server, cfg.DefaultSymbol, true, 30); err != nil { return err }
if _, err := account.AccountSummary(ctx); err != nil { return err }
if _, err := account.Quote(ctx, cfg.DefaultSymbol); err != nil { return err }
```

---

## 🔐 Security notes

* `config.json` in this repo is meant for **local development**.
* Do **not** commit real trading credentials or API keys to a public repo. Use demo credentials, or keep the API key in `MRPC_API_KEY`.
* If you need secrets isolation, read credentials from environment variables in your own code (see [Security & Secrets](../../Security_Secrets.md)).

---

## ⚠️ Common errors

* **`symbol not found`** → wrong symbol; check suffix (`EURUSD.m`).
* **connect error** → verify login/password/server against MetaTrader's login dialog.
* **`API_KEY_INVALID`** → the key is unknown; copy it again from [https://mrpc.pro/my](https://mrpc.pro/my).
* **Timeouts on first run** → the first connect starts a terminal for the account; retry, or raise the connect timeout (last argument of `ConnectByServerName`).

---

## 🔗 See also

* [EnsureSymbolVisible](EnsureSymbolVisible.md) — make sure `DefaultSymbol` is available.
* [AccountSummary](../Account/AccountSummary.md) — quick health snapshot after config is loaded.
* [GetQuote](../Market_Info/GetQuote.md) — confirm quotes for the default symbol.
