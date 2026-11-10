# gRPC Streaming & Subscription Management

> How to subscribe to live data in GoMT4, consume it cleanly, and shut it down — the Go/channel way.

---

## The model in one line

A streaming call returns **two channels**: `(<-chan *pb.…Data, <-chan error)`. Read frames off the data channel with `select`; the error channel delivers a terminal error. **To unsubscribe, cancel the context** — the SDK's background goroutine stops and closes the channels.

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel() // stops the stream

data, errs := acc.OnSymbolTick(ctx, []string{"EURUSD"})
for {
    select {
    case d := <-data:
        t := d.GetSymbolTick()
        fmt.Println(t.GetSymbol(), t.GetBid(), t.GetAsk())
    case err := <-errs:
        if err != nil { return err }
    case <-ctx.Done():
        return nil
    }
}
```

---

## The five streams

On **[`MT4Account`](../MT4Account/Overview.md)**:

| Method | Pushes |
|--------|--------|
| `OnSymbolTick(ctx, symbols)` | live ticks for the given symbols |
| `OnTrade(ctx)` | order/trade events |
| `OnOpenedOrdersTickets(ctx, intervalMs)` | the open-order ticket list on a timer |
| `OnOpenedOrdersProfit(ctx, intervalMs)` | per-order profit + account snapshot on a timer |
| `OnConnectState(ctx, terminalIds)` | terminal connection-state events |

Plus the paged history streams `OrdersHistoryStream` and `QuoteHistoryStream`.

On **[`MT4Service`](../MT4Service/Overview.md)** the same streams are renamed `StreamTicks`, `StreamTrades`, `StreamOpenedTickets`, `StreamOpenedProfit`, `StreamConnectState`.

On **[`MT4Sugar`](../MT4Sugar/Overview.md)** three are exposed as `SubscribeQuotes`, `WatchProfit`, `WatchTickets`.

---

## Bounding a stream

Streams run until you stop them. Common bounds:

### By time

```go
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
data, errs := acc.OnSymbolTick(ctx, []string{"EURUSD"})
// the loop ends when ctx fires
```

### By frame count

```go
n := 0
for {
    select {
    case d := <-data:
        handle(d)
        if n++; n >= 10 { return } // defer cancel() stops the stream
    case err := <-errs:
        if err != nil { return }
    case <-ctx.Done():
        return
    }
}
```

> For "wait until price X" or "wait until filled", `MT4Sugar` already has `WaitPrice` / `WaitFilled` / `WaitClosed`.

---

## Reading several streams at once

Each stream is independent. Run one per goroutine, or multiplex them in a single `select`:

```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

ticks, tickErr := acc.OnSymbolTick(ctx, []string{"EURUSD"})
prof, profErr := acc.OnOpenedOrdersProfit(ctx, 1000)

for {
    select {
    case d := <-ticks:
        t := d.GetSymbolTick(); fmt.Println("tick", t.GetBid())
    case d := <-prof:
        if a := d.GetAccountInfo(); a != nil { fmt.Println("equity", a.GetEquity()) }
    case err := <-tickErr:
        if err != nil { return err }
    case err := <-profErr:
        if err != nil { return err }
    case <-ctx.Done():
        return nil
    }
}
```

> The gRPC client is safe for concurrent use — opening several streams on one `MT4Account` is fine.

---

## Reconnection

Pair streams with `OnConnectState` to know when the broker link drops, and re-open your data streams after a reconnect. The unary `call` wrapper already auto-retries transient transport errors and terminal-not-found; long-lived streams are re-opened by your own logic.

```go
cs, cErr := acc.OnConnectState(ctx, nil)
for {
    select {
    case d := <-cs:
        if !d.GetIsConnected() { /* pause trading; re-open data streams on reconnect */ }
    case err := <-cErr:
        if err != nil { return err }
    case <-ctx.Done():
        return nil
    }
}
```

---

## ✅ Do / ❌ Don't

**✅ Do**
- Cancel the context to unsubscribe — it's deterministic.
- Bound streams (count, or `context.WithTimeout`).
- Handle the **error channel**.
- Keep draining the data channel — do slow work in another goroutine so you don't stall the stream.
- Prefer **one stream with many symbols** over many single-symbol streams.

**❌ Don't**
- Don't forget `defer cancel()` — otherwise the goroutine leaks.
- Don't poll `Quote` in a tight loop when a tick stream exists.
- Don't assume "no frames" means "broken" — a quiet market simply sends nothing.

---

**See also:** [MT4Account overview](../MT4Account/Overview.md) · [Concurrency Model](SYNC_VS_ASYNC.md)
