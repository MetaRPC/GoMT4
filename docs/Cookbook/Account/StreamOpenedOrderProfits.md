# 📈 StreamOpenedOrderProfits (GoMT4)

**Goal:** subscribe to **live profit updates** for currently opened orders with auto‑reconnect and clean shutdown.

> Real code refs you already ship:
>
> * Streaming helper: `examples/mt4/MT4Account.go` (`OnOpenedOrdersProfit`)
> * Service method: `examples/mt4/MT4_service.go` (`StreamOpenedOrderProfits`)
> * Docs: `docs/Streaming/StreamOpenedOrderProfits.md`

---

## ✅ 1) Preconditions

* MT4 terminal is running & connected.
* There are **open positions**; otherwise, the stream will be mostly silent.

---

## ▶️ 2) Service method example

```go
func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Second argument = buffer size for channel
    profitCh, errCh := s.account.OnOpenedOrdersProfit(ctx, 1000)

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

        case <-time.After(30 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}
```

---

## 🧭 3) Backoff & health checks

* Under the hood, `OnOpenedOrdersProfit` auto‑reconnects on transient errors (`io.EOF`, `codes.Unavailable`).
* Reconnect delays are controlled by constants in `MT4Account.go` (backoff+jitter).

---

## 🧰 4) Processing tips

* Keep the loop **non‑blocking** — heavy work (DB, reporting) should go to a worker via buffered channel.
* Aggregate by ticket to compute **running P/L** per position.
* Combine with `StreamQuotes` for unrealized P/L analysis.

---

## ⚠️ Pitfalls

* **No positions** → stream may emit nothing; normal if no orders open.
* **Context canceled** → stream ends; always handle `<-ctx.Done()>`.
* **Backpressure** → if you stop reading `profitCh`, the stream can stall.
* **Timeout** → in this demo, `time.After(30s)` closes the loop; in production use context deadlines.

---

## 🔗 See also

* `AccountSummary.md` — snapshot of balance/equity/margins.
* `StreamQuotes.md` — combine with quotes to analyze P/L dynamics.
* `Reliability (en)` — timeouts, reconnects & backoff patterns.
