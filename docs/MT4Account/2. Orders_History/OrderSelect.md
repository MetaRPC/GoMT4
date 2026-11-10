# OrderSelect

> Look up one open order by ticket — the full detail for a single order.

**API Information:**

* **Low-level API:** `MT4Account.OrderSelect(ctx, ticket)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **Built on:** [`OpenedOrders`](./OpenedOrders.md) — it is a client-side lookup, **not** a separate RPC. It fetches the open orders and returns the one whose ticket matches.
* **Payload:** `OpenedOrderInfo` (one row of `OpenedOrders`)

## 💬 Just the essentials

* **What it is.** The single-order version of [`OpenedOrders`](./OpenedOrders.md) — one ticket, full detail.
* **Why you need it.** Read one order's price/SL/TP/profit without scanning the whole list.
* **When to use.** After you already have a ticket (from a send, a watcher, or `OpenedOrdersTickets`).

---

## 🎯 Purpose

* Fetch one open order by its ticket
* Read its side, lots, open price, SL/TP, floating profit
* Confirm an order still exists / inspect it before modifying

---

## Method Signature

```go
func (a *MT4Account) OrderSelect(ctx context.Context, ticket int32) (*pb.OpenedOrderInfo, error)
```

**How it works:** it calls [`OpenedOrders`](./OpenedOrders.md) and scans for the matching ticket:

```go
// (paraphrased from mt4/MT4Account.go)
for _, order := range opened.GetOrderInfos() {
    if order.GetTicket() == ticket {
        return order, nil
    }
}
return nil, fmt.Errorf("order with ticket %d not found", ticket)
```

**Payload:** `OpenedOrderInfo` — same shape as one row of [`OpenedOrders`](./OpenedOrders.md).

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `ticket` | `int32` | yes | Ticket of the order to read. |

---

## ⬆️ Output — `*pb.OpenedOrderInfo`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `ticket` | int32 | `GetTicket() int32` | Order ticket |
| `symbol` | string | `GetSymbol() string` | Instrument |
| `order_type` | OpenedOrderType | `GetOrderType()` | Buy/sell, market or pending (see [`OpenedOrders`](./OpenedOrders.md#related-enums-from-proto)) |
| `lots` | double | `GetLots() float64` | Volume |
| `open_price` | double | `GetOpenPrice() float64` | Entry / pending price |
| `stop_loss` / `take_profit` | double | `GetStopLoss()` / `GetTakeProfit()` | Protective levels |
| `profit` | double | `GetProfit() float64` | Floating P/L |
| `swap` / `commision` | double | `GetSwap()` / `GetCommision()` | Costs |
| `magic_number` | int32 | `GetMagicNumber() int32` | EA / strategy tag |

> An **unknown / closed ticket** returns `nil` and the error `"order with ticket N not found"`. A transport/API failure of the underlying `OpenedOrders` surfaces as its text error (`"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Open orders only.** For closed trades use [`OrdersHistory`](./OrdersHistory.md).
* **Check the type** with `GetOrderType()` before acting — a pending behaves differently from a position.
* **Pairs with watchers.** Use it to refresh one order's state inside a management loop.

---

## 🔗 Usage Examples

### 1) Read one order

```go
func showOrder(ctx context.Context, acc *mt4.MT4Account, ticket int32) error {
    o, err := acc.OrderSelect(ctx, ticket)
    if err != nil { return err }
    fmt.Printf("#%d %s %s %.2f @ %.5f  SL=%.5f TP=%.5f  P/L %+.2f\n",
        o.GetTicket(), o.GetSymbol(), o.GetOrderType(), o.GetLots(),
        o.GetOpenPrice(), o.GetStopLoss(), o.GetTakeProfit(), o.GetProfit())
    return nil
}
```

### 2) Guard: is it a position or a pending?

```go
import pb "github.com/MetaRPC/GoMT4/package"

func isPosition(ctx context.Context, acc *mt4.MT4Account, ticket int32) (bool, error) {
    o, err := acc.OrderSelect(ctx, ticket)
    if err != nil { return false, err }
    t := o.GetOrderType()
    return t == pb.OpenedOrderType_OO_OP_BUY || t == pb.OpenedOrderType_OO_OP_SELL, nil
}
```

### 3) Does the order still exist?

```go
func stillOpen(ctx context.Context, acc *mt4.MT4Account, ticket int32) bool {
    _, err := acc.OrderSelect(ctx, ticket)
    return err == nil
}
```

---

## ⚠️ Pitfalls

* **Unknown / closed ticket** returns an error, not a nil order — branch on `err`.
* **Snapshot.** Values are point-in-time; re-read inside a loop for fresh profit.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open an order, select it by ticket, assert the fields match the send.
* **Closed ticket.** After closing, `OrderSelect` on that ticket returns an error.
* **Type.** A pending returns a `*LIMIT`/`*STOP` `GetOrderType()`.

---

## 📚 See also

* [OpenedOrders](./OpenedOrders.md) · [OpenedOrdersTickets](./OpenedOrdersTickets.md) · [OrderModify](../4.%20Trading_Operations/OrderModify.md)
* **Back to:** [2. Orders & History — Overview](./Orders_History.Overview.md) · [MT4Account — Overview](../Overview.md)
