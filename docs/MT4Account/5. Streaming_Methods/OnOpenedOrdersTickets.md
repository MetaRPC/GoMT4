# OnOpenedOrdersTickets

> Periodic snapshot of which tickets are currently open — positions and pendings, as two slices.

**API Information:**

* **Low-level API:** `MT4Account.OnOpenedOrdersTickets(ctx, intervalMs)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.Subscriptions`
* **RPC:** `OnOpenedOrdersTickets(OnOpenedOrdersTicketsRequest) → stream OnOpenedOrdersTicketsReply`
* **Payload:** `OnOpenedOrdersTicketsData`

## 💬 Just the essentials

* **What it is.** A lightweight ticket list on a timer — just the IDs, split into positions and pendings.
* **Why you need it.** Cheap "what's open right now" polling to detect additions/removals.
* **When to use.** Reconciliation loops where full order detail would be overkill.

---

## 🎯 Purpose

* Stream the open-ticket set at a fixed interval
* Separate positions from pendings
* Diff between frames to detect fills/closes

---

## Method Signature

```go
func (a *MT4Account) OnOpenedOrdersTickets(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersTicketsData, <-chan error)
```

**Request message:**

```protobuf
message OnOpenedOrdersTicketsRequest { int32 pull_interval_milliseconds = 1; }
```

**Payload message:**

```protobuf
message OnOpenedOrdersTicketsData {
  repeated int32 position_tickets = 1;
  repeated int32 pending_order_tickets = 2;
  Timestamp server_time = 3;
  Timestamp server_time_utc = 5;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to unsubscribe and close the channels. |
| `intervalMs` | `int32` | yes | Pull period in ms. |

---

## ⬆️ Frame — `*pb.OnOpenedOrdersTicketsData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `position_tickets` | repeated int32 | `GetPositionTickets() []int32` | Open market positions |
| `pending_order_tickets` | repeated int32 | `GetPendingOrderTickets() []int32` | Resting pendings |
| `server_time` / `server_time_utc` | Timestamp | `GetServerTime()` / `GetServerTimeUtc()` | Snapshot time |

> A non-nil value on the **error channel** ends the stream; both channels close on `ctx` cancel.

---

## 🧩 Notes & Tips

* **Diff between frames** to spot new fills or closes without the heavier [`OnTrade`](./OnTrade.md) payload.
* **Two separate slices** — positions and pendings never mix.
* **IDs only** — for price/volume/SL you still need [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).

---

## 🔗 Usage Examples

### 1) Print the open ticket sets

```go
func watchTickets(ctx context.Context, acc *mt4.MT4Account) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    data, errs := acc.OnOpenedOrdersTickets(ctx, 1000)
    for {
        select {
        case d := <-data:
            fmt.Printf("positions=%v pendings=%v\n",
                d.GetPositionTickets(), d.GetPendingOrderTickets())
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 2) Detect a new fill

```go
func onNewPosition(ctx context.Context, acc *mt4.MT4Account) error {
    seen := map[int32]bool{}
    first := true
    data, errs := acc.OnOpenedOrdersTickets(ctx, 1000)
    for {
        select {
        case d := <-data:
            for _, t := range d.GetPositionTickets() {
                if !seen[t] && !first { fmt.Printf("new position #%d\n", t) }
                seen[t] = true
            }
            first = false
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 3) Stop when flat

```go
func untilFlat(ctx context.Context, acc *mt4.MT4Account) error {
    data, errs := acc.OnOpenedOrdersTickets(ctx, 1000)
    for {
        select {
        case d := <-data:
            if len(d.GetPositionTickets()) == 0 && len(d.GetPendingOrderTickets()) == 0 {
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

* **IDs only** — no detail; pair with `OpenedOrders` when you need more.
* **First frame** is a baseline, not a change — diff from the second frame on.
* **Leaking the context** leaks the goroutine — `defer cancel()`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open an order; its ticket appears in `GetPositionTickets()`.
* **Diff.** Close it and assert it drops out on the next frame.
* **Flat.** Both slices empty on a flat account.

---

## 📚 See also

* [OnOpenedOrdersProfit](./OnOpenedOrdersProfit.md) · [OpenedOrdersTickets](../2.%20Orders_History/OpenedOrdersTickets.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
