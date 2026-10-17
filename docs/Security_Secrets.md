# 🔐 Security & Secrets (GoMT4)

This page explains **how to handle credentials and protect your setup** based on the current GoMT4 codebase (your Go program → gRPC over TLS → MetaRPC cloud API at `mt4.mrpc.pro:443` → MT4 terminal managed by MetaRPC).

> TL;DR: keep secrets **out of Git**, prefer **per‑machine env vars** or Windows **Credential Manager**, never log the password or API key, and redact logs.

---

## 🧭 Scope & threat model

* **Transport**: every call goes to `mt4.mrpc.pro:443` over **TLS** (`NewMT4Account` uses `credentials.NewTLS` with certificate verification on). Nothing listens on your machine; there is no port to open or protect.
* **What leaves your machine**: the MT4 login/password (in `ConnectEx` / `Connect`, so MetaRPC can log the terminal in), your **API key** (`apikey` metadata on every call) and the **terminal id** (`id` metadata).
* **What you must protect locally**: the files and environment variables that hold those values, and your logs.

---

## 🗝️ Secrets inventory (what we actually store)

* **MT4 login** *(account number)*
* **MT4 password** *(investor or trade)*
* **MT4 server name** *(e.g., `MetaQuotes-Demo`)* — not secret
* **MetaRPC API key** — identifies your MetaRPC account; calls with it run under your account
* Optional **default symbol** (not sensitive)

The example reads them from `examples/config/config.json` (`Login`, `Password`, `Server`, `DefaultSymbol`, `ApiKey`); the API key can also come from the `MRPC_API_KEY` environment variable or the `--api-key` argument.

---

## 📁 Where to store secrets

### 1) Development (recommended)

* Put the **API key** into the `MRPC_API_KEY` environment variable and leave `"ApiKey": "TRIAL"` in `config.json`: `LoadConfig` then takes the key from the environment.
* For login/password in **your own program**, read environment variables (as [Your First Project](All_Guides/Your_First_Project.md) does with `MT4_LOGIN`, `MT4_PASSWORD`, `MT4_SERVER`, `MRPC_API_KEY`). You can keep them in a local **`.env`** file (ignored by Git) and let VS Code inject it (via `launch.json` → `envFile`).
* Commit **only** `.env.example` (placeholder values).

**`.env.example`**

```
MT4_LOGIN=12345678
MT4_PASSWORD=replace-me
MT4_SERVER=YourBroker-Server
MRPC_API_KEY=replace-me
```

**.gitignore** (fragment)

```
# Secrets
.env
```

**Why**: the real `.env` never leaves your machine; teammates create their own.

### 2) Production / headless

* Prefer **per‑machine environment variables** (user scope):

  * PowerShell (current user):

    ```powershell
    [Environment]::SetEnvironmentVariable("MRPC_API_KEY","<secret>","User")
    ```
* Or use **Windows Credential Manager** and load at runtime (requires a small helper; see ▶ Optional hardening).

### 3) JSON config

* `examples/config/config.json` is tracked by Git and convenient, but **do not commit real credentials**. Keep it with placeholders (or `"Login": 0`, which opens a demo account) and tell Git to ignore your local edits:

  ```bash
  git update-index --skip-worktree examples/config/config.json
  ```

---

## 🔒 Process & file permissions (Windows)

* Run GoMT4 under a **standard user**, not Administrator.
* Keep the repo in your user profile; avoid world‑writable paths.
* NTFS permissions: only your user needs read access to `.env` and local configs.

---

## 🧱 Network exposure

* The client only makes **outbound** TLS connections to `mt4.mrpc.pro:443`; no inbound firewall rule is needed.
* Do not disable certificate verification (`InsecureSkipVerify`) — the client keeps it on.

---

## 📝 Logging hygiene

* Never log credentials, the API key or full requests. Redact sensitive fields.
* Treat the **terminal id** (`account.TerminalInstanceGuid`, `account.Id`) and the **API key** as private: calls carrying them act on your terminal.

**Redaction helper (Go)**

```go
func redact(v string) string {
    if len(v) <= 4 { return "***" }
    return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
// log.Printf("login=%d server=%s pwd=%s key=%s", login, server, redact(password), redact(apiKey))
```

* When printing orders/quotes, it’s fine; avoid dumping the `MT4Account` struct — it contains `Password` and `ApiKey`.

---

## 🔗 Dependencies & supply chain

* You import pb as a Go module: `github.com/MetaRPC/GoMT4/package` (the client in `examples/mt4` comes from your clone of the repository).
* **Pin versions**: use **tags/commits** in `go.mod`, keep `go.sum` committed. Example (pseudo‑version of a commit):

  ```
  require github.com/MetaRPC/GoMT4/package v0.0.0-<date>-<commit>
  ```
* For reproducible/offline builds you can `go mod vendor`; this copies deps into `vendor/` (bigger repo, but no network at build time).

---

## 🧪 Secrets in tests & examples

* Don’t embed real credentials in `main.go` or examples; read from env/config.
* Add a quick **startup guard**:

```go
if os.Getenv("MT4_PASSWORD") == "replace-me" {
    log.Fatal("Refusing to start with placeholder password — set MT4_PASSWORD")
}
```

---

## ✅ Checklist (quick)

* [ ] `.env` exists locally; `.env.example` in Git; **real `.env` is ignored**.
* [ ] No real secrets in `config.json` committed.
* [ ] Logs redact passwords, API keys and terminal ids.
* [ ] Certificate verification stays on.
* [ ] `go.mod` pins the pb module; `go.sum` committed.

---

## ▶ Optional hardening

### API key handling

* Use a separate API key per deployment, so one can be replaced without touching the others (manage keys at [https://mrpc.pro/my](https://mrpc.pro/my)).
* If a key leaks, replace it in the portal and update `MRPC_API_KEY`; calls with an unknown key are refused (`API_KEY_INVALID`).

### Investor password for read‑only tools

* Dashboards and reports that never trade can log in with the **investor password**: they can read account data, quotes and history, but cannot trade (`AccountSummary` reports `IsInvestor`).

### Windows Credential Manager

> Optional replacement for `.env` in production.

* Store a Generic Credential (name `GoMT4/MT4_PASSWORD`).
* Load via a small helper lib (e.g., `github.com/danieljoos/wincred`).

---

## 🔍 Quick self‑audit before push

* Grep for forbidden strings:

  ```powershell
  Select-String -Path (Get-ChildItem -Recurse -File -Exclude .env) -Pattern '"Password":\s*"[^"<]', "MRPC_API_KEY=", "MT4_PASSWORD="
  ```
* Double‑check `.gitignore` catches `.env`.

---

## 📎 References

* Cookbook → [Config Example](Cookbook/Utils_Helpers/ConfigExample.md) (how to structure dev configs)
* Setup → [setup.md](setup.md) (launch.json)
* Performance Notes → logging hints (hot paths)
