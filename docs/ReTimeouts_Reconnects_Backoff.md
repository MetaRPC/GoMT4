# 🔒 Reliability: Timeouts, Reconnects, Backoff (GoMT4)

This chapter shows **how reliability is implemented in this repo** and how to tune it.
Everything below references real code locations so you can cross‑check quickly.

---

## ⏳ 1) Backoff & Jitter (central knobs)

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

Helpers:

```go
// waitWithCtx sleeps for d unless ctx is done.
func waitWithCtx(ctx context.Context, d time.Duration) error { /* ... */ }

// backoffDelay returns exponential backoff with jitter, capped.
func backoffDelay(attempt int) time.Duration { /* base<<attempt, cap, jitter */ }
```

💡 **Why it matters:** all retry loops (unary *and* streaming) use these values.
Increase `backoffBase` for slower retry cadence; raise `backoffMax` for noisy networks; widen `jitterRange` to desync reconnect storms.

---

## ⏱️ 2) Per‑call timeouts (unary RPC)

Pattern across calls: add a **short per‑call timeout** around read‑only RPCs.
Example in `ConnectByServerName` health‑check (`examples/mt4/MT4Account.go`):

```go
hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()
if _, err := a.AccountSummary(hctx); err != nil { /* disconnect & error */ }
```

📌 Guidelines:

* Read‑only calls: **2–5s**.
* Trading actions: **3–8s**, depending on broker latency.
* Always `defer cancel()` to avoid goroutine leaks.

---

## 🔁 3) Retrying unary calls (transport errors)

Unary execution uses a retry loop that:

* calls the RPC,
* if error is **`codes.Unavailable`** (transient transport), waits `backoffDelay(attempt)` with `waitWithCtx`, then retries,
* aborts on context cancel/deadline.

Pseudo‑excerpt (`examples/mt4/MT4Account.go` logic):

```go
for attempt := 0; attempt < maxRetries; attempt++ {
    res, err := grpcCall(headers)
    if err == nil { return res, nil }
    if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
        if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil { return zero, err }
        continue
    }
    return zero, err // non‑transient
}
return zero, fmt.Errorf("max retries reached: %w", lastErr)
```

⚙️ Tuning tips:

* For LAN/VPS, `backoffBase=150ms` often feels snappier.
* For unstable links, keep `base=300ms`, maybe `max=8–10s`.

---

## 📡 4) Streaming with auto‑reconnect

Streaming helpers reopen the stream on recoverable errors (**`io.EOF` / `codes.Unavailable`**), using the same backoff+jitter.
Used by methods like **`OnSymbolTick`** and history/updates streams (`examples/mt4/MT4Account.go`).

Consumer pattern:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Start stream
dataCh, errCh := account.OnSymbolTick(ctx, []string{"EURUSD"})

for {
    select {
    case <-ctx.Done():
        return
    case err := <-errCh:
        if err != nil {
            log.Printf("stream stopped: %v", err)
            return
        }
    case tick := <-dataCh:
        // process tick (DB, strategy, etc.)
    }
}
```

🔎 Notes:

* Helper **closes both channels** when it gives up (`maxRetries` reached) or `ctx` canceled.
* Don’t call `Recv()` yourself; consume from `dataCh`.

---

## 🧹 5) Clean cancellation & health checks

* Use **one parent `ctx`** per workflow; cancel on shutdown.
* Example (`examples/main.go`): account closed via `defer account.Disconnect()`.
* After connect, code performs **AccountSummary health‑check** with 3s timeout.

✅ Shutdown checklist:

* Cancel stream contexts first.
* Wait for goroutines to finish.
* Then `Disconnect()`.

---

## ⚖️ 6) Choosing sensible defaults

| Scenario                | Timeouts      | Backoff (base → max)            |
| ----------------------- | ------------- | ------------------------------- |
| Read‑only (quotes/info) | 2–3s per call | 300ms → 5s (default)            |
| Trading actions         | 3–8s per call | 300ms → 5–8s                    |
| Unstable/home Wi‑Fi     | 4–6s per call | 500ms → 8–10s, jitter 300–400ms |
| VPS / local LAN         | 1–2s per call | 150ms → 3–5s                    |

📌 All map to constants in `MT4Account.go`. Adjust once, apply everywhere.

---

## ⚠️ 7) Common pitfalls (and fixes)

* ❌ Leak: forgot `cancel()` → ✅ Always `defer cancel()`.
* ❌ Hammering retries → ✅ Increase `backoffBase` or `jitterRange`.
* ❌ Permanent errors treated as transient → ✅ Retry only on `codes.Unavailable`/`io.EOF`.
* ❌ Dead app on shutdown → ✅ Handle `<-ctx.Done()` properly.

---

## 📂 8) Where to look in code

* Constants & helpers → `examples/mt4/MT4Account.go` (retry/backoff, jitter, `waitWithCtx`).
* Unary patterns & health‑check → `examples/mt4/MT4Account.go`.
* Streaming patterns → `examples/mt4/MT4Account.go` (`OnSymbolTick`).
* Entry point & cleanup → `examples/main.go` (`Disconnect()` on exit).
