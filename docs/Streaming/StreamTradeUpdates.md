# 🌐 Streaming Trade Updates

> **Request:** subscribe to real-time trade update stream
> Starts a server-side stream to receive trade activity as it happens.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Streams trade events and prints the first new order of each event; the demo stops after 3s without an event.
svc.StreamTradeUpdates(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

tradeCh, errCh := account.OnTrade(ctx)

fmt.Println("📊 Streaming trade updates...")
for {
    select {
    case trade, ok := <-tradeCh:
        if !ok {
            fmt.Println("✅ Trade stream ended.")
            return
        }

        ev := trade.GetEventData()
        for _, order := range ev.GetNewOrders() {
            fmt.Printf("[New] Ticket: %d | Symbol: %s | Type: %v | Volume: %.2f | Profit: %.2f\n",
                order.GetTicket(),
                order.GetSymbol(),
                order.GetType(),
                order.GetLots(),
                order.GetOrderProfit())
        }
        for _, u := range ev.GetUpdatedOrders() {
            fmt.Printf("[Updated] Ticket: %d | SL: %.5f -> %.5f\n",
                u.GetCurrent().GetTicket(), u.GetPrevious().GetStopLoss(), u.GetCurrent().GetStopLoss())
        }
        for _, h := range ev.GetNewHistoryOrders() {
            fmt.Printf("[Closed/Deleted] Ticket: %d | Profit: %.2f\n", h.GetTicket(), h.GetOrderProfit())
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
func (s *MT4Service) StreamTradeUpdates(ctx context.Context)
func (a *MT4Account) OnTrade(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error)
```

---

## 🔽 Input

| Field | Type              | Description                           |
| ----- | ----------------- | ------------------------------------- |
| `ctx` | `context.Context` | Used to cancel or control the stream. |

---

## ⬆️ Output

Stream of `*pb.OnTradeData` messages:

| Field                    | Type                               | Description                                  |
| ------------------------ | ---------------------------------- | -------------------------------------------- |
| `EventData`              | `*pb.OnTadeEventData`              | What changed (see below).                    |
| `AccountInfo`            | `*pb.OnEventAccountInfo`           | `Balance`, `Credit`, `Equity`, `Margin`, `FreeMargin`, `Profit`, `MarginLevel`, `Login`. |
| `Type`                   | `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE` | Event group (`OrderUpdate`).                 |
| `TerminalInstanceGuidId` | `string`                           | Id of the terminal that sent it.             |

### Structure: `pb.OnTadeEventData`

(The type name is spelled `OnTadeEventData` in the generated code.)

| Field              | Type                            | Description                                  |
| ------------------ | ------------------------------- | -------------------------------------------- |
| `NewOrders`        | `[]*pb.OnTradeOrderInfo`        | Orders that appeared (opened / placed).      |
| `UpdatedOrders`    | `[]*pb.OnTradeUpdatedOrderInfo` | Changed orders: `Previous` and `Current` state. |
| `RemovedOrders`    | `[]*pb.OnTradeOrderInfo`        | Orders that left the open list.              |
| `NewHistoryOrders` | `[]*pb.OnTradeOrderInfo`        | Orders that appeared in history (closed / deleted). |

### Structure: `pb.OnTradeOrderInfo`

| Field         | Type                          | Description                           |
| ------------- | ----------------------------- | ------------------------------------- |
| `Ticket`      | `int32`                       | Unique identifier of the trade order  |
| `Symbol`      | `string`                      | Trading symbol (e.g., "EURUSD")       |
| `Type`        | `pb.SUB_ORDER_OPERATION_TYPE` | Trade type (Buy/Sell/Stop/Limit etc)  |
| `Lots`        | `float64`                     | Volume in lots                        |
| `OpenPrice`   | `float64`                     | Opening price of the trade            |
| `OrderProfit` | `float64`                     | Profit/loss                           |
| `StopLoss` / `TakeProfit` | `float64`         | SL / TP levels                        |
| `OpenTime` / `OpenTimeUtc` | `*timestamppb.Timestamp` | Open time (broker server time / UTC) |
| `ClosePrice`, `CloseTime` / `CloseTimeUtc` | `float64`, `*timestamppb.Timestamp` | Only for history orders |
| `IsHistory`   | `bool`                        | Whether this is a history order       |
| `Comment`     | `string`                      | Trade comment                         |
| `MagicNumber` | `int32`                       | Magic number for programmatic tagging |

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

This method lets you **monitor live trading activity** in real time — it’s the central source for:

* Updating user dashboards and UI widgets
* Triggering alerts and post-trade actions
* Building audit trails and analytics

> ⚠️ The stream is continuous. Make sure to implement cancellation or filtering logic as needed for production use.
