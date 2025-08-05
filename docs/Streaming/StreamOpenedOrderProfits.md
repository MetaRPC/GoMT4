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

fmt.Println("\uD83D\uDD04 Streaming order profits...")

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

> `intervalMs` is set internally in the wrapper call to 1000 ms by default.

---

## ⬆️ Output

Stream of `OnOpenedOrdersProfitOrderInfo` objects:

| Field          | Type                 | Description                                 |
| -------------- | -------------------- | ------------------------------------------- |
| `Ticket`       | `int32`              | Order ticket ID                             |
| `Symbol`       | `string`             | Trading symbol (e.g., "EURUSD")             |
| `Lots`         | `float64`            | Trade volume in lots                        |
| `Profit`       | `float64`            | Current floating profit/loss                |
| `OpenPrice`    | `float64`            | Price at which the order was opened         |
| `CurrentPrice` | `float64`            | Current market price                        |
| `OpenTime`     | `timestamp`          | Order open time                             |
| `OrderType`    | `ENUM_ORDER_TYPE_TF` | Type of trade: Buy, Sell, etc.              |
| `Magic`        | `int32`              | Magic number identifying source or strategy |
| `Comment`      | `string`             | Custom comment attached to the order        |

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

This method allows **real-time tracking of floating P/L per open order**.

Use cases include:

* Updating per-order PnL widgets in dashboards
* Monitoring exposure and live risk
* Alerting systems based on drawdown/profit triggers

Optimized for continuous updates using a polling interval in milliseconds.
