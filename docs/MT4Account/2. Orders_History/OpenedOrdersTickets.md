# OpenedOrdersTickets

> Just the ticket numbers of all open orders — the lightweight "what's open?" read.

**API Information:**

* **Low-level API:** `MT4Account.OpenedOrdersTickets(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **Payload:** `OpenedOrdersTicketsData`

## 💬 Just the essentials

* **What it is.** The list of open-order ticket IDs — no prices, no volumes.
* **Why you need it.** Cheap polling to detect new fills / closes without the full payload.
* **When to use.** Reconciliation loops where full detail would be overkill.

---

## 🎯 Purpose

* Fetch all open-order tickets in one small call
* Diff between calls to spot additions / removals
* Drive follow-up detail reads via [`OrderSelect`](./OrderSelect.md)

---

## Method Signature

```go
func (a *MT4Account) OpenedOrdersTickets(ctx context.Context) (*pb.OpenedOrdersTicketsData, error)
```

**Request message:**

```protobuf
message OpenedOrdersTicketsRequest {}   // no arguments
```

**Payload message:**

```protobuf
message OpenedOrdersTicketsData { repeated int32 tickets = 1; }
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output — `*pb.OpenedOrdersTicketsData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `tickets` | repeated int32 | `GetTickets() []int32` | All open-order ticket IDs |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **IDs only.** For symbol/price/SL use [`OpenedOrders`](./OpenedOrders.md) or look one up with [`OrderSelect`](./OrderSelect.md).
* **No position/pending split** here — this unary call returns a flat ticket list. The streaming [`OnOpenedOrdersTickets`](../5.%20Streaming_Methods/OnOpenedOrdersTickets.md) splits them.
* **Diff two reads** to detect what opened or closed between polls.

---

## 🔗 Usage Examples

### 1) Print the open tickets

```go
func printTickets(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.OpenedOrdersTickets(ctx)
    if err != nil { return err }
    fmt.Printf("%d open: %v\n", len(d.GetTickets()), d.GetTickets())
    return nil
}
```

### 2) Detect closes between two polls

```go
func closedSince(ctx context.Context, acc *mt4.MT4Account, before []int32) ([]int32, error) {
    d, err := acc.OpenedOrdersTickets(ctx)
    if err != nil { return nil, err }
    now := map[int32]bool{}
    for _, t := range d.GetTickets() { now[t] = true }
    var gone []int32
    for _, t := range before {
        if !now[t] { gone = append(gone, t) }
    }
    return gone, nil
}
```

### 3) Fetch details for each ticket

```go
func detailAll(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.OpenedOrdersTickets(ctx)
    if err != nil { return err }
    for _, t := range d.GetTickets() {
        if o, err := acc.OrderSelect(ctx, t); err == nil {
            fmt.Printf("#%d %s %.2f\n", o.GetTicket(), o.GetSymbol(), o.GetLots())
        }
    }
    return nil
}
```

---

## ⚠️ Pitfalls

* **Snapshot.** Tickets can change the instant after the read — treat as point-in-time.
* **Empty is normal.** A flat account returns an empty slice — guard `len(...)`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Open an order; its ticket appears here.
* **Diff.** Close it and assert `closedSince` reports its ticket.
* **Consistency.** `len(GetTickets())` equals `len(OpenedOrders().GetOrderInfos())`.

---

## 📚 See also

* [OpenedOrders](./OpenedOrders.md) · [OrderSelect](./OrderSelect.md) · [OnOpenedOrdersTickets](../5.%20Streaming_Methods/OnOpenedOrdersTickets.md)
* **Back to:** [2. Orders & History — Overview](./Orders_History.Overview.md) · [MT4Account — Overview](../Overview.md)
