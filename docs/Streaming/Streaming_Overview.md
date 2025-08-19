# Streaming — Overview

This section groups together methods for **real-time streaming** from MT4.
Use them to subscribe to continuous updates: order profits, tickets, quotes, trade events — and now **paged order history** and **chunked quote history**.

---

## 📂 Methods in this Section

* [StreamOpenedOrderProfits.md](StreamOpenedOrderProfits.md)
  Subscribe to live updates of **floating profit/loss** for all active orders.

* [StreamOpenedOrderTickets.md](StreamOpenedOrderTickets.md)
  Stream the list of **open order ticket IDs** in real time.

* [StreamQuotes.md](StreamQuotes.md)
  Get a continuous stream of **tick data** (bid/ask) for selected symbols.

* [StreamTradeUpdates.md](StreamTradeUpdates.md)
  Receive live **trade activity events** as they happen.

* [StreamOrdersHistoryExample.md](StreamOrdersHistoryExample.md)
  Demo wrapper: last 30 days, sorted by close time (DESC), **page size = 200**.

* [StreamQuoteHistoryExample.md](StreamQuoteHistoryExample.md)
  Demo wrapper: last 90 days, **H1 timeframe**, **weekly chunks**.

---

## ⚡ Example Workflow

```go
// Example: subscribe to profit updates and stop when threshold hit
profitCh, errCh := mt4.OnOpenedOrdersProfit(ctx, 1000)

for {
    select {
    case pkt := <-profitCh:
        for _, info := range pkt.OpenedOrdersWithProfitUpdated {
            if info.OrderProfit < -100 { // your risk threshold
                fmt.Println("⚠️ Drawdown alert!", info.Ticket)
                cancel() // stop the stream
            }
        }
    case err := <-errCh:
        log.Println("Stream error:", err)
        return
    }
}
```

```go
// Example: stream historical orders in pages (service wrapper)
pagesCh, errCh := svc.StreamOrdersHistory(
    ctx,
    pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
    &from, &to,
    200, // page size
)

for {
    select {
    case page, ok := <-pagesCh:
        if !ok { return }
        for _, o := range page.GetOrdersInfo() {
            fmt.Printf("[HIST] %s | %d | %s | PnL: %.2f\n",
                o.GetOrderType(), o.GetTicket(), o.GetSymbol(), o.GetProfit())
        }
    case err := <-errCh:
        log.Println("history stream error:", err)
        return
    }
}
```

```go
// Example: stream OHLC history by time chunks (service wrapper)
barsCh, errCh := svc.StreamQuoteHistory(
    ctx,
    "EURUSD",
    pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1,
    from, to,
    7*24*time.Hour, // weekly chunks
)

for {
    select {
    case batch, ok := <-barsCh:
        if !ok { return }
        for _, c := range batch.GetHistoricalQuotes() {
            fmt.Printf("[%s] O: %.5f C: %.5f\n",
                c.GetTime().AsTime().Format("2006-01-02 15:04:05"),
                c.GetOpen(), c.GetClose(),
            )
        }
    case err := <-errCh:
        log.Println("quote history stream error:", err)
        return
    }
}
```

---

## ✅ Best Practices

1. Always manage **cancellation** via `context.Context` to avoid goroutine leaks.
2. Pick sensible **polling intervals** (e.g., `intervalMs`) to balance freshness vs. load.
3. Combine multiple streams for dashboards (e.g., trades + quotes + profits).
4. Always read from the **error channel**; log and exit/cancel on terminal failures.
5. For **orders history**, tune **page size** to your throughput (larger pages ⇒ fewer RPCs).
6. For **quote history**, adjust **chunk duration** to timeframe and volume (bigger chunks ⇒ fewer round-trips).
7. Keep print/formatting light in the stream loop; offload heavy processing to workers.

---

## 🎯 Purpose

The streaming block is designed for:

* Real-time dashboards and monitoring
* Automated alerting systems
* Synchronizing UIs with terminal/account state
* Backfilling analytics from historical orders and OHLC

---

👉 Use this overview as a **map**, then jump into each `.md` file for full method details and code snippets.
