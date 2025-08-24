# 🩺 HealthCheck (GoMT4)

**Goal:** perform a fast **terminal health check** right after connect — using the same calls and patterns that exist in this repo.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`AccountSummary`, `Quote`, retry/backoff helpers)
> * Service demo: `examples/mt4/MT4_service.go` (summary/quote examples)
> * Config: `examples/config/config.json`

---

## ✅ What we consider "healthy"

* MT4 terminal is running and **connected to broker**.
* `AccountSummary` returns within a **short timeout** (≈3s).
* Optionally: a `Quote(symbol)` succeeds for your `DefaultSymbol` from `config.json`.

---

## ⏱️ Quick summary-based check (3s timeout + retry on transport)

```go
func HealthCheck(ctx context.Context, a *MT4Account) error {
    // Short deadline: we want a quick yes/no
    hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // Retry only transient transport errors (codes.Unavailable)
    return a.callWithRetry(hctx, func(c context.Context) error {
        _, err := a.AccountSummary(c)
        return err
    })
}
```

**Why this works:** `AccountSummary` is lightweight and exercises the session; your `callWithRetry` already implements exponential backoff + jitter and respects context.

---

## 💱 Optional: include a quote probe for DefaultSymbol

```go
func HealthCheckWithQuote(ctx context.Context, a *MT4Account, symbol string) error {
    if err := HealthCheck(ctx, a); err != nil { return err }

    qctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    return a.callWithRetry(qctx, func(c context.Context) error {
        _, err := a.Quote(c, symbol)
        return err
    })
}
```

Use the `DefaultSymbol` from `examples/config/config.json`, making sure the symbol is visible in *Market Watch* (suffixes like `EURUSD.m` are broker‑specific).

---

## 🧪 Interpreting results

* **OK:** both calls return under the deadlines.
* **Timeout (`context deadline exceeded`):** terminal is not ready or network stalls — wait a bit and retry; consider higher `backoffMax`.
* **Business error:** bubble it up (do **not** retry): fix login/server/symbol.

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

Tune them for your environment (home Wi‑Fi vs VPS/LAN). Timeouts in the health‑check are **per‑call** and independent from backoff caps.

---

## 🧰 Usage example (service layer)

⚠️ Note: Your base project does not include a method named `EnsureHealthy`. This section is shown only as a **convenience wrapper idea** — a way to combine the health‑check calls (`AccountSummary` + `Quote`) into one place. You may add such a method if you want a single entry point to verify that MT4 is ready.

Why useful?

* Before running strategies or bots, you can call `EnsureHealthy` to quickly confirm the terminal is connected and symbols are available.
* It saves copy‑pasting the same two checks (`AccountSummary` + `Quote`) everywhere.
* If something is wrong (wrong server, hidden symbol, no connection) you fail fast with a clear error.

Example implementation:

```go
func (s *MT4Service) EnsureHealthy(ctx context.Context) error {
    // 1) Summary probe
    if err := HealthCheck(ctx, s.account); err != nil {
        return fmt.Errorf("health summary failed: %w", err)")
    }

    // 2) Quote probe for default symbol
    sym := s.cfg.DefaultSymbol
    if sym != "" {
        if err := HealthCheckWithQuote(ctx, s.account, sym); err != nil {
            return fmt.Errorf("health quote failed for %s: %w", sym, err)
        }
    }

    log.Println("✅ MT4 is healthy and ready")
    return nil
}
```

---

## ⚠️ Pitfalls

* **Hidden symbol / suffix mismatch** → `EURUSD` vs `EURUSD.m`; fix config or show all symbols in MT4.
* **No cancel** → always `defer cancel()` after `WithTimeout` to avoid leaks.
* **Retrying business errors** → don’t; adjust credentials/server/symbol.

---

## 📎 See also

* `UnaryRetries.md` — per‑call retry wrapper used above.
* `AccountSummary.md` — details on fields and risk checks.
* `GetQuote.md` — one‑shot quote used in the probe.
