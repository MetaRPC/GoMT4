# OrderCloseBy

> Net two opposite orders on the same symbol against each other — close one *by* the other.

**API Information:**

* **Low-level API:** `MT4Account.OrderCloseBy(ctx, ticketToClose, oppositeTicket)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.TradingHelper`
* **RPC:** `OrderCloseBy(OrderCloseByRequest) → OrderCloseByReply`
* **Payload:** `OrderCloseByData`

## 💬 Just the essentials

* **What it is.** MQL `OrderCloseBy()` — closes a buy against a sell on the same symbol, saving one spread.
* **Why you need it.** Flatten a hedged pair in a single operation.
* **When to use.** You hold opposite positions on one symbol and want to net them.

---

## 🎯 Purpose

* Close one position *by* an opposite one on the same symbol
* Save a spread vs closing both separately
* Read back the netting price and realized profit

---

## Method Signature

```go
func (a *MT4Account) OrderCloseBy(ctx context.Context, ticketToClose, oppositeTicket int32) (*pb.OrderCloseByData, error)
```

**Request message:**

```protobuf
message OrderCloseByRequest {
  int32 ticket_to_close = 1;
  int32 opposite_ticket_closing_by = 2;
}
```

**Payload message:**

```protobuf
message OrderCloseByData {
  double close_price = 1;
  double profit = 2;
  Timestamp close_time = 3;
  Timestamp close_time_utc = 4;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `ticketToClose` | `int32` | yes | The order being closed. |
| `oppositeTicket` | `int32` | yes | An **opposite-side** order on the **same symbol**. |

---

## ⬆️ Output — `*pb.OrderCloseByData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `close_price` | double | `GetClosePrice() float64` | Price at which the netting happened |
| `profit` | double | `GetProfit() float64` | Realized P/L |
| `close_time` / `close_time_utc` | Timestamp | `GetCloseTime()` / `GetCloseTimeUtc()` | Close time (server / UTC) |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md). Not retried.

---

## 🧩 Notes & Tips

* **Same symbol, opposite side.** A buy closes a sell and vice versa; different symbols are rejected.
* **Unequal volumes** leave a remainder open (possibly a new ticket) — re-read [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).
* **Broker support varies** — some brokers disable close-by; handle the error.

---

## 🔗 Usage Examples

### 1) Net a known hedged pair

```go
func net(ctx context.Context, acc *mt4.MT4Account, buyTicket, sellTicket int32) error {
    d, err := acc.OrderCloseBy(ctx, buyTicket, sellTicket)
    if err != nil { return err }
    fmt.Printf("closed by @ %.5f  profit %+.2f\n", d.GetClosePrice(), d.GetProfit())
    return nil
}
```

### 2) Find and net an opposite pair on a symbol

```go
import pb "github.com/MetaRPC/GoMT4/package"

func netSymbol(ctx context.Context, acc *mt4.MT4Account, symbol string) error {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return err }
    var buy, sell int32
    for _, o := range d.GetOrderInfos() {
        if o.GetSymbol() != symbol { continue }
        switch o.GetOrderType() {
        case pb.OpenedOrderType_OO_OP_BUY:  buy = o.GetTicket()
        case pb.OpenedOrderType_OO_OP_SELL: sell = o.GetTicket()
        }
    }
    if buy == 0 || sell == 0 {
        return fmt.Errorf("no opposite pair on %s", symbol)
    }
    _, err = acc.OrderCloseBy(ctx, buy, sell)
    return err
}
```

### 3) Net and report the leftover

```go
func netAndCheck(ctx context.Context, acc *mt4.MT4Account, a, b int32) error {
    if _, err := acc.OrderCloseBy(ctx, a, b); err != nil { return err }
    open, err := acc.OpenedOrders(ctx)
    if err != nil { return err }
    fmt.Printf("remaining open: %d\n", len(open.GetOrderInfos()))
    return nil
}
```

---

## ⚠️ Pitfalls

* **Different symbols** are rejected — both tickets must be the same instrument.
* **Partial net.** Unequal lots leave a remainder; don't assume full flatten.
* **Not supported everywhere** — branch on the error if the broker disables close-by.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open equal buy+sell on one symbol, net them, assert both leave `OpenedOrders`.
* **Unequal.** Net 0.10 against 0.05 and confirm a 0.05 remainder remains.
* **Cross-symbol.** Netting tickets on different symbols errors.

---

## 📚 See also

* [OrderClose](./OrderClose.md) · [OrderSend](./OrderSend.md) · [OpenedOrders](../2.%20Orders_History/OpenedOrders.md)
* **Back to:** [4. Trading Operations — Overview](./Trading_Operations.Overview.md) · [MT4Account — Overview](../Overview.md)
