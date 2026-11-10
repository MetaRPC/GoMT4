# MT4Sugar — High-Level API Overview

> One-liners and business logic over [`MT4Service`](../MT4Service/Overview.md): quick prices, symbol math, risk-based lot sizing, order entry by side, order management, trailing/breakeven, watchers/waiters and order info.

Defined in [`mt4/mt4sugar.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/mt4sugar.go). Composes an `*MT4Service`; borrow it with `sugar.Service()`.

```go
sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(777, 20, "my-bot")

lots, _ := sugar.CalcLotByRisk(ctx, "EURUSD", 50, 20)        // normalized lot risking 50 over 20 pips
r, _    := sugar.BuyMarket(ctx, "EURUSD", lots, nil, nil)    // *pb.OrderSendData
_ = sugar.ModifySLTPByPips(ctx, r.GetTicket(), ptr(20.0), ptr(40.0))
```

---

## Conventions

* **Defaults.** `NewMT4Sugar(svc).SetDefaults(magic, slippage, comment)` stamps the order-entry helpers.
* **Pips.** A pip is `point × 10` on 3/5-digit symbols, else `point`.
* **Orders, not positions.** `Orders` returns everything; `Positions` the open market ones.

---

## Navigation by group

| # | Group | Methods |
|---|-------|---------|
| 1 | [Connection & Defaults](./1.%20Connection_Defaults.md) | `ConnectByServerName` · `Disconnect` · `Ping` · `SetDefaults` |
| 2 | [Prices & Quotes](./2.%20Prices_Quotes.md) | `LastQuote` · `Quotes` · `MidPrice` · `SpreadPips` · `EnsureSymbol` · `Bars` |
| 3 | [Symbol Math](./3.%20Symbol_Math.md) | `Point` · `Digits` · `PipSize` · `PipsToPrice` · `PriceToPips` · `NormalizeLot` · `NormalizePrice` · `TickValue` |
| 4 | [Risk Sizing](./4.%20Risk_Sizing.md) | `CalcCashRisk` · `CalcLotByRisk` |
| 5 | [Order Entry](./5.%20Order_Entry.md) | `BuyMarket` · `SellMarket` · `BuyLimit` · `SellLimit` · `BuyStop` · `SellStop` |
| 6 | [Order Management](./6.%20Order_Management.md) | `Close` · `ClosePartial` · `CloseBy` · `CloseAll` · `CancelPendings` · `ModifySLTPByPrice` · `ModifySLTPByPips` |
| 7 | [Trailing & Breakeven](./7.%20Trailing_Breakeven.md) | `AutoBreakeven` · `TrailStop` |
| 8 | [Watchers & Waiters](./8.%20Watchers_Waiters.md) | `SubscribeQuotes` · `WatchProfit` · `WatchTickets` · `WaitPrice` · `WaitFilled` · `WaitClosed` |
| 9 | [Order Info](./9.%20Order_Info.md) | `Orders` · `Positions` · `OrdersForSymbol` · `OrdersByMagic` · `FindTicket` · `TotalProfit` · `PositionsValue` · `ExposureSummary` · `GetAccountInfo` |

---

**See also:** [MT4Service](../MT4Service/Overview.md) · [MT4Account](../MT4Account/Overview.md)
