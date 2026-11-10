# OnTrade

> Order-lifecycle events — fires whenever an order is opened, changed, closed or moved to history.

**API Information:**

* **Low-level API:** `MT4Account.OnTrade(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.Subscriptions`
* **RPC:** `OnTrade(OnTradeRequest) → stream OnTradeReply`
* **Payload:** `OnTradeData`

## 💬 Just the essentials

* **What it is.** A push notification for every trade event on the account.
* **Why you need it.** React to fills/closes without polling `OpenedOrders`.
* **When to use.** Live order tracking, audit, copy-trading.

---

## 🎯 Purpose

* Receive new / updated / removed / history order deltas as they happen
* Read an equity/margin snapshot captured at each event
* Drive copy-trading or an audit log
* Avoid polling for order changes

---

## Method Signature

```go
func (a *MT4Account) OnTrade(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error)
```

**No arguments** beyond `ctx` — it reports **all** account events (no filter). Cancel the context to unsubscribe.

**Request message:**

```protobuf
message OnTradeRequest {}   // no arguments
```

**Payload message:**

```protobuf
message OnTradeData {
  MT4_SUB_ENUM_EVENT_GROUP_TYPE type = 1;
  OnEventAccountInfo account_info = 2;
  OnTadeEventData event_data = 3;
}
```

---

## ⬆️ Frame — `*pb.OnTradeData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `type` | MT4_SUB_ENUM_EVENT_GROUP_TYPE | `GetType()` | Event group |
| `account_info` | OnEventAccountInfo | `GetAccountInfo()` | Equity/margin snapshot (`GetEquity()`, `GetProfit()`, `GetMarginLevel()`, …) |
| `event_data` | OnTadeEventData | `GetEventData()` | The order deltas |

`GetEventData()` carries four slices: `GetNewOrders()` / `GetRemovedOrders()` / `GetNewHistoryOrders()` (`[]*OnTradeOrderInfo`) and `GetUpdatedOrders()` (`[]*OnTradeUpdatedOrderInfo`).

> A non-nil value on the **error channel** ends the stream; both channels close on `ctx` cancel.

---

## 🧱 Related enums (from proto)

### `MT4_SUB_ENUM_EVENT_GROUP_TYPE` — field `type`

| Constant | Value | Meaning |
|----------|-------|---------|
| `OrderProfit` | 0 | A profit-update group |
| `OrderUpdate` | 1 | An order add/change/remove group |

> In Go: `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Quiet by design** — no frames until something trades; idle is not an error.
* **`GetAccountInfo()`** gives equity/profit/margin level at the moment of the event — handy for risk.
* **Note the spelling** — the payload type is `OnTadeEventData` (as generated).
* **Guard sub-messages** — `GetEventData()` may be `nil` on a keep-alive.

---

## 🔗 Usage Examples

### 1) Count the deltas on each event

```go
func watchTrades(ctx context.Context, acc *mt4.MT4Account) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    data, errs := acc.OnTrade(ctx)
    for {
        select {
        case d := <-data:
            ev := d.GetEventData()
            if ev == nil { continue }
            fmt.Printf("new=%d updated=%d removed=%d\n",
                len(ev.GetNewOrders()), len(ev.GetUpdatedOrders()), len(ev.GetRemovedOrders()))
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 2) Log every newly opened order

```go
func logFills(ctx context.Context, acc *mt4.MT4Account) error {
    data, errs := acc.OnTrade(ctx)
    for {
        select {
        case d := <-data:
            ev := d.GetEventData()
            if ev == nil { continue }
            for _, o := range ev.GetNewOrders() {
                fmt.Printf("opened #%d %s  equity now %.2f\n",
                    o.GetTicket(), o.GetSymbol(), d.GetAccountInfo().GetEquity())
            }
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 3) Stop when an order closes

```go
func untilClose(ctx context.Context, acc *mt4.MT4Account) error {
    data, errs := acc.OnTrade(ctx)
    for {
        select {
        case d := <-data:
            if ev := d.GetEventData(); ev != nil && len(ev.GetRemovedOrders()) > 0 {
                fmt.Println("an order closed — done")
                return nil
            }
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

---

## ⚠️ Pitfalls

* **Nil sub-messages** on keep-alives — guard `GetEventData()` / `GetAccountInfo()`.
* **All events, no filter** — filter client-side by symbol/magic if you only care about some.
* **Leaking the context** leaks the goroutine — `defer cancel()`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open/close an order and assert a frame arrives with the right delta slice populated.
* **Idle.** With nothing trading, no frames arrive (not an error).
* **Equity.** `GetAccountInfo().GetEquity()` is a sensible number at each event.

---

## 📚 See also

* [OnOpenedOrdersProfit](./OnOpenedOrdersProfit.md) · [OpenedOrders](../2.%20Orders_History/OpenedOrders.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
