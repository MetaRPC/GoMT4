# ⚡ Streaming Opened Order Profits

> **Request:** subscribe to real-time stream of floating profit/loss per open order
> Starts a streaming channel to monitor floating profit for all active trades.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Streams profit updates for opened orders; the demo stops after 3s without an update.
svc.StreamOpenedOrderProfits(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Interval in milliseconds (server pulls/aggregates updates)
profitCh, errCh := account.OnOpenedOrdersProfit(ctx, 1000)

fmt.Println("🔄 Streaming order profits...")

for {
    select {
    case profit, ok := <-profitCh:
        if !ok {
            fmt.Println("✅ Profit stream ended.")
            return
        }
        for _, info := range profit.GetOpenedOrdersWithProfitUpdated() {
            fmt.Printf("[Profit] Ticket: %d | Symbol: %s | Profit: %.2f\n",
                info.GetTicket(), info.GetSymbol(), info.GetOrderProfit())
        }
        if ai := profit.GetAccountInfo(); ai != nil {
            fmt.Printf("[Account] Equity: %.2f | Margin: %.2f | FreeMargin: %.2f\n",
                ai.GetEquity(), ai.GetMargin(), ai.GetFreeMargin())
        }

    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return

    case <-time.After(30 * time.Second): // demo timeout
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}

```

---

### Method Signature

```go
func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context)
func (a *MT4Account) OnOpenedOrdersProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersProfitData, <-chan error)
```

---

## 🔽 Input

| Field        | Type              | Description                                       |
| ------------ | ----------------- | ------------------------------------------------- |
| `ctx`        | `context.Context` | For stream lifetime control and cancellation.     |
| `intervalMs` | `int32`           | Interval between updates in milliseconds (sent as `TimerPeriodMilliseconds`). |

> The wrapper uses **1000 ms**.

---

## ⬆️ Output

Stream of `*pb.OnOpenedOrdersProfitData`:

| Field                           | Type                                  | Description                                |
| ------------------------------- | ------------------------------------- | ------------------------------------------ |
| `OpenedOrdersWithProfitUpdated` | `[]*pb.OnOpenedOrdersProfitOrderInfo` | Per-order profit entries.                  |
| `AccountInfo`                   | `*pb.OnEventAccountInfo`              | `Balance`, `Credit`, `Equity`, `Margin`, `FreeMargin`, `Profit`, `MarginLevel`, `Login`. |
| `Type`                          | `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE`    | Event group (`OrderProfit`).               |
| `TerminalInstanceGuidId`        | `string`                              | Id of the terminal that sent it.           |

Each `*pb.OnOpenedOrdersProfitOrderInfo`:

| Field          | Type                          | Description                      |
| -------------- | ----------------------------- | -------------------------------- |
| `Ticket`       | `int32`                       | Order ticket ID.                 |
| `Symbol`       | `string`                      | Trading symbol (e.g., "EURUSD"). |
| `Lots`         | `float64`                     | Trade volume in lots.            |
| `OrderProfit`  | `float64`                     | Current floating profit/loss.    |
| `OpenPrice`    | `float64`                     | Order open price.                |
| `StopLoss` / `TakeProfit` | `float64`          | SL / TP levels.                  |
| `Swap` / `Commission` | `float64`              | Swap and commission.             |
| `OpenTime` / `OpenTimeUtc` | `*timestamppb.Timestamp` | Order open time (broker server time / UTC). |
| `Type`         | `pb.SUB_ORDER_OPERATION_TYPE` | Trade type.                      |
| `MagicNumber`  | `int32`                       | Strategy/magic number.           |
| `Comment`      | `string`                      | Order comment.                   |

---

### ENUM: `pb.SUB_ORDER_OPERATION_TYPE`

| Value                                  | Description        |
| -------------------------------------- | ------------------ |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_BUY`       | Buy order          |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_SELL`      | Sell order         |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_BUYLIMIT`  | Pending Buy Limit  |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_SELLLIMIT` | Pending Sell Limit |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_BUYSTOP`   | Pending Buy Stop   |
| `SUB_ORDER_OPERATION_TYPE_SUB_OP_SELLSTOP`  | Pending Sell Stop  |

---

## 🎯 Purpose

Enable **real-time tracking of floating P/L per open order** for dashboards, exposure monitoring, and alerting.

---

## 🧩 Notes & Tips

* **Interval trade-off:** Lower `intervalMs` → fresher updates, higher CPU/network use. 500–2000 ms works well for dashboards.
* **Keep state by ticket:** Maintain a map keyed by `Ticket` and update it from each packet.

---

## ⚠️ Pitfalls

* **Timeout vs cancel:** The example uses a demo timeout (`time.After`); prefer explicit `cancel()` for controlled shutdowns.
* **Ordering:** Do not assume chronological ordering across entries; sort by `Ticket` or timestamp if needed.
* **Float formatting:** Round only for display; keep raw values for calculations.

---

## 🧪 Testing Suggestions

* **Smoke test:** `intervalMs=1000` for a few minutes; verify steady updates and no leaks.
* **Burst test:** Open/close several orders quickly; ensure map/delta logic stays consistent.
* **Shutdown:** Cancel context and assert both channels close/return without blocking.
