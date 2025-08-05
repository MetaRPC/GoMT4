# Streaming Opened Order Tickets

> **Request:** subscribe to stream of currently open order ticket numbers
> Returns only the **ticket IDs** of all open orders as they change in real time.

---

### Code Example

```go
// Using service wrapper
service.StreamOpenedOrderTickets(context.Background())

// Or directly from MT4Account
ticketCh, errCh := mt4.OnOpenedOrdersTickets(context.Background(), 1000)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

fmt.Println("\uD83D\uDD04 Streaming opened order tickets...")

for {
    select {
    case pkt, ok := <-ticketCh:
        if !ok {
            fmt.Println("✅ Ticket stream ended.")
            return
        }
        tix := append(pkt.PositionTickets, pkt.PendingOrderTickets...)
        fmt.Printf("[Tickets] %d open tickets: %v\n", len(tix), tix)

    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return

    case <-time.After(30 * time.Second):
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}
```

---

### Method Signature

```go
func (s *MT4Service) StreamOpenedOrderTickets(ctx context.Context)
```

---

## 🔽 Input

| Field        | Type              | Description                                       |
| ------------ | ----------------- | ------------------------------------------------- |
| `ctx`        | `context.Context` | Controls streaming lifecycle and cancellation.    |
| `intervalMs` | `int`             | Polling interval between updates in milliseconds. |

> Note: `intervalMs` is hardcoded to 1000ms in the wrapper implementation, but can be made configurable.

---

## ⬆️ Output

Stream of `OpenedOrderTicketsData` objects, each including:

| Field                 | Type      | Description                            |
| --------------------- | --------- | -------------------------------------- |
| `PositionTickets`     | `[]int32` | List of position order ticket numbers. |
| `PendingOrderTickets` | `[]int32` | List of pending order ticket numbers.  |

Together, they represent all open orders in the terminal at each polling cycle.

---

## 🎯 Purpose

Use this method to **track open order ticket numbers** in real time.
Useful for:

* Updating active trade lists in UIs
* Detecting order creation/deletion events
* Triggering related updates or monitoring logic

It’s a **minimal-overhead** solution compared to full order detail streaming.
