# The Concurrency Model — How GoMT4 Executes

> Unlike SDKs in `async`/`await` languages, GoMT4 has **no async keyword**. Unary calls are ordinary blocking functions that take a `context.Context`; streams hand you **Go channels** you read in a loop or a goroutine. This page explains the model.

---

## 🎯 Two shapes of call

| Shape | What it is | How you consume it |
|-------|-----------|--------------------|
| **Unary (request / response)** | `AccountSummary(ctx)`, `Quote(ctx, sym)`, `OrderSend(…)` | call it; it blocks until the reply (or the context deadline) |
| **Server stream** | `OnSymbolTick(ctx, …)`, `StreamOpenedProfit(…)` | returns `(<-chan *pb.…Data, <-chan error)` — read frames off the channel |

### Unary

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

q, err := sugar.LastQuote(ctx, "EURUSD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(q.GetBid(), q.GetAsk())
```

The call runs on the calling goroutine and returns `(value, error)`. There is nothing to `await`; the gRPC round-trip happens inside the function.

### Streaming

A stream returns two channels — data and error. Read them with `select`, and stop by cancelling the context:

```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

data, errs := acc.OnSymbolTick(ctx, []string{"EURUSD"})
for {
    select {
    case d := <-data:
        t := d.GetSymbolTick()
        fmt.Println(t.GetSymbol(), t.GetBid(), t.GetAsk())
    case err := <-errs:
        if err != nil { return err }
    case <-ctx.Done():
        return nil // deadline/cancel → stop
    }
}
```

The SDK runs the gRPC receive loop on a background goroutine and forwards frames to the channel, retrying transient drops. See **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)**.

---

## 🧵 `context.Context` is the control knob

Every SDK method takes a `ctx`. It carries the **deadline** and the **cancel** signal:

- `context.WithTimeout(parent, d)` — bound a call; it errors with `DeadlineExceeded` if it overruns.
- `context.WithCancel(parent)` — stop a stream (or a batch of calls) on demand by calling `cancel()`.
- Passing `nil` is accepted — the low-level layer substitutes `context.Background()`. When the context carries no deadline, the method adds a short default per-call timeout: **~3 s** for simple reads (quote, account, opened orders), **~5 s** for trade calls and stats, and **~8 s** for `OrdersHistory`. Pass your own `ctx` with a deadline to override it.

---

## 🚀 Running work concurrently

Go's concurrency is **goroutines + channels**, not futures. Fan out independent calls yourself:

```go
// price several symbols at once
type res struct{ sym string; q *pb.QuoteData; err error }
out := make(chan res, 3)
for _, s := range []string{"EURUSD", "GBPUSD", "USDJPY"} {
    go func(s string) {
        q, err := sugar.LastQuote(ctx, s)
        out <- res{s, q, err}
    }(s)
}
for i := 0; i < 3; i++ {
    r := <-out
    // handle r
}
```

Or run a management loop in the background while the main flow continues:

```go
go func() {
    for {
        if o, _ := sugar.FindTicket(ctx, ticket); o == nil { return } // closed
        sugar.TrailStop(ctx, ticket, 15)
        time.Sleep(2 * time.Second)
    }
}()
```

> Prefer a **batched call** (`QuoteMany`, one `OnSymbolTick` with all symbols) over many single calls where the API offers one.

---

## 💡 Rules of thumb

* **Always pass a `ctx`** with a deadline for unary calls; cancel it to stop streams.
* **One stream, many symbols** beats many single-symbol streams.
* **Issuing RPCs concurrently is fine once connected** — the underlying gRPC calls run in parallel safely. But `MT4Account` has **no internal locking around its connection**, so run `Connect` / `Reconnect` / `Disconnect` from a single goroutine, never while other calls are in flight.
* **Don't block a stream's reader** — do slow work in a separate goroutine so you keep draining the channel.

---

**See also:** [gRPC Streaming](GRPC_STREAM_MANAGEMENT.md) · [MT4Sugar overview](../MT4Sugar/Overview.md)
