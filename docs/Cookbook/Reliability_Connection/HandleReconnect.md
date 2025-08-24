# 🔁 HandleReconnect (GoMT4)

**Goal:** robust reconnects for **unary** and **streaming** RPCs using the helpers already present in this repo.

> Real code refs:
>
> * Backoff & helpers: `examples/mt4/MT4Account.go` (`backoffDelay`, `waitWithCtx`, `maxRetries`, etc.)
> * Unary pattern: `examples/mt4/MT4Account.go` (retry on `codes.Unavailable`)
> * Streams: `OnSymbolTick`, `OnOpenedOrdersProfit` wrappers

---

## ✅ 1) Principles

* Retry **only transient transport** errors: `codes.Unavailable`, `io.EOF`.
* Respect **context** (timeouts/cancel) to avoid leaks.
* Use **exponential backoff + jitter** (central constants in `MT4Account.go`).

---

## 🔹 2) Unary RPC with built-in retry (pattern)

Most account methods already follow this template: try → on `Unavailable` wait `backoffDelay(attempt)` → retry.

```go
func (a *MT4Account) callWithRetry(ctx context.Context, fn func(context.Context) error) error {
    var last error
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := fn(ctx); err != nil {
            st, ok := status.FromError(err)
            if ok && st.Code() == codes.Unavailable {
                // transient transport → backoff, then retry
                if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil { return err }
                last = err
                continue
            }
            return err // non-transient → bubble up
        }
        return nil // success
    }
    return fmt.Errorf("max retries reached: %w", last)
}
```

**Usage (example: health-check AccountSummary):**

```go
hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()
err := a.callWithRetry(hctx, func(c context.Context) error {
    _, err := a.AccountSummary(c)
    return err
})
if err != nil { return err }
```

> The same pattern is used internally by methods like `OrderSend`, `Quote`, etc., in your account layer.

---

## 🔸 3) Streaming reconnect loop (structure)

Your stream helpers (`OnSymbolTick`, `OnOpenedOrdersProfit`) already encapsulate the reconnect loop. The core logic looks like this:

```go
func (a *MT4Account) runStreamWithReconnect(ctx context.Context, start func(context.Context) (recv func() (*pb.Tick, error), close func() error, err error),
) (<-chan *pb.Tick, <-chan error) {
    dataCh := make(chan *pb.Tick, 1024)
    errCh  := make(chan error, 1)

    go func() {
        defer close(dataCh)
        defer close(errCh)

        for attempt := 0; attempt < maxRetries; attempt++ {
            // (re)open stream
            recv, closeFn, err := start(ctx)
            if err != nil {
                // cannot open → transient?
                if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
                    if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil { errCh <- err; return }
                    continue
                }
                errCh <- err; return
            }

            // receive loop
            for {
                msg, err := recv()
                if err == nil {
                    select {
                    case dataCh <- msg:
                    case <-ctx.Done(): _ = closeFn(); return
                    }
                    continue
                }
                // stream error → decide if reconnect
                if err == io.EOF {
                    // server closed → reconnect with backoff
                } else if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
                    // transient transport → reconnect
                } else {
                    // permanent
                    _ = closeFn(); errCh <- err; return
                }
                _ = closeFn()
                if err := waitWithCtx(ctx, backoffDelay(attempt)); err != nil { errCh <- err; return }
                break // out to reopen
            }
        }
        errCh <- fmt.Errorf("max stream retries reached")
    }()

    return dataCh, errCh
}
```

> In your repo, this logic is packaged in concrete helpers: `OnSymbolTick(ctx, symbols)`, `OnOpenedOrdersProfit(ctx, bufSize)`.

---

## ▶️ 4) Consumer pattern (don’t block!)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Example: quotes stream
dataCh, errCh := account.OnSymbolTick(ctx, []string{"EURUSD","GBPUSD"})

for {
    select {
    case <-ctx.Done():
        return // graceful stop
    case err := <-errCh:
        if err != nil { log.Printf("stream stopped: %v", err); return }
    case t := <-dataCh:
        // offload heavy work
        processAsync(t)
    }
}
```

* Heavy work → отправляй в воркер через буферизованный канал.
* Не забывай `ctx.Done()` для чистого завершения.

---

## 🧭 5) Tuning backoff (central knobs)

Constants found in `examples/mt4/MT4Account.go`:

```go
const (
    backoffBase = 300 * time.Millisecond
    backoffMax  = 5 * time.Second
    jitterRange = 200 * time.Millisecond
    maxRetries  = 10
)
```

* **Home Wi‑Fi / unstable** → попробуй `backoffMax=8–10s`, `jitterRange=300–400ms`.
* **VPS / LAN** → `backoffBase=150ms`, `backoffMax=3–5s`.

---

## ⚠️ Pitfalls

* **Ретраим бизнес-ошибки** → нельзя. Ретраим только транспортные (`Unavailable`, `EOF`).
* **Забыли отменить контекст** → утечки горутин. Всегда `defer cancel()`.
* **Блокируем dataCh** → поток встанет. Либо буфер, либо быстрый приём.
* **Бесконечные ретраи** → ограничивай `maxRetries`, логируй финальную ошибку.

---

## 🔗 See also

* `Reliability (en)` — timeouts, reconnects & backoff summary.
* `StreamQuotes.md`, `StreamOpenedOrderProfits.md` — готовые обёртки.
* `UnaryRetries.md` — точечные примеры для отдельных методов.
