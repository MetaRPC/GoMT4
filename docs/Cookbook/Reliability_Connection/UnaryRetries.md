# 🔂 UnaryRetries (GoMT4)

**Goal:** apply a consistent retry pattern for **unary** RPC calls (quotes, account, orders) — on top of what the client already does.

> Real code refs:
>
> * Built-in retry: `examples/mt4/MT4Account.go` (`ExecuteWithReconnect`, with the unexported `waitWithCtx`, `backoffDelay`, `maxRetries`)
> * Typical calls: `Quote`, `AccountSummary`, `OrderSend`, `OrderModify`, `OrderClose`

---

## ✅ Principles

* Every `MT4Account` method **already retries** transient errors: `codes.Unavailable` and the API errors `TERMINAL_INSTANCE_NOT_FOUND` / `TERMINAL_REGISTRY_TERMINAL_NOT_FOUND`, up to 10 attempts with backoff + jitter, within the call's context.
* Each method also bounds itself with a **default timeout** when your `ctx` has no deadline (3s reads, 5s trading, 8s history).
* Extra retries in your code are only useful for longer outages (e.g. the terminal restarting for longer than the built-in attempts take) — and only for **idempotent** calls (reads). Do not blindly retry `OrderSend`: a request that timed out may still have opened an order.
* Stop immediately on **business** errors (invalid volume/price, not enough money, …).

---

## 🧱 Skeleton helper (in your own package)

The client's `waitWithCtx` / `backoffDelay` are unexported, so write your own small helper:

```go
func callWithRetry(ctx context.Context, attempts int, fn func(context.Context) error) error {
    var last error
    for attempt := 0; attempt < attempts; attempt++ {
        err := fn(ctx)
        if err == nil {
            return nil // success
        }
        last = err
        if st, ok := status.FromError(err); !ok || st.Code() != codes.Unavailable {
            if !strings.Contains(err.Error(), "exceeded retries") {
                return err // non‑transient → bubble up
            }
        }
        // transient: back off and retry
        delay := time.Duration(1<<attempt) * 500 * time.Millisecond
        select {
        case <-ctx.Done():
            return ctx.Err() // context cancelled/deadline
        case <-time.After(delay):
        }
    }
    return fmt.Errorf("max retries reached: %w", last)
}
```

> `"exceeded retries"` is the prefix of the error `ExecuteWithReconnect` returns when its own attempts ran out on transient errors.

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
err := callWithRetry(ctx, 3, func(c context.Context) error {
    return withTimeout(c, 3*time.Second, func(cc context.Context) error {
        q, err := account.Quote(cc, symbol)
        if err != nil { return err }
        log.Printf("%s %.5f/%.5f @ %s", symbol, q.GetBid(), q.GetAsk(), q.GetDateTime().AsTime())
        return nil
    })
})
if err != nil { return fmt.Errorf("quote failed: %w", err) }
```

---

## 🧾 Example: `AccountSummary` health‑check

```go
err := withTimeout(ctx, 3*time.Second, func(c context.Context) error {
    _, err := account.AccountSummary(c) // already retries inside
    return err
})
if err != nil { return fmt.Errorf("health‑check failed: %w", err) }
```

---

## 🛒 Example: `OrderSend` (market) — no extra retry

```go
err := withTimeout(ctx, 8*time.Second, func(c context.Context) error {
    _, err := account.OrderSend(c, symbol, side, volume, nil, &slip, sl, tp, &comment, &magic, nil)
    return err
})
if err != nil {
    // Before sending again, check account.OpenedOrders for an order with your magic number/comment.
    return fmt.Errorf("OrderSend failed: %w", err)
}
```

---

## 🎛️ Tuning

Constants in `examples/mt4/MT4Account.go` (unexported — change them in your clone):

```go
const (
    backoffBase = 300 * time.Millisecond
    backoffMax  = 5 * time.Second
    jitterRange = 200 * time.Millisecond
    maxRetries  = 10
)
```

* **Fast VPS**: `backoffBase=150ms`, `backoffMax=3–5s`, timeouts 2–3s for reads.
* **Home/unstable**: `backoffMax=8–10s`, timeouts 4–6s (reads) / 6–10s (trades).

---

## ⚠️ Pitfalls

* **Retrying business errors** (for example, invalid volume/price) — no need, return it immediately.
* **No `defer cancel()`** — goroutine leaks.
* **Too aggressive backoff** — "pounding" the network; increase the delay and add jitter.
* **One global context for all** — it is better to have a separate timeout for each call.

---

## 🔗 See also

* [HandleReconnect](HandleReconnect.md) — for streaming and general strategy.
* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — summary recommendations on timeouts/retries.
* [GetQuote](../Market_Info/GetQuote.md), [PlaceMarketOrder](../Orders/PlaceMarketOrder.md) — where it is applied live.
