# Reliability: Timeouts, Reconnects, Backoff (GoMT4)

This chapter shows **how reliability is implemented in this repo** and how to tune it.
Everything below references real code locations so you can cross‑check quickly.

---

## 1) Backoff & Jitter (central knobs)

Defined in `examples/mt4/MT4Account.go`:

```go
// Retry/backoff settings.
const (
    backoffBase = 300 * time.Millisecond // initial backoff
    backoffMax  = 5 * time.Second        // cap for backoff
    jitterRange = 200 * time.Millisecond // ± jitter range
    maxRetries  = 10                     // attempts before giving up
)
```

And the helpers:

```go
// waitWithCtx sleeps for d unless ctx is done.
func waitWithCtx(ctx context.Context, d time.Duration) error { /* ... */ }

// backoffDelay returns exponential backoff with jitter, capped.
func backoffDelay(attempt int) time.Duration { /* base<<attempt, cap, jitter */ }
```

> **Why it matters:** all retry loops (unary *and* streaming) use these values. Increase `backoffBase` for slower retry cadence; raise `backoffMax` for noisy networks; widen `jitterRange` to desync reconnect storms.

---

## 2) Per‑call timeouts (unary RPC)

Pattern used across calls: add a **short per‑call timeout** around read‑only RPCs.
Real example in `ConnectByServerName` health‑check (`examples/mt4/MT4Account.go`):

```go
hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()
if _, err := a.AccountSummary(hctx); err != nil { /* disconnect & error */ }
```

Guidelines:

* Read‑only calls: **2–5s**.
* Trading actions: **3–8s**, depending on broker latency.
* Always `defer cancel()` to avoid goroutine leaks.

---

## 3) Retrying unary calls (transport errors)

Unary execution uses a retry loop that:

* calls the RPC,
* if error is **`codes.Unavailable`** (transient transport), waits `backoffDelay(attempt)` with `waitWithCtx`, then retries,
* aborts on context cancel/deadline.

Pseudo‑excerpt (mirrors `examples/mt4/MT4Account.go` logic):

```go
for attempt := 0; attempt < maxRetries; attempt++ {
    res, err := grpcCall(headers)
    if err == nil { return res, nil }
    if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
        if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil { return zero, err }
        continue
    }
    // non‑transient error: stop
    return zero, err
}
return zero, fmt.Errorf("max retries reached: %w", lastErr)
```

Tuning tips:

* For LAN/VPS, `backoffBase=150ms` often feels snappier.
* For unstable links, keep `base=300ms`, maybe `max=8–10s`.

---

## 4) Streaming with auto‑reconnect

Streaming helpers reopen the stream on recoverable errors (**`io.EOF`/`codes.Unavailable`**), using the same backoff+jitter.
Used by high‑level methods like **`OnSymbolTick`** and history/updates streams (`examples/mt4/MT4Account.go`).

Consumer pattern:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Start stream
dataCh, errCh := account.OnSymbolTick(ctx, []string{"EURUSD"})

for {
    select {
    case <-ctx.Done():
        return             // stops stream + cleanup
    case err := <-errCh:
        if err != nil {    // non‑recoverable error surfaced by helper
            log.Printf("stream stopped: %v", err)
            return
        }
    case tick := <-dataCh:
        // process tick here (send to DB, strategy, etc.)
    }
}
```

Notes:

* The helper **closes both channels** when it gives up (`maxRetries` reached) or `ctx` is canceled.
* You **should not** call `Recv()` yourself; consume from `dataCh`.

---

## 5) Clean cancellation & health checks

* Use **a single parent `ctx`** per workflow and pass it through; cancel on shutdown.
* Main example (`examples/main.go`): the account is closed via `defer account.Disconnect()`.
* After connect, the code performs an **AccountSummary health‑check** with a 3s timeout (see §2) to ensure the terminal is ready.

Shutdown checklist:

* Cancel contexts of long‑living streams first.
* Wait for goroutines to finish (if you used your own workers), then `Disconnect()`.

---

## 6) Choosing sensible defaults

| Scenario                | Suggested timeouts | Backoff (base → max)            |
| ----------------------- | ------------------ | ------------------------------- |
| Read‑only (quotes/info) | 2–3s per call      | 300ms → 5s (default)            |
| Trading actions         | 3–8s per call      | 300ms → 5–8s                    |
| Unstable/home Wi‑Fi     | 4–6s per call      | 500ms → 8–10s, jitter 300–400ms |
| VPS / local LAN         | 1–2s per call      | 150ms → 3–5s                    |

> All these map to the constants at the top of `MT4Account.go`. Adjust them in one place to affect all retry loops.

---

## 7) Common pitfalls (and fixes)

* **Leak: forgot `cancel()`** → Always `defer cancel()` after `WithTimeout/WithCancel`.
* **Hammering retries** → Increase `backoffBase` or `jitterRange`.
* **Permanent errors treated as transient** → Only retry on `codes.Unavailable`/`io.EOF` (transport). Propagate business errors immediately.
* **Dead app on shutdown** → Ensure your select reads `<-ctx.Done()` and returns.

---

## 8) Where to look in code

* Constants & helpers: `examples/mt4/MT4Account.go` (retry/backoff, jitter, `waitWithCtx`).
* Unary patterns & health‑check: `examples/mt4/MT4Account.go` (`context.WithTimeout` around calls).
* Streaming patterns: `examples/mt4/MT4Account.go` (stream wrapper + `OnSymbolTick`).
* Entry point & cleanup: `examples/main.go` (config load, `Disconnect()` on exit).
