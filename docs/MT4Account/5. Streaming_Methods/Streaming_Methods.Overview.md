# 5. Streaming Methods — Overview

> Long-lived subscriptions. Each returns **two channels** — `(<-chan *pb.…Data, <-chan error)` — and runs until you cancel the context.

Unlike the request/response methods in groups 1–4, these open a server stream and keep pushing. The Go pattern is always the same:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()                       // cancel() == unsubscribe

data, errs := acc.OnSymbolTick(ctx, []string{"EURUSD"})
for {
    select {
    case d := <-data:
        // handle frame
    case err := <-errs:
        if err != nil { return err } // stream ended / fatal
    case <-ctx.Done():
        return nil
    }
}
```

> 🔁 **Auto-reconnect.** The `On*` subscriptions go through `ExecuteStreamWithReconnect`, which transparently re-dials on `codes.Unavailable`. **`cancel()` is how you stop** — a cancelled context closes both channels.

> ℹ️ The two `*Stream` helpers (`OrdersHistoryStream`, `QuoteHistoryStream`) are **not** live server streams — they page/chunk a request method and feed the results through a channel for convenience.

---

## Methods in this group

| Method | Frame type | One-liner |
|--------|-----------|-----------|
| [`OnSymbolTick(ctx, symbols)`](./OnSymbolTick.md) | `*pb.OnSymbolTickData` | Real-time bid/ask ticks. |
| [`OnTrade(ctx)`](./OnTrade.md) | `*pb.OnTradeData` | Order lifecycle events (new/updated/removed). |
| [`OnOpenedOrdersProfit(ctx, intervalMs)`](./OnOpenedOrdersProfit.md) | `*pb.OnOpenedOrdersProfitData` | Periodic floating P/L per open order. |
| [`OnOpenedOrdersTickets(ctx, intervalMs)`](./OnOpenedOrdersTickets.md) | `*pb.OnOpenedOrdersTicketsData` | Periodic snapshot of open tickets. |
| [`OnConnectState(ctx, terminalIds)`](./OnConnectState.md) | `*pb.ConnectStateData` | Terminal connection-state changes. |
| [`OrdersHistoryStream(ctx, sort, from, to, pageSize)`](./OrdersHistoryStream.md) | `*pb.OrdersHistoryData` | Paged history, delivered page by page. |
| [`QuoteHistoryStream(ctx, symbol, tf, from, to, chunk)`](./QuoteHistoryStream.md) | `*pb.QuoteHistoryData` | History bars, chunked by time window. |

---

## Lifecycle rules

* **Always `defer cancel()`** — leaking a context leaks the stream goroutine.
* **Drain the error channel** — a non-nil error means the stream is done; reads after that return zero values.
* **One stream, one goroutine** — don't share a channel across consumers.
* **Idle ≠ error** — many streams simply go quiet when nothing happens; see the examples' `common.ReadStream`, which treats `DeadlineExceeded` as "no events in window".

---

**Back to:** [MT4Account — Overview](../Overview.md)
