# Streaming Opened Order Profits

> **Request:** subscribe to real-time stream of floating profit/loss per open order
> Starts a streaming channel to monitor floating profit for all active trades.

---

### Code Example

```go
// Using service wrapper
service.StreamOpenedOrderProfits(context.Background())

// Or directly from MT4Account
profitCh, errCh := mt4.OnOpenedOrdersProfit(context.Background(), 1000)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

fmt.Println("🔄 Streaming order profits...")

for {
    select {
    case profit, ok := <-profitCh:
        if !ok {
            fmt.Println("✅ Profit stream ended.")
            return
        }

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
```

---

### Method Signature

```go
func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context)
```

---

## 🔽 Input

| Field        | Type              | Description                                       |
| ------------ | ----------------- | ------------------------------------------------- |
| `ctx`        | `context.Context` | For stream lifetime control and cancellation.     |
| `intervalMs` | `int`             | Polling interval between updates in milliseconds. |

> In the wrapper, `intervalMs` defaults to **1000 ms**.

---

## ⬆️ Output

Stream of per-order profit updates.
Underlying response: `*pb.OnOpenedOrdersProfitData`

Each update contains entries like `OnOpenedOrdersProfitOrderInfo`:

| Field          | Type                 | Description                      |
| -------------- | -------------------- | -------------------------------- |
| `Ticket`       | `int32`              | Order ticket ID.                 |
| `Symbol`       | `string`             | Trading symbol (e.g., "EURUSD"). |
| `Lots`         | `float64`            | Trade volume in lots.            |
| `Profit`       | `float64`            | Current floating profit/loss.    |
| `OpenPrice`    | `float64`            | Order open price.                |
| `CurrentPrice` | `float64`            | Current market price.            |
| `OpenTime`     | `timestamp`          | Order open time.                 |
| `OrderType`    | `ENUM_ORDER_TYPE_TF` | Trade type.                      |
| `Magic`        | `int32`              | Strategy/magic number.           |
| `Comment`      | `string`             | Order comment.                   |

---

### ENUM: `ENUM_ORDER_TYPE_TF`

| Value                  | Description        |
| ---------------------- | ------------------ |
| `OrderTypeTfBuy`       | Buy order          |
| `OrderTypeTfSell`      | Sell order         |
| `OrderTypeTfBuyLimit`  | Pending Buy Limit  |
| `OrderTypeTfSellLimit` | Pending Sell Limit |
| `OrderTypeTfBuyStop`   | Pending Buy Stop   |
| `OrderTypeTfSellStop`  | Pending Sell Stop  |

---

## 🎯 Purpose

Enable **real-time tracking of floating P/L per open order** for dashboards, exposure monitoring, and alerting.

---

## 🧩 Notes & Tips

* **Interval trade-off:** Lower `intervalMs` → fresher updates, higher CPU/network use. 500–2000 ms works well for dashboards.
* **Delta batches:** Updates often include only orders that changed; maintain a map keyed by `Ticket` and apply deltas.
* **Backpressure:** Always read **both** data and error channels to avoid goroutine leaks.
* **UI smoothing:** Debounce rendering if bursts arrive at interval boundaries.

---

## ⚠️ Pitfalls

* **Timeout vs cancel:** Your example uses a 30s timeout; prefer explicit `cancel()` for controlled shutdowns.
* **Ordering:** Do not assume chronological ordering across entries; sort by `Ticket` or timestamp if needed.
* **Float formatting:** Round only for display; keep raw values for calculations.

---

## 🧪 Testing Suggestions

* **Smoke test:** `intervalMs=1000` for a few minutes; verify steady updates and no leaks.
* **Burst test:** Open/close several orders quickly; ensure map/delta logic stays consistent.
* **Shutdown:** Cancel context and assert both channels close/return without blocking.
