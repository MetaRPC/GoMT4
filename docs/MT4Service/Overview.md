# MT4Service — Mid-Level API Overview

> The friendly facade over [`MT4Account`](../MT4Account/Overview.md). Same operations, but with **clean return values** (wrapper `…Data` unwrapped to slices/scalars), simplified arguments, and clearly-named streaming methods. What most code uses.

Defined in [`mt4/mt4service.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/mt4service.go). Wraps one `*MT4Account`; borrow it with `svc.Account()`.

```go
svc := mt4.NewMT4Service(acc)

orders, _ := svc.OpenedOrders(ctx)      // []*pb.OpenedOrderInfo (unwrapped, not a wrapper)
tickets, _ := svc.OpenedOrdersTickets(ctx) // []int32
syms, _ := svc.Symbols(ctx)             // []*pb.SymbolNameInfo
q, _ := svc.Quote(ctx, "EURUSD")        // *pb.QuoteData
```

---

## What it adds over MT4Account

* **Unwrapped returns** — `OpenedOrders` → `[]*OpenedOrderInfo`, `OpenedOrdersTickets` → `[]int32`, `Symbols` → `[]*SymbolNameInfo`.
* **Simpler arguments** — e.g. `TradeStats(ctx, from)` instead of `(ctx, from, fromUtc)`.
* **Clean stream names** — `StreamTicks`, `StreamTrades`, `StreamOpenedTickets`, `StreamOpenedProfit`, `StreamConnectState`.

---

## Section pages

Each group has its own page with signatures, arguments, returns and examples:

* **[1. Account Information](./1.%20Account_Information.md)** — summary + trade stats / equity / unrealized P/L.
* **[2. Orders & History](./2.%20Orders_History.md)** — open orders, tickets, closed history.
* **[3. Symbols & Quotes](./3.%20Symbols_Quotes.md)** — symbols, select, params, tick values, quotes, bars, timezone.
* **[4. Trading Operations](./4.%20Trading_Operations.md)** — send / modify / close / close-by.
* **[5. Streaming Methods](./5.%20Streaming_Methods.md)** — ticks, trades, ticket lists, per-order profit, connect state.

---

## Method index (25)

| Group | Methods |
|-------|---------|
| [Account Info](./1.%20Account_Information.md) | `AccountSummary`, `TradeStats`, `TradeStatsEquityHistory`, `EquityHistory`, `TradeUnrealizedPL` |
| [Orders & History](./2.%20Orders_History.md) | `OpenedOrders`, `OpenedOrdersTickets`, `OrdersHistory` |
| [Symbols & Quotes](./3.%20Symbols_Quotes.md) | `Symbols`, `SymbolSelect`, `SymbolParamsMany`, `TickValueWithSize`, `Quote`, `QuoteMany`, `QuoteHistory`, `ServerTimezone` |
| [Trading](./4.%20Trading_Operations.md) | `OrderSend`, `OrderModify`, `OrderClose`, `OrderCloseBy` |
| [Streaming](./5.%20Streaming_Methods.md) | `StreamTicks`, `StreamTrades`, `StreamOpenedTickets`, `StreamOpenedProfit`, `StreamConnectState` |

> Order construction by side, pip math and risk sizing live one layer up in [`MT4Sugar`](../MT4Sugar/Overview.md).

---

**See also:** [MT4Account](../MT4Account/Overview.md) · [MT4Sugar](../MT4Sugar/Overview.md)
