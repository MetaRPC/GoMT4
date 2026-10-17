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

Every unary method of `MT4Account` adds its **own per‑call timeout when your `ctx` has no deadline**: 3s for read‑only calls (`AccountSummary`, `Quote`, `OpenedOrders`, `SymbolParams`, …), 5s for trading calls (`OrderSend`, `OrderModify`, `OrderClose`, `OrderCloseBy`), 8s for history (`OrdersHistory`, `QuoteHistory`). If you pass a `ctx` with a deadline, yours is used instead.

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

Unary execution (`ExecuteWithReconnect`) uses a retry loop that:

* calls the RPC,
* if the transport error is **`codes.Unavailable`**, waits `backoffDelay(attempt)` with `waitWithCtx`, then retries,
* if the reply carries the API error **`TERMINAL_INSTANCE_NOT_FOUND`** or **`TERMINAL_REGISTRY_TERMINAL_NOT_FOUND`** (the terminal is not up yet / is restarting), waits and retries the same way,
* returns any other error immediately, and aborts on context cancel/deadline.

Excerpt (`examples/mt4/MT4Account.go`, simplified):

```go
for attempt := 0; attempt < maxRetries; attempt++ {
    res, err := grpcCall(a.getHeaders())
    if err != nil {
        if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
            lastErr = err
            if werr := waitWithCtx(ctx, backoffDelay(attempt)); werr != nil { return zeroT, werr }
            continue
        }
        return zeroT, err // non‑transient
    }
    if apiErr := errorSelector(res); apiErr != nil {
        code := apiErr.GetErrorCode()
        if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
            /* wait and retry */
            continue
        }
        return zeroT, wrapAPIError(apiErr)
    }
    return res, nil
}
return zeroT, fmt.Errorf("exceeded retries: %w", lastErr)
```

⚙️ Tuning tips:

* The constants are unexported; to change them, edit `examples/mt4/MT4Account.go` in your clone.
* On a fast, stable link (e.g. a VPS), `backoffBase=150ms` often feels snappier.
* For unstable links, keep `base=300ms`, maybe `max=8–10s`.

---

## 📡 4) Streaming with auto‑reconnect

The streaming helper (`ExecuteStreamWithReconnect`) reopens the stream on recoverable errors (**`io.EOF`, `codes.Unavailable`, and the API errors `TERMINAL_INSTANCE_NOT_FOUND` / `TERMINAL_REGISTRY_TERMINAL_NOT_FOUND`**), using the same backoff+jitter.
Used by the subscriptions **`OnSymbolTick`**, **`OnTrade`**, **`OnOpenedOrdersProfit`**, **`OnOpenedOrdersTickets`** (`examples/mt4/MT4Account.go`). (`OrdersHistoryStream` and `QuoteHistoryStream` are not gRPC streams: they page through unary calls, each of which retries as in §3.)

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
    case err, ok := <-errCh:
        if ok && err != nil {
            log.Printf("stream stopped: %v", err)
        }
        return
    case tick, ok := <-dataCh:
        if !ok {
            return
        }
        _ = tick // process tick (DB, strategy, etc.)
    }
}
```

🔎 Notes:

* The helper sends one error to `errCh` and **closes both channels** when it gives up (`maxRetries` reached), on a non‑recoverable error, or when `ctx` is canceled.
* Don’t call `Recv()` yourself; consume from `dataCh`.

---

## 🧹 5) Clean cancellation & health checks

* Use **one parent `ctx`** per workflow; cancel on shutdown.
* Example (`examples/main.go`): account closed via `defer account.Disconnect(true)` (`true` also deletes the terminal instance on the server).
* After connect, code performs **AccountSummary health‑check** with 3s timeout; if it fails, the account is disconnected and the connect returns an error.

✅ Shutdown checklist:

* Cancel stream contexts first.
* Wait for goroutines to finish.
* Then `Disconnect()` (or `Disconnect(true)`); after it the `MT4Account` cannot be used again — create a new one with `NewMT4Account` to reconnect.

---

## ⚖️ 6) Choosing sensible defaults

| Scenario                | Timeouts      | Backoff (base → max)            |
| ----------------------- | ------------- | ------------------------------- |
| Read‑only (quotes/info) | 2–3s per call | 300ms → 5s (default)            |
| Trading actions         | 3–8s per call | 300ms → 5–8s                    |
| Unstable/home Wi‑Fi     | 4–6s per call | 500ms → 8–10s, jitter 300–400ms |
| Fast VPS                | 1–2s per call | 150ms → 3–5s                    |

📌 Timeouts are set per call through the `ctx` you pass (or the method defaults from §2). Backoff values are the constants in `MT4Account.go`: adjust once, apply everywhere.

---

## ⚠️ 7) Common pitfalls (and fixes)

* ❌ Leak: forgot `cancel()` → ✅ Always `defer cancel()`.
* ❌ Hammering retries → ✅ Increase `backoffBase` or `jitterRange`.
* ❌ Permanent errors treated as transient → ✅ Retry only on `codes.Unavailable`/`io.EOF` and the terminal-not-found API errors (what the built-in helpers do).
* ❌ Dead app on shutdown → ✅ Handle `<-ctx.Done()` properly.

---

## 📂 8) Where to look in code

* Constants & helpers → `examples/mt4/MT4Account.go` (retry/backoff, jitter, `waitWithCtx`).
* Unary patterns & health‑check → `examples/mt4/MT4Account.go` (`ExecuteWithReconnect`, `ConnectByServerName`).
* Streaming patterns → `examples/mt4/MT4Account.go` (`ExecuteStreamWithReconnect`, `OnSymbolTick`).
* Entry point & cleanup → `examples/main.go` (`Disconnect(true)` on exit).
