# 📡 API Reference — Streaming (GoMT4)

This page documents the **streaming messages** used in GoMT4. The four subscriptions are long-lived gRPC server streams (`SubscriptionService`) that continuously send updates; `MT4Account` wraps each in a method that returns a data channel and an error channel and reconnects on transient errors.

```go
dataCh, errCh := account.OnSymbolTick(ctx, []string{"EURUSD"})
```

* The data channel is unbuffered; read it promptly.
* On a final error the client sends it to `errCh` and closes both channels. Cancelling `ctx` ends the stream the same way.
* Every reply is `{ Data | Error }`; `Error` with `TERMINAL_INSTANCE_NOT_FOUND` / `TERMINAL_REGISTRY_TERMINAL_NOT_FOUND` makes the client reopen the stream, any other `Error` ends it.

---

## 🔔 Quotes Streaming — `account.OnSymbolTick(ctx, symbols)`

Request `OnSymbolTickRequest`: `SymbolNames []string`.

### OnSymbolTickData

| Field                    | Type                      |
| ------------------------ | ------------------------- |
| `SymbolTick`             | `*pb.OnSymbolMqlTickInfo` |
| `TerminalInstanceGuidId` | `string`                  |

### OnSymbolMqlTickInfo

| Field                      | Type                     |
| -------------------------- | ------------------------ |
| `Symbol`                   | `string`                 |
| `Bid`                      | `float64`                |
| `Ask`                      | `float64`                |
| `Last`                     | `float64`                |
| `Volume`                   | `uint64`                 |
| `Time` / `TimeUtc`         | `*timestamppb.Timestamp` |
| `TimeMsc` / `TimeMscUtc`   | `int64` (milliseconds)   |

**Notes:** one tick. Compute spread: `(Ask - Bid) / Point` (`Point` from `SymbolParams`).

---

## 🔄 Trade Updates — `account.OnTrade(ctx)`

Request `OnTradeRequest`: no fields.

### OnTradeData

| Field                    | Type                               |
| ------------------------ | ---------------------------------- |
| `Type`                   | `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE` |
| `EventData`              | `*pb.OnTadeEventData`              |
| `AccountInfo`            | `*pb.OnEventAccountInfo`           |
| `TerminalInstanceGuidId` | `string`                           |

### OnTadeEventData

| Field              | Type                            |
| ------------------ | ------------------------------- |
| `NewOrders`        | `[]*pb.OnTradeOrderInfo`        |
| `UpdatedOrders`    | `[]*pb.OnTradeUpdatedOrderInfo` (`Previous`, `Current`) |
| `RemovedOrders`    | `[]*pb.OnTradeOrderInfo`        |
| `NewHistoryOrders` | `[]*pb.OnTradeOrderInfo`        |

### OnTradeOrderInfo

| Field                        | Type                          |
| ---------------------------- | ----------------------------- |
| `Ticket`                     | `int32`                       |
| `Symbol`                     | `string`                      |
| `Type`                       | `pb.SUB_ORDER_OPERATION_TYPE` |
| `Lots`                       | `float64`                     |
| `OpenPrice`                  | `float64`                     |
| `StopLoss` / `TakeProfit`    | `float64`                     |
| `OrderProfit`                | `float64`                     |
| `Swap` / `Commission`        | `float64`                     |
| `OpenTime` / `OpenTimeUtc`   | `*timestamppb.Timestamp`      |
| `Expiration` / `ExpirationUtc` | `*timestamppb.Timestamp`    |
| `ClosePrice`                 | `float64` (history orders only) |
| `CloseTime` / `CloseTimeUtc` | `*timestamppb.Timestamp` (history orders only) |
| `IsHistory`                  | `bool`                        |
| `MagicNumber`                | `int32`                       |
| `Comment`                    | `string`                      |
| `Index`                      | `int32`                       |
| `AccountLogin`               | `int64`                       |

### OnEventAccountInfo

`Balance`, `Credit`, `Equity`, `Margin`, `FreeMargin`, `Profit`, `MarginLevel` (`float64`), `Login` (`int64`).

---

## 🎟️ Opened Orders Tickets — `account.OnOpenedOrdersTickets(ctx, intervalMs)`

Request `OnOpenedOrdersTicketsRequest`: `PullIntervalMilliseconds int32` (= `intervalMs`).

### OnOpenedOrdersTicketsData

| Field                          | Type                     |
| ------------------------------ | ------------------------ |
| `PositionTickets`              | `[]int32`                |
| `PendingOrderTickets`          | `[]int32`                |
| `ServerTime` / `ServerTimeUtc` | `*timestamppb.Timestamp` |
| `TerminalInstanceGuidId`       | `string`                 |

---

## 💰 Opened Orders Profit Stream — `account.OnOpenedOrdersProfit(ctx, intervalMs)`

Request `OnOpenedOrdersProfitRequest`: `TimerPeriodMilliseconds int32` (= `intervalMs`).

### OnOpenedOrdersProfitData

| Field                           | Type                                  |
| ------------------------------- | ------------------------------------- |
| `Type`                          | `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE`    |
| `OpenedOrdersWithProfitUpdated` | `[]*pb.OnOpenedOrdersProfitOrderInfo` |
| `AccountInfo`                   | `*pb.OnEventAccountInfo`              |
| `TerminalInstanceGuidId`        | `string`                              |

### OnOpenedOrdersProfitOrderInfo

| Field                        | Type                          |
| ---------------------------- | ----------------------------- |
| `Ticket`                     | `int32`                       |
| `Symbol`                     | `string`                      |
| `Type`                       | `pb.SUB_ORDER_OPERATION_TYPE` |
| `Lots`                       | `float64`                     |
| `OpenPrice`                  | `float64`                     |
| `OrderProfit`                | `float64`                     |
| `StopLoss` / `TakeProfit`    | `float64`                     |
| `Swap` / `Commission`        | `float64`                     |
| `OpenTime` / `OpenTimeUtc`   | `*timestamppb.Timestamp`      |
| `Expiration` / `ExpirationUtc` | `*timestamppb.Timestamp`    |
| `MagicNumber`                | `int32`                       |
| `Comment`                    | `string`                      |
| `Index`                      | `int32`                       |

**Notes:** real-time PnL per order plus account margin figures.

---

## 🕒 History helpers (client-side, not gRPC streams)

### `account.OrdersHistoryStream(ctx, sortType, from, to, pageSize)` → `<-chan *pb.OrdersHistoryData`

Calls `OrdersHistory` with page 1, 2, 3, … (`ItemsPerPage = pageSize`) and sends each non-empty page; stops at the first empty page. Same as `svc.StreamOrdersHistory(...)`.

### `account.QuoteHistoryStream(ctx, symbol, timeframe, from, to, chunk)` → `<-chan *pb.QuoteHistoryData`

Splits `[from, to]` into windows of `chunk` (default 24h), calls `QuoteHistory` per window and sends each non-empty result. Same as `svc.StreamQuoteHistory(...)`.

---

📌 **Tip:**

* Use `context.WithCancel` (not short timeouts) to control stream lifetime.
* Always read the error channel and check `ok` when receiving from the data channel.
* Streaming is ideal for real-time dashboards (quotes, PnL, trades).
