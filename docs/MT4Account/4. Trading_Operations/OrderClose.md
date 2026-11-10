# OrderClose

> Close a market order (full or partial) **or** delete a pending — the server picks based on the ticket.

**API Information:**

* **Low-level API:** `MT4Account.OrderClose(ctx, ticket, lots, price, slippage)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.TradingHelper`
* **RPC:** `OrderCloseDelete(OrderCloseDeleteRequest) → OrderCloseDeleteReply`
* **Payload:** `OrderCloseDeleteData`

## 💬 Just the essentials

* **What it is.** MQL `OrderClose()` for market orders and `OrderDelete()` for pendings, merged into one call.
* **Why you need it.** Exit a position (fully or partially) or cancel a pending.
* **When to use.** Closing out; scaling out with a partial close.

---

## 🎯 Purpose

* Fully close a market position
* Partially close a position by volume
* Delete a pending order
* Read back whether it closed a market order or deleted a pending

---

## Method Signature

```go
func (a *MT4Account) OrderClose(
    ctx context.Context,
    ticket int32,
    lots, price *float64,
    slippage *int32,
) (*pb.OrderCloseDeleteData, error)
```

**Request message:**

```protobuf
message OrderCloseDeleteRequest {
  int32 order_ticket = 1;
  optional double lots = 2;           // nil = whole order; partial (market only)
  optional double closing_price = 3;  // nil = market price
  optional int32  slippage = 4;       // points
}
```

**Payload message:**

```protobuf
message OrderCloseDeleteData {
  OrderCloseDeleteMode mode = 1;
  string history_order_comment = 2;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `ticket` | `int32` | yes | Order to close/delete. |
| `lots` | `*float64` | no | `nil` = close the whole order; `&x` = partial (market only). |
| `price` | `*float64` | no | `nil` = close at market. |
| `slippage` | `*int32` | no | Max deviation in points. |

---

## ⬆️ Output — `*pb.OrderCloseDeleteData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `mode` | OrderCloseDeleteMode | `GetMode()` | `OCD_MARKET_ORDER` (closed) / `OCD_PENDING_ORDER` (deleted) |
| `history_order_comment` | string | `GetHistoryOrderComment() string` | Broker's close comment |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md). Not retried.

---

## 🧱 Related enums (from proto)

### `OrderCloseDeleteMode` — field `mode`

| Constant | Value | Meaning |
|----------|-------|---------|
| `OCD_MARKET_ORDER` | 0 | A market order was closed |
| `OCD_PENDING_ORDER` | 1 | A pending order was deleted |

> Compare via `d.GetMode() == pb.OrderCloseDeleteMode_OCD_PENDING_ORDER`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **One call, two behaviours.** `GetMode()` tells you whether it closed a market order or deleted a pending.
* **Partial remainder.** A partial close can leave the leftover under a new ticket (MT4 rule) — re-read [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).
* **`lots` on a pending is ignored** — pendings are deleted whole (see [`OrderDelete`](./OrderDelete.md)).

---

## 🔗 Usage Examples

### 1) Full close

```go
func closeFull(ctx context.Context, acc *mt4.MT4Account, ticket int32) error {
    d, err := acc.OrderClose(ctx, ticket, nil, nil, nil)
    if err != nil { return err }
    fmt.Println("mode:", d.GetMode())
    return nil
}
```

### 2) Partial close (scale out 0.05 lot)

```go
func scaleOut(ctx context.Context, acc *mt4.MT4Account, ticket int32) error {
    lots := 0.05
    _, err := acc.OrderClose(ctx, ticket, &lots, nil, nil)
    return err
}
```

### 3) Flatten every open position

```go
import pb "github.com/MetaRPC/GoMT4/package"

func flattenAll(ctx context.Context, acc *mt4.MT4Account) (int, error) {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return 0, err }
    n := 0
    for _, o := range d.GetOrderInfos() {
        if o.GetOrderType() == pb.OpenedOrderType_OO_OP_BUY || o.GetOrderType() == pb.OpenedOrderType_OO_OP_SELL {
            if _, err := acc.OrderClose(ctx, o.GetTicket(), nil, nil, nil); err == nil { n++ }
        }
    }
    return n, nil
}
```

---

## ⚠️ Pitfalls

* **`lots` on a pending is ignored** — use [`OrderDelete`](./OrderDelete.md) for clarity.
* **Not retried.** Reconcile with [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md) on failure.
* **Partial leaves a new ticket** — don't assume the same ticket persists.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open then full-close; `GetMode()` is `OCD_MARKET_ORDER` and the ticket leaves `OpenedOrders`.
* **Partial.** Close part of a position; the remainder (possibly new ticket) still shows the reduced volume.
* **Pending.** Closing a pending reports `OCD_PENDING_ORDER`.

---

## 📚 See also

* [OrderDelete](./OrderDelete.md) · [OrderCloseBy](./OrderCloseBy.md) · [OrderSend](./OrderSend.md) · [OpenedOrders](../2.%20Orders_History/OpenedOrders.md)
* **Back to:** [4. Trading Operations — Overview](./Trading_Operations.Overview.md) · [MT4Account — Overview](../Overview.md)
