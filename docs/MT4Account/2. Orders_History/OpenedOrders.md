# OpenedOrders (RECOMMENDED)

> Every open order — market positions **and** resting pendings — in one call.

**API Information:**

* **Low-level API:** `MT4Account.OpenedOrders(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **RPC:** `OpenedOrders(OpenedOrdersRequest) → OpenedOrdersReply`
* **Payload:** `OpenedOrdersData` → `GetOrderInfos()` → `OpenedOrderInfo`

## 💬 Just the essentials

* **What it is.** A snapshot of all currently open orders with full detail (price, lots, SL/TP, profit, magic).
* **Why you need it.** The source of truth for "what do I have open right now?" — sizing, reconciliation, management.
* **When to use.** Before/after trading, in dashboards, to reconcile after a timed-out send.

---

## 🎯 Purpose

* List open **positions** (BUY/SELL) and **pendings** (LIMIT/STOP)
* Read each order's lots, open price, SL/TP, floating profit
* Filter by symbol or magic for "my" orders
* Reconcile after a failed/timed-out `OrderSend`

---

## Method Signature

```go
func (a *MT4Account) OpenedOrders(ctx context.Context) (*pb.OpenedOrdersData, error)
```

**Request message:**

```protobuf
message OpenedOrdersRequest {}   // no arguments
```

**Payload message:**

```protobuf
message OpenedOrdersData { repeated OpenedOrderInfo order_infos = 1; }

message OpenedOrderInfo {
  int32  ticket = ...;
  string symbol = ...;
  OpenedOrderType order_type = ...;
  double lots = ...;
  double open_price = ...;
  double stop_loss = ...;
  double take_profit = ...;
  double profit = ...;
  double swap = ...;
  double commision = ...;
  int32  magic_number = ...;
  string comment = ...;
  Timestamp open_time = ...;
  Timestamp expiration_time = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output — `*pb.OpenedOrdersData`

`GetOrderInfos()` returns `[]*pb.OpenedOrderInfo`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `ticket` | int32 | `GetTicket() int32` | Order ticket |
| `symbol` | string | `GetSymbol() string` | Instrument |
| `order_type` | OpenedOrderType | `GetOrderType()` | Buy/sell, market or pending |
| `lots` | double | `GetLots() float64` | Volume |
| `open_price` | double | `GetOpenPrice() float64` | Entry / pending price |
| `stop_loss` / `take_profit` | double | `GetStopLoss()` / `GetTakeProfit()` | Protective levels (0 = none) |
| `profit` | double | `GetProfit() float64` | Floating P/L (money) |
| `swap` / `commision` | double | `GetSwap()` / `GetCommision()` | Costs |
| `magic_number` | int32 | `GetMagicNumber() int32` | EA / strategy tag |
| `comment` | string | `GetComment() string` | Order comment |
| `open_time` / `expiration_time` | Timestamp | `GetOpenTime()` / `GetExpirationTime()` | Open / pending expiry |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧱 Related enums (from proto)

### `OpenedOrderType` — field `order_type`

| Constant | Value | Meaning |
|----------|-------|---------|
| `OO_OP_BUY` | 0 | Market buy (position) |
| `OO_OP_SELL` | 1 | Market sell (position) |
| `OO_OP_BUYLIMIT` | 2 | Buy limit (pending) |
| `OO_OP_SELLLIMIT` | 3 | Sell limit (pending) |
| `OO_OP_BUYSTOP` | 4 | Buy stop (pending) |
| `OO_OP_SELLSTOP` | 5 | Sell stop (pending) |
| `OO_OP_BALANCE` | 6 | Balance entry (non-trading) |
| `OO_OP_CREDIT` | 7 | Credit entry (non-trading) |

> `OO_OP_BUY`/`OO_OP_SELL` are open **positions**; `*LIMIT`/`*STOP` are **pendings**. Compare via `o.GetOrderType() == pb.OpenedOrderType_OO_OP_BUY`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Positions vs pendings.** Filter on `GetOrderType()` — the two market types are positions, the four `*LIMIT`/`*STOP` are pendings.
* **Floating profit** is `GetProfit()`; add `GetSwap()` + `GetCommision()` for the net.
* **Reconcile after errors.** If an `OrderSend` times out, call this to see whether it actually filled before retrying.
* **For just the tickets** (lighter), use [`OpenedOrdersTickets`](./OpenedOrdersTickets.md).

---

## 🔗 Usage Examples

### 1) Print all open orders

```go
func listOpen(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return err }
    for _, o := range d.GetOrderInfos() {
        fmt.Printf("#%d %s %s %.2f @ %.5f  P/L %+.2f\n",
            o.GetTicket(), o.GetSymbol(), o.GetOrderType(),
            o.GetLots(), o.GetOpenPrice(), o.GetProfit())
    }
    return nil
}
```

### 2) Split positions from pendings

```go
import pb "github.com/MetaRPC/GoMT4/package"

func split(ctx context.Context, acc *mt4.MT4Account) (positions, pendings []*pb.OpenedOrderInfo, err error) {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return nil, nil, err }
    for _, o := range d.GetOrderInfos() {
        switch o.GetOrderType() {
        case pb.OpenedOrderType_OO_OP_BUY, pb.OpenedOrderType_OO_OP_SELL:
            positions = append(positions, o)
        case pb.OpenedOrderType_OO_OP_BUYLIMIT, pb.OpenedOrderType_OO_OP_SELLLIMIT,
            pb.OpenedOrderType_OO_OP_BUYSTOP, pb.OpenedOrderType_OO_OP_SELLSTOP:
            pendings = append(pendings, o)
        }
    }
    return positions, pendings, nil
}
```

### 3) Total floating P/L for "my" magic

```go
func myFloatingPL(ctx context.Context, acc *mt4.MT4Account, magic int32) (float64, error) {
    d, err := acc.OpenedOrders(ctx)
    if err != nil { return 0, err }
    var pl float64
    for _, o := range d.GetOrderInfos() {
        if o.GetMagicNumber() == magic {
            pl += o.GetProfit() + o.GetSwap() + o.GetCommision()
        }
    }
    return pl, nil
}
```

---

## ⚠️ Pitfalls

* **Snapshot, not a stream.** Re-query for fresh values, or use [`OnOpenedOrdersProfit`](../5.%20Streaming_Methods/OnOpenedOrdersProfit.md).
* **Balance/credit rows.** `OO_OP_BALANCE` / `OO_OP_CREDIT` aren't trades — skip them when summing exposure.
* **Empty is normal.** A flat account returns zero `OrderInfos` — guard `len(...)`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open a 0.01 order; it appears here with the right `GetOrderType()` and `GetLots()`.
* **Round-trip.** After [`OrderClose`](../4.%20Trading_Operations/OrderClose.md) it no longer appears.
* **Magic filter.** Tag an order with a magic and assert the filtered sum includes only it.

---

## 📚 See also

* [OpenedOrdersTickets](./OpenedOrdersTickets.md) · [OrdersHistory](./OrdersHistory.md) · [OrderSelect](./OrderSelect.md) · [OrdersTotal](./OrdersTotal.md)
* **Back to:** [2. Orders & History — Overview](./Orders_History.Overview.md) · [MT4Account — Overview](../Overview.md)
