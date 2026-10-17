# 🩺 HealthCheck (GoMT4)

**Goal:** perform a fast **terminal health check** right after connect — using the same calls and patterns that exist in this repo.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`AccountSummary`, `Quote`, the health-check inside `ConnectByServerName` / `ConnectByHostPort`)
> * Service demo: `examples/mt4/MT4_service.go` (summary/quote examples)
> * Config: `examples/config/config.json`

---

## ✅ What we consider "healthy"

* The MT4 terminal MetaRPC runs for your account is up and **connected to the broker**.
* `AccountSummary` returns within a **short timeout** (≈3s).
* Optionally: a `Quote(symbol)` succeeds for your `DefaultSymbol` from `config.json`.

`ConnectByServerName` / `ConnectByHostPort` already do the first check: right after connecting they call `AccountSummary` with a 3s timeout and, if it fails, disconnect and return an error. The helpers below repeat it later (e.g. before each trading session).

---

## ⏱️ Quick summary-based check (3s timeout)

```go
func HealthCheck(ctx context.Context, a *mt4.MT4Account) error {
    // Short deadline: we want a quick yes/no
    hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // AccountSummary retries transient errors itself (within hctx)
    _, err := a.AccountSummary(hctx)
    return err
}
```

**Why this works:** `AccountSummary` is lightweight and exercises the whole path (API key, terminal id, terminal, broker connection).

---

## 💱 Optional: include a quote probe for DefaultSymbol

```go
func HealthCheckWithQuote(ctx context.Context, a *mt4.MT4Account, symbol string) error {
    if err := HealthCheck(ctx, a); err != nil { return err }

    qctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    _, err := a.Quote(qctx, symbol)
    return err
}
```

Use the `DefaultSymbol` from `examples/config/config.json` (suffixes like `EURUSD.m` are broker‑specific).

---

## 🔎 Optional: ask the server for the terminal state

The generated `ConnectionClient` has `ConnectionStatus` (not wrapped by `MT4Account`). Call it through `ExecuteWithReconnect` so the `apikey` / `id` metadata is attached:

```go
reply, err := mt4.ExecuteWithReconnect(account, ctx,
    func(md metadata.MD) (*pb.ConnectionStatusReply, error) {
        return account.ConnectionClient.ConnectionStatus(metadata.NewOutgoingContext(ctx, md), &pb.ConnectionStatusRequest{})
    },
    func(r *pb.ConnectionStatusReply) *pb.Error { return r.GetError() },
)
if err != nil { return err }
d := reply.GetData()
fmt.Printf("state=%s connected=%v alive=%v server=%s\n", d.GetState(), d.GetIsConnected(), d.GetIsAlive(), d.GetServer())
```

---

## 🧪 Interpreting results

* **OK:** the calls return under the deadlines.
* **Timeout (`context deadline exceeded`) / `exceeded retries`:** the terminal is not ready or not connected under this account's id — wait a bit and retry, or connect again (see [HandleReconnect](HandleReconnect.md)).
* **Business error** (wrong symbol, `API_KEY_INVALID`, …): bubble it up (do **not** retry): fix the API key/login/server/symbol.

---

## 🧭 Where the knobs live (no magic numbers)

Backoff/jitter/retry limits are centralized in `examples/mt4/MT4Account.go`:

```go
const (
    backoffBase = 300 * time.Millisecond
    backoffMax  = 5 * time.Second
    jitterRange = 200 * time.Millisecond
    maxRetries  = 10
)
```

Tune them for your environment (home Wi‑Fi vs VPS). Timeouts in the health‑check are **per‑call** and independent from backoff caps.

---

## 🧰 Usage example (your own wrapper)

⚠️ Note: the project does not include a method named `EnsureHealthy`. This section shows a **convenience wrapper idea** — a way to combine the health‑check calls (`AccountSummary` + `Quote`) into one place in your own code.

Why useful?

* Before running strategies or bots, you can call `EnsureHealthy` to quickly confirm the terminal is connected and symbols are available.
* It saves copy‑pasting the same two checks (`AccountSummary` + `Quote`) everywhere.
* If something is wrong (wrong server, wrong symbol, no connection) you fail fast with a clear error.

Example implementation:

```go
func EnsureHealthy(ctx context.Context, account *mt4.MT4Account, defaultSymbol string) error {
    // 1) Summary probe
    if err := HealthCheck(ctx, account); err != nil {
        return fmt.Errorf("health summary failed: %w", err)
    }

    // 2) Quote probe for default symbol
    if defaultSymbol != "" {
        if err := HealthCheckWithQuote(ctx, account, defaultSymbol); err != nil {
            return fmt.Errorf("health quote failed for %s: %w", defaultSymbol, err)
        }
    }

    log.Println("✅ MT4 is healthy and ready")
    return nil
}
```

---

## ⚠️ Pitfalls

* **Suffix mismatch** → `EURUSD` vs `EURUSD.m`; fix config (check names with `account.Symbols(ctx)`).
* **No cancel** → always `defer cancel()` after `WithTimeout` to avoid leaks.
* **Retrying business errors** → don’t; adjust credentials/server/symbol/API key.

---

## 📎 See also

* [UnaryRetries](UnaryRetries.md) — per‑call retry wrapper.
* [AccountSummary](../Account/AccountSummary.md) — details on fields and risk checks.
* [GetQuote](../Market_Info/GetQuote.md) — one‑shot quote used in the probe.
