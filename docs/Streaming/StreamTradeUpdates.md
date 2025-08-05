# Streaming Trade Updates

> **Request:** subscribe to real-time trade update stream
> Starts a server-side stream to receive trade activity as it happens.

---

### Code Example

```go
// Using service wrapper
service.StreamTradeUpdates(context.Background())

// Or directly from MT4Account
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

stream, err := mt4.OnTrade(ctx)
if err != nil {
    log.Fatalf("Stream error: %v", err)
}

fmt.Println("\uD83D\uDCCA Streaming trade updates...")
for trade := range stream {
    fmt.Printf("Trade ticket: %d, Symbol: %s, Profit: %.2f\n",
        trade.TradeInfo.Ticket,
        trade.TradeInfo.Symbol,
        trade.TradeInfo.Profit)
    break // for demo
}
```

---

### Method Signature

```go
func (s *MT4Service) StreamTradeUpdates(ctx context.Context)
```

---

## 🔽 Input

| Field | Type              | Description                           |
| ----- | ----------------- | ------------------------------------- |
| `ctx` | `context.Context` | Used to cancel or control the stream. |

---

## ⬆️ Output

Stream of `OnTradeData` messages. Each message contains a `TradeInfo` structure:

### Structure: `TradeInfo`

| Field       | Type     | Description                           |
| ----------- | -------- | ------------------------------------- |
| `Ticket`    | `int`    | Unique identifier of the trade order  |
| `Symbol`    | `string` | Trading symbol (e.g., "EURUSD")       |
| `Lots`      | `float`  | Volume in lots                        |
| `OpenPrice` | `float`  | Opening price of the trade            |
| `Profit`    | `float`  | Current floating P/L                  |
| `OpenTime`  | `string` | UTC timestamp of order open time      |
| `OrderType` | `int32`  | Trade type (Buy/Sell/Stop/Limit etc)  |
| `Comment`   | `string` | Trade comment                         |
| `Magic`     | `int`    | Magic number for programmatic tagging |

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

This method lets you **monitor live trading activity** in real time — it’s the central source for:

* Updating user dashboards and UI widgets
* Triggering alerts and post-trade actions
* Building audit trails and analytics

> ⚠️ The stream is continuous. Make sure to implement cancellation or filtering logic as needed for production use.
