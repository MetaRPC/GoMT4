# 🔂 UnaryRetries (GoMT4)

**Goal:** apply a consistent retry pattern for **unary** RPC calls (quotes, account, orders) using helpers already present in this repo.

> Real code refs:
>
> * Backoff & helpers: `examples/mt4/MT4Account.go` (`waitWithCtx`, `backoffDelay`, `maxRetries`)
> * Typical calls: `Quote`, `AccountSummary`, `OrderSend`, `OrderModify`, `OrderClose`

---

## ✅ Principles

* Retry **only transient transport** errors: `codes.Unavailable` (and network I/O errors if wrapped accordingly).
* Use **per‑call timeout** via `context.WithTimeout` to bound latency.
* Between attempts, **sleep with backoff + jitter** using `waitWithCtx(ctx, backoffDelay(attempt))`.
* Stop immediately on **non‑transient** (business) errors.

---

## 🧱 Skeleton helper (mirrors your account layer)

```go
func (a *MT4Account) callWithRetry(ctx context.Context, fn func(context.Context) error) error {
    var last error
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := fn(ctx); err != nil {
            if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
                // transient: back off and retry
                if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil {
                    return err // context cancelled/deadline
                }
                last = err
                continue
            }
            return err // non‑transient → bubble up
        }
        return nil // success
    }
    return fmt.Errorf("max retries reached: %w", last)
}
```

> In this project, a similar logic is called directly from the account methods so as not to duplicate backups throughout the project.

---

## ⏱️ Per‑call timeout wrapper

```go
func withTimeout(parent context.Context, d time.Duration, fn func(context.Context) error) error {
    ctx, cancel := context.WithTimeout(parent, d)
    defer cancel()
    return fn(ctx)
}
```

---

## 💱 Example: robust `Quote`

```go
err := withTimeout(ctx, 3*time.Second, func(c context.Context) error {
    return a.callWithRetry(c, func(cc context.Context) error {
        q, err := a.Quote(cc, symbol)
        if err != nil { return err }
        log.Printf("%s %.5f/%.5f @ %s", symbol, q.GetBid(), q.GetAsk(), q.GetTime().AsTime())
        return nil
    })
})
if err != nil { return fmt.Errorf("quote failed: %w", err) }
```

---

## 🧾 Example: `AccountSummary` health‑check

```go
err := withTimeout(ctx, 3*time.Second, func(c context.Context) error {
    return a.callWithRetry(c, func(cc context.Context) error {
        _, err := a.AccountSummary(cc)
        return err
    })
})
if err != nil { return fmt.Errorf("health‑check failed: %w", err) }
```

---

## 🛒 Example: `OrderSend` (market)

```go
err := withTimeout(ctx, 8*time.Second, func(c context.Context) error {
    return a.callWithRetry(c, func(cc context.Context) error {
        _, err := a.OrderSend(cc, symbol, side, volume, nil, &slip, sl, tp, &comment, &magic, nil)
        return err
    })
})
if err != nil { return fmt.Errorf("OrderSend failed: %w", err) }
```

---

## 🎛️ Tuning

Constants in `examples/mt4/MT4Account.go`:

```go
const (
    backoffBase = 300 * time.Millisecond
    backoffMax  = 5 * time.Second
    jitterRange = 200 * time.Millisecond
    maxRetries  = 10
)
```

* **VPS/LAN**: `backoffBase=150ms`, `backoffMax=3–5s`, timeouts 2–3s for reads.
* **Home/unstable**: `backoffMax=8–10s`, timeouts 4–6s (reads) / 6–10s (trades).

---

## ⚠️ Pitfalls

* **We will delete business errors** (for example, invalid volume/price) — no need, return it immediately.
* **There is no `defer cancel()`** — goroutin leaks.
* **Too aggressive backoff** — "pounding" on the network; increase the "backoffBase" and jitter spread.
* **One global context for all** — it is better to have a separate timeout for each call.

---

## 🔗 See also

* `HandleReconnect.md ` — for streaming and general strategy.
* `Reliability (en)' — summary recommendations on timeouts/retreats.
* `GetQuote.md `, `PlaceMarketOrder.md ` — where it is applied live.
