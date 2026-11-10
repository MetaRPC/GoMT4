# OrderDelete

> Delete a pending order by ticket.

**API Information:**

* **Low-level API:** `MT4Account.OrderDelete(ctx, ticket)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.TradingHelper`
* **RPC:** `OrderCloseDelete(OrderCloseDeleteRequest) → OrderCloseDeleteReply`
* **Payload:** `OrderCloseDeleteData`

## 💬 Just the essentials

* **What it is.** A convenience for cancelling a resting **pending** order (an `OrderClose` with no lots/price).
* **Why you need it.** Cancel a limit/stop you no longer want.
* **When to use.** Removing pendings. For market orders use [`OrderClose`](./OrderClose.md).

---

## 🎯 Purpose

* Cancel a resting pending order
* Confirm the deletion via the returned mode
* Clean up unfilled pendings (e.g. after a timeout)

---

## Method Signature

```go
func (a *MT4Account) OrderDelete(ctx context.Context, ticket int32) (*pb.OrderCloseDeleteData, error)
```

**Request message:**

```protobuf
message OrderCloseDeleteRequest {
  int32 order_ticket = 1;   // the pending to delete (no lots/price)
}
```

**Payload message:** `OrderCloseDeleteData` — same as [`OrderClose`](./OrderClose.md).

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `ticket` | `int32` | yes | Pending order to delete. |

---

## ⬆️ Output — `*pb.OrderCloseDeleteData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `mode` | OrderCloseDeleteMode | `GetMode()` | Should be `OCD_PENDING_ORDER` for a pending |
| `history_order_comment` | string | `GetHistoryOrderComment() string` | Broker's comment |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md). Not retried.

---

## 🧱 Related enums (from proto)

### `OrderCloseDeleteMode` — field `mode`

| Constant | Value | Meaning |
|----------|-------|---------|
| `OCD_MARKET_ORDER` | 0 | A market order was closed |
| `OCD_PENDING_ORDER` | 1 | A pending order was deleted |

> See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Pendings only.** Calling it on a filled market order is not its purpose — use [`OrderClose`](./OrderClose.md).
* **Same RPC as `OrderClose`.** `OrderDelete` just omits lots/price; the server treats the ticket as a pending.
* **Verify the mode** is `OCD_PENDING_ORDER` to be sure you cancelled a pending.

---

## 🔗 Usage Examples

### 1) Delete one pending

```go
import pb "github.com/MetaRPC/GoMT4/package"

func deletePending(ctx context.Context, acc *mt4.MT4Account, ticket int32) error {
    d, err := acc.OrderDelete(ctx, ticket)
    if err != nil { return err }
    if d.GetMode() == pb.OrderCloseDeleteMode_OCD_PENDING_ORDER {
        fmt.Println("pending deleted")
    }
    return nil
}
```

### 2) Cancel all pendings on a symbol

```go
func cancelPendings(ctx context.Context, acc *mt4.MT4Account, symbol string) (int, error) {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return 0, err }
    n := 0
    for _, o := range d.GetOrderInfos() {
        switch o.GetOrderType() {
        case pb.OpenedOrderType_OO_OP_BUYLIMIT, pb.OpenedOrderType_OO_OP_SELLLIMIT,
            pb.OpenedOrderType_OO_OP_BUYSTOP, pb.OpenedOrderType_OO_OP_SELLSTOP:
            if symbol == "" || o.GetSymbol() == symbol {
                if _, err := acc.OrderDelete(ctx, o.GetTicket()); err == nil { n++ }
            }
        }
    }
    return n, nil
}
```

### 3) Delete after a fill timeout

```go
func cancelIfUnfilled(ctx context.Context, acc *mt4.MT4Account, ticket int32) error {
    _, err := acc.OrderDelete(ctx, ticket)
    return err
}
```

---

## ⚠️ Pitfalls

* **Filled already?** If the pending filled before you delete it, expect an error — reconcile with [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).
* **Not retried.** Handle the error path.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Place a far pending, delete it, assert `OCD_PENDING_ORDER` and that it leaves `OpenedOrders`.
* **Already filled.** Deleting a filled ticket errors.
* **Bulk.** `cancelPendings` removes only pendings, leaving positions untouched.

---

## 📚 See also

* [OrderClose](./OrderClose.md) · [OrderSend](./OrderSend.md) · [OpenedOrders](../2.%20Orders_History/OpenedOrders.md)
* **Back to:** [4. Trading Operations — Overview](./Trading_Operations.Overview.md) · [MT4Account — Overview](../Overview.md)
