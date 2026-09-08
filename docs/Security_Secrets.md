# 🔐 Security & Secrets (GoMT4)

This page explains **how to handle credentials and protect your setup** based on the current GoMT4 codebase (Windows, local MT4 terminal, gRPC on `mt4.mrpc.pro:443`).

> TL;DR: keep secrets **out of Git**, prefer **per‑machine env vars** or Windows **Credential Manager**, lock down the gRPC port, and redact logs.

---

## 🧭 Scope & threat model

* **Default setup**: gRPC server listens on **localhost only** → attacks must come from the same PC.
* **If you expose it to LAN/WAN** (not default): you must add **TLS** and **authentication** (see ▶ Optional hardening).

---

## 🗝️ Secrets inventory (what we actually store)

* **MT4 login** *(account number)*
* **MT4 password** *(investor or trade)*
* **MT4 server** *(e.g., `RoboForex-Demo`)*
* Optional **symbols list / defaults** (not sensitive)

In your repo these live in `examples/config/config.json` (dev) or **env vars** (recommended).

---

## 📁 Where to store secrets

### 1) Development (recommended)

* Create **`.env`** in repo root (ignored by Git) and let VS Code inject it (via `launch.json` → `envFile`).
* Commit **only** `.env.example` (placeholder values).

**`.env.example`**

```
MT4_LOGIN=12345678
MT4_PASSWORD=replace-me
MT4_SERVER=YourBroker-Server
DEFAULT_SYMBOL=EURUSD
```

**.gitignore** (fragment)

```
# Secrets
.env
examples/config/config.local.json
```

**Why**: the real `.env` never leaves your machine; teammates create their own.

### 2) Production / headless

* Prefer **per‑machine environment variables** (user scope):

  * PowerShell (current user):

    ```powershell
    [Environment]::SetEnvironmentVariable("MT4_PASSWORD","<secret>","User")
    ```
* Or use **Windows Credential Manager** and load at runtime (requires a small helper; see ▶ Optional hardening).

### 3) JSON config

* `examples/config/config.json` is convenient but **do not commit real credentials**. Use a local variant (`config.local.json`) in `.gitignore`.

---

## 🔒 Process & file permissions (Windows)

* Run GoMT4 under a **standard user**, not Administrator.
* Keep repo and MT4 data folders in user profile; avoid world‑writable paths.
* NTFS permissions: only your user needs read access to `.env` and local configs.

---

## 🧱 Network exposure

* Keep the listener on **`mt4.mrpc.pro:443`** for local dev.
* The firewall rule you created opens the port inbound; limit it to **Local subnet** or remove it if not required for local‑only.
* If you must listen on `0.0.0.0` or a LAN IP: enable **TLS** and some **auth** (token or mTLS). See below.

---

## 📝 Logging hygiene

* Never log credentials or full requests. Redact sensitive fields.

**Redaction helper (Go)**

```go
func redact(v string) string {
    if len(v) <= 4 { return "***" }
    return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
// log.Printf("login=%d server=%s pwd=%s", login, server, redact(password))
```

* When printing orders/quotes, it’s fine; avoid dumping entire structs that may include headers/metadata.

---

## 🔗 Dependencies & supply chain

* You import pb as a Go module: `github.com/MetaRPC/GoMT4`.
* **Pin versions**: use **tags/commits** in `go.mod`, keep `go.sum` committed. Example:

  ```
  require github.com/MetaRPC/GoMT4 v0.1.3
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
* [ ] Logs redact passwords/tokens.
* [ ] gRPC bound to **127.0.0.1** unless TLS+auth is configured.
* [ ] `go.mod` pins pb module; `go.sum` committed.

---

## ▶ Optional hardening (when exposing beyond localhost)

### TLS for gRPC

> Not enabled by default in examples. Enable if you bind to non‑localhost.

* Generate a server cert (self‑signed for lab):

  * PowerShell:

    ```powershell
    New-SelfSignedCertificate -DnsName "gomt4.local" -CertStoreLocation Cert:\LocalMachine\My
    ```
  * Export PFX/CRT/KEY and configure your server to use it.
* In Go server, add:

```go
creds := credentials.NewTLS(&tls.Config{ /* MinVersion: tls.VersionTLS12, Certificates: [...] */ })
s := grpc.NewServer(grpc.Creds(creds))
```

* In client, trust the CA or use `RootCAs` with the server cert.

### Simple token auth (metadata header)

> Not present in your current examples. Add only if you need LAN/WAN.

* Client adds header:

```go
md := metadata.Pairs("x-api-key", os.Getenv("GOMT4_API_KEY"))
ctx := metadata.NewOutgoingContext(ctx, md)
```

* Server interceptor checks it:

```go
func apiKeyUnary(next grpc.UnaryHandler) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    if md, ok := metadata.FromIncomingContext(ctx); ok {
      if keys := md.Get("x-api-key"); len(keys) == 1 && keys[0] == os.Getenv("GOMT4_API_KEY") {
        return handler(ctx, req)
      }
    }
    return nil, status.Error(codes.Unauthenticated, "invalid api key")
  }
}
```

### Windows Credential Manager

> Optional replacement for `.env` in production.

* Store a Generic Credential (name `GoMT4/MT4_PASSWORD`).
* Load via a small helper lib (e.g., `github.com/danieljoos/wincred`).

---

## 🔍 Quick self‑audit before push

* Grep for forbidden strings:

  ```powershell
  Select-String -Path . -Pattern "Password=", "MT4_PASSWORD", "x-api-key" -NotMatch "\.env$" -Recurse
  ```
* Double‑check `.gitignore` catches `.env` and any `*.local.json`.

---

## 📎 References

* Cookbook → `ConfigExample.md` (how to structure dev configs)
* Setup → `setup.md` (launch.json with `envFile`)
* Performance Notes → logging hints (hot paths)
