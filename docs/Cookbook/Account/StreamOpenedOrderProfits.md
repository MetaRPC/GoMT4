# 📈 StreamOpenedOrderProfits (GoMT4)

**Goal:** subscribe to **live profit updates** for currently opened orders with auto‑reconnect and clean shutdown.

> Real code refs you already ship:
>
> * Streaming helper: `examples/mt4/MT4Account.go` (`OnOpenedOrdersProfit`)
> * Service method: `examples/mt4/MT4_service.go` (`StreamOpenedOrderProfits`)
> * Docs: [Streaming → StreamOpenedOrderProfits](../../Streaming/StreamOpenedOrderProfits.md)

---

## ✅ 1) Preconditions

* `account` is connected.
* There are **open positions**; otherwise, the stream will be mostly silent.

---

## ▶️ 2) Service method example

This is `MT4Service.StreamOpenedOrderProfits` from `examples/mt4/MT4_service.go` (inside the `mt4` package, so it can use `s.account`):

```go
func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context) {

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    profitCh, errCh := s.account.OnOpenedOrdersProfit(ctx, 1000) // interval: 1000 ms

    fmt.Println("🔄 Streaming order profits...")

    for {
        select {
        case profit, ok := <-profitCh:
            if !ok {
                fmt.Println("✅ Profit stream ended.")
                return
            }

            // profit.OpenedOrdersWithProfitUpdated is []*OnOpenedOrdersProfitOrderInfo
            for _, info := range profit.OpenedOrdersWithProfitUpdated {
                fmt.Printf("[Profit] Ticket: %d | Symbol: %s | Profit: %.2f\n",
                    info.Ticket, info.Symbol, info.OrderProfit)
            }

        case err := <-errCh:
            log.Printf("❌ Stream error: %v", err)
            return

        case <-time.After(3 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}
```

In your own code, call `account.OnOpenedOrdersProfit(ctx, intervalMs)` the same way. Each packet also carries `AccountInfo` with `Balance`, `Equity`, `Margin`, `FreeMargin`, `MarginLevel`:

```go
if ai := profit.GetAccountInfo(); ai != nil && ai.GetMarginLevel() > 0 && ai.GetMarginLevel() < 100 {
    log.Printf("⚠️ Low margin level: %.1f%%", ai.GetMarginLevel())
}
```

---

## 🧭 3) Backoff & health checks

* Under the hood, `OnOpenedOrdersProfit` auto‑reconnects on transient errors (`io.EOF`, `codes.Unavailable`, terminal-not-found).
* Reconnect delays are controlled by constants in `MT4Account.go` (backoff+jitter).

---

## 🧰 4) Processing tips

* Keep the loop **non‑blocking** — heavy work (DB, reporting) should go to a worker via buffered channel.
* Aggregate by ticket to compute **running P/L** per position.
* Combine with `OnSymbolTick` for unrealized P/L analysis.

---

## ⚠️ Pitfalls

* **No positions** → stream may emit nothing; normal if no orders open.
* **Context canceled** → stream ends; always handle `<-ctx.Done()`.
* **Backpressure** → if you stop reading `profitCh`, the stream stalls.
* **Timeout** → in this demo, `time.After(3s)` ends the loop when no update arrives for 3 seconds; in production cancel the context instead.

---

## 🔗 See also

* [AccountSummary](AccountSummary.md) — snapshot of balance/equity.
* [StreamQuotes](../Market_Info/StreamQuotes.md) — combine with quotes to analyze P/L dynamics.
* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — timeouts, reconnects & backoff patterns.
