# OrderModify

> Change an order's price (pendings) and/or SL / TP / expiry.

**API Information:**

* **Low-level API:** `MT4Account.OrderModify(ctx, ticket, price, stoploss, takeprofit, expiration)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.TradingHelper`
* **RPC:** `OrderModify(OrderModifyRequest) → OrderModifyReply`
* **Payload:** `OrderModifyData` → `GetOrderWasModified()` (returned as a plain `bool`)

## 💬 Just the essentials

* **What it is.** MQL `OrderModify()` — adjust protective levels or a pending's trigger.
* **Why you need it.** Move SL to breakeven, trail, retarget a pending, set an expiry.
* **When to use.** Managing a live order after entry.

---

## 🎯 Purpose

* Set / move stop-loss and take-profit (positions and pendings)
* Change a pending order's trigger price
* Set or clear a pending's expiration
* Leave any field unchanged by passing `nil`

---

## Method Signature

```go
func (a *MT4Account) OrderModify(
    ctx context.Context,
    ticket int32,
    price, stoploss, takeprofit *float64,
    expiration *timestamppb.Timestamp,
) (bool, error)
```

**Request message:**

```protobuf
message OrderModifyRequest {
  int32 order_ticket = 1;
  optional double new_price = 2;        // pendings only
  optional double new_stop_loss = 3;
  optional double new_take_profit = 4;
  optional Timestamp new_expiration = 5;   // (+ new_expiration_utc)
}
```

**Payload message:**

```protobuf
message OrderModifyData { bool order_was_modified = 1; }
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `ticket` | `int32` | yes | Order to modify. |
| `price` | `*float64` | no | New trigger price (**pendings only**); `nil` leaves it. |
| `stoploss` | `*float64` | no | New SL price; `nil` leaves it. |
| `takeprofit` | `*float64` | no | New TP price; `nil` leaves it. |
| `expiration` | `*timestamppb.Timestamp` | no | New pending expiry; `nil` leaves it. |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| modified | `bool` | `true` if the server accepted the change (`OrderModifyData.order_was_modified`) |

> On failure the call returns `false` and a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md). Not retried.

---

## 🧩 Notes & Tips

* **`nil` = unchanged.** Pass `nil` for any field you don't want to touch — safe partial updates.
* **Prices, not pips.** Convert pip offsets to absolute prices first ([`MT4Sugar.ModifySLTPByPips`](../../MT4Sugar/6.%20Order_Management.md) does this).
* **Stops distance.** Respect `GetTradeStopsLevel()` from [`SymbolParams`](../3.%20Symbols_Quotes/SymbolParams.md) (130 otherwise).
* **`price` only matters for pendings** — on an open position it's ignored.

---

## 🔗 Usage Examples

### 1) Move SL to breakeven

```go
func toBreakeven(ctx context.Context, acc *mt4.MT4Account, ticket int32, entry float64) (bool, error) {
    sl := entry
    return acc.OrderModify(ctx, ticket, nil, &sl, nil, nil)
}
```

### 2) Set both SL and TP (absolute prices)

```go
func setStops(ctx context.Context, acc *mt4.MT4Account, ticket int32, sl, tp float64) (bool, error) {
    return acc.OrderModify(ctx, ticket, nil, &sl, &tp, nil)
}
```

### 3) Retarget a pending's trigger price

```go
func retarget(ctx context.Context, acc *mt4.MT4Account, ticket int32, newPrice float64) (bool, error) {
    return acc.OrderModify(ctx, ticket, &newPrice, nil, nil, nil)
}
```

---

## ⚠️ Pitfalls

* **"No changes" (error 1).** Re-sending identical levels is rejected — only modify when something differs.
* **`price` on a market position is ignored** — only SL/TP apply there.
* **Too-close stops** are rejected (130) — honor the stops level.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open an order, set an SL, assert `true` and that [`OrderSelect`](../2.%20Orders_History/OrderSelect.md) shows it.
* **No-op.** Re-send the same SL and assert the error path (no change).
* **Pending retarget.** Move a pending's price and confirm via `OrderSelect`.

---

## 📚 See also

* [OrderSend](./OrderSend.md) · [OrderClose](./OrderClose.md) · [MT4Sugar · Order Management](../../MT4Sugar/6.%20Order_Management.md)
* **Back to:** [4. Trading Operations — Overview](./Trading_Operations.Overview.md) · [MT4Account — Overview](../Overview.md)
