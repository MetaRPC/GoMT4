# Streaming — Overview

This section groups together methods for **real-time streaming** from MT4.
These allow you to subscribe to continuous updates: order profits, tickets, quotes, and trade events.

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

---

## ⚡ Example Workflow

```go
// Example: subscribe to profit updates and stop when threshold hit
profitCh, errCh := mt4.OnOpenedOrdersProfit(ctx, 1000)

for {
    select {
    case pkt := <-profitCh:
        for _, info := range pkt.OpenedOrdersWithProfitUpdated {
            if info.OrderProfit < -100 {
                fmt.Println("⚠️ Drawdown alert!", info.Ticket)
                cancel()
            }
        }
    case err := <-errCh:
        log.Println("Stream error:", err)
        return
    }
}
```

---

## ✅ Best Practices

1. Always manage **cancellation** via `context.Context`.
2. Use **polling intervals** (`intervalMs`) wisely to balance speed vs. load.
3. Combine different streams for full dashboards (e.g., trades + quotes).
4. Implement error handling and automatic reconnects in production.

---

## 🎯 Purpose

The streaming block is designed for:

* Real-time dashboards and monitoring
* Building automated alert systems
* Keeping UIs synchronized with terminal state
* Tracking both market data and account activity

---

👉 Use this overview as a **map**, and jump into each `.md` file for full method details.
