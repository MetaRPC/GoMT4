# OnOpenedOrdersProfit

> Periodic floating P/L for every open order, pushed on a timer you set.

**API Information:**

* **Low-level API:** `MT4Account.OnOpenedOrdersProfit(ctx, intervalMs)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.Subscriptions`
* **RPC:** `OnOpenedOrdersProfit(OnOpenedOrdersProfitRequest) → stream OnOpenedOrdersProfitReply`
* **Payload:** `OnOpenedOrdersProfitData`

## 💬 Just the essentials

* **What it is.** A heartbeat of open-position profit, every `intervalMs` milliseconds.
* **Why you need it.** Live P/L widgets, equity-based risk rules.
* **When to use.** When you want a steady cadence rather than event-driven [`OnTrade`](./OnTrade.md).

---

## 🎯 Purpose

* Stream per-order floating P/L at a fixed interval
* Read a live equity/profit/margin snapshot each tick
* Drive P/L dashboards and drawdown guards

---

## Method Signature

```go
func (a *MT4Account) OnOpenedOrdersProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersProfitData, <-chan error)
```

**Request message:**

```protobuf
message OnOpenedOrdersProfitRequest { int32 timer_period_milliseconds = 1; }
```

**Payload message:**

```protobuf
message OnOpenedOrdersProfitData {
  OnEventAccountInfo account_info = 1;
  repeated OnOpenedOrdersProfitOrderInfo opened_orders_with_profit_updated = 2;
  MT4_SUB_ENUM_EVENT_GROUP_TYPE type = 3;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to unsubscribe and close the channels. |
| `intervalMs` | `int32` | yes | Push period in ms (e.g. `1000` = once per second). |

---

## ⬆️ Frame — `*pb.OnOpenedOrdersProfitData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `account_info` | OnEventAccountInfo | `GetAccountInfo()` | `GetEquity()`, `GetProfit()`, `GetMarginLevel()`, `GetBalance()` |
| `opened_orders_with_profit_updated` | repeated OnOpenedOrdersProfitOrderInfo | `GetOpenedOrdersWithProfitUpdated()` | Per-order P/L |

Each per-order info: `GetTicket()` `int32`, `GetSymbol()` `string`, `GetLots()` `float64`, `GetOpenPrice()` `float64`, `GetOrderProfit()` `float64`, `GetSwap()` / `GetCommission()` `float64`.

> A non-nil value on the **error channel** ends the stream; both channels close on `ctx` cancel.

---

## 🧩 Notes & Tips

* **Pick a sane interval** — 500–1000 ms is plenty; small values add load for no benefit.
* **Empty slice** — a frame with no open orders still arrives; `GetAccountInfo()` is still valid.
* **Getter name** is `GetOpenedOrdersWithProfitUpdated()`, not `GetOrders()`.

---

## 🔗 Usage Examples

### 1) Live equity heartbeat (1 Hz)

```go
func heartbeat(ctx context.Context, acc *mt4.MT4Account) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    data, errs := acc.OnOpenedOrdersProfit(ctx, 1000)
    for {
        select {
        case d := <-data:
            ai := d.GetAccountInfo()
            fmt.Printf("equity=%.2f floatP/L=%.2f (%d orders)\n",
                ai.GetEquity(), ai.GetProfit(), len(d.GetOpenedOrdersWithProfitUpdated()))
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 2) Per-order P/L, bail on drawdown

```go
func bailOnDrawdown(ctx context.Context, acc *mt4.MT4Account, limit float64) error {
    data, errs := acc.OnOpenedOrdersProfit(ctx, 1000)
    for {
        select {
        case d := <-data:
            for _, o := range d.GetOpenedOrdersWithProfitUpdated() {
                fmt.Printf("#%d %s %.2f lot  P/L %+.2f\n",
                    o.GetTicket(), o.GetSymbol(), o.GetLots(), o.GetOrderProfit())
            }
            if d.GetAccountInfo().GetProfit() < -limit {
                return fmt.Errorf("drawdown past %.2f — stopping", limit)
            }
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 3) Alert on margin level

```go
func marginAlert(ctx context.Context, acc *mt4.MT4Account, floor float64) error {
    data, errs := acc.OnOpenedOrdersProfit(ctx, 2000)
    for {
        select {
        case d := <-data:
            if ml := d.GetAccountInfo().GetMarginLevel(); ml > 0 && ml < floor {
                fmt.Printf("⚠️ margin level %.0f%%\n", ml)
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

* **Interval too small** wastes CPU/bandwidth for no gain.
* **Nil account info** on a keep-alive — guard before dereferencing.
* **Leaking the context** leaks the goroutine — `defer cancel()`.

---

## 🧪 Testing Suggestions

* **Happy path (demo with an open order).** Frames arrive at ~`intervalMs`; the order appears with a finite P/L.
* **Flat account.** Frames still arrive; the per-order slice is empty.
* **Drawdown guard.** Simulate a loss and assert `bailOnDrawdown` returns.

---

## 📚 See also

* [OnOpenedOrdersTickets](./OnOpenedOrdersTickets.md) · [OnTrade](./OnTrade.md) · [AccountSummary](../1.%20Account_Information/AccountSummary.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
