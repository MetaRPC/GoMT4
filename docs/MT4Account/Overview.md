# MT4Account — Low-Level API Overview

> The low-level gRPC layer. Each method builds the protobuf `…Request`, attaches auth metadata, and unwraps the reply to `(*pb.…Data, error)`. It owns the channel, connection lifecycle, retry and streaming. Everything else builds on it.

Defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go) (package `mt4`).

```text
Your code → MT4Sugar → MT4Service → MT4Account → gRPC → MetaRPC gateway → MT4 terminal
```

Create and connect:

```go
acc, err := mt4.NewMT4Account(login, password, "", uuid.Nil, apiKey)
if err != nil { log.Fatal(err) }
defer acc.Disconnect(true)

ctx := context.Background()
if err := acc.ConnectByServerName(ctx, "MetaQuotes-Demo", "EURUSD", true, 30); err != nil {
    log.Fatal(err)
}
```

---

## 🔌 Connection (on this layer)
`NewMT4Account` · `NewMT4AccountWithApiKey` · `GetId` · `ConnectByHostPort` · `ConnectByServerName` · `Close` · `Disconnect` — build the account, connect the terminal, tear down. See the connection example above.

---

## Section pages

Each group has its own overview that drills into per-method pages (signature, request/reply, output getters, examples, pitfalls):

* **[1. Account Information](./1.%20Account_Information/Account_Information.Overview.md)** — summary, balance/equity getters, trade stats.
* **[2. Orders & History](./2.%20Orders_History/Orders_History.Overview.md)** — `OpenedOrders`, `OpenedOrdersTickets`, `OrdersHistory`, `OrderSelect`, `OrdersTotal`.
* **[3. Symbols & Quotes](./3.%20Symbols_Quotes/Symbols_Quotes.Overview.md)** — `Quote`, `QuoteMany`, `QuoteHistory`, `Symbols`, `SymbolSelect`, `SymbolParams`, `SymbolParamsMany`, `TickValueWithSize`, `ServerTimezone`.
* **[4. Trading Operations](./4.%20Trading_Operations/Trading_Operations.Overview.md)** — `OrderSend`, `OrderModify`, `OrderClose`, `OrderDelete`, `OrderCloseBy`.
* **[5. Streaming Methods](./5.%20Streaming_Methods/Streaming_Methods.Overview.md)** — `OnSymbolTick`, `OnTrade`, `OnOpenedOrdersProfit`, `OnOpenedOrdersTickets`, `OnConnectState`, `OrdersHistoryStream`, `QuoteHistoryStream`.

---

## Conventions

* Every method takes a `context.Context` (nil is accepted; a short per-call timeout is added for read-only calls).
* Reads are retried on transient transport errors and terminal-not-found (one reconnect); trade calls are **not** retried.
* Streams return a `(<-chan *pb.…Data, <-chan error)` pair — read frames off the data channel; drop/cancel the context to stop.

---

**See also:** [MT4Service](../MT4Service/Overview.md) · [MT4Sugar](../MT4Sugar/Overview.md)
