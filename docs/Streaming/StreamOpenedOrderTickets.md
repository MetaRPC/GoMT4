# Streaming Opened Order Tickets

> **Request:** subscribe to stream of currently open order ticket numbers
> Returns only the **ticket IDs** of all open orders as they change in real time.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Streams all opened order tickets (positions + pendings); stops after ~30s in demo.
svc.StreamOpenedOrderTickets(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Interval in milliseconds (server polls updates)
ticketCh, errCh := account.OnOpenedOrdersTickets(ctx, 1000)

fmt.Println("🔄 Streaming opened order tickets...")

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

    case <-time.After(30 * time.Second): // demo timeout
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

> Wrapper default: **1000 ms**.

---

## ⬆️ Output

Stream of `*pb.OnOpenedOrdersTicketsData` packets. Each packet contains:

| Field                 | Type      | Description                   |
| --------------------- | --------- | ----------------------------- |
| `PositionTickets`     | `[]int32` | Ticket IDs of open positions. |
| `PendingOrderTickets` | `[]int32` | Ticket IDs of pending orders. |

Combined, they represent all open orders at the time of the poll.

---

## 🎯 Purpose

Track **open order ticket numbers** in real time for:

* Updating active trade lists in UIs
* Detecting order creation/closure events
* Triggering targeted follow-ups (fetch details, modify/close)

---

## 🧩 Notes & Tips

* **Diffing:** Maintain a `map[int32]bool` of previous tickets. On each packet, compute added/removed sets to detect events.
* **Batch ops:** When many changes occur at once, process `PositionTickets` and `PendingOrderTickets` separately if logic differs.
* **Minimal overhead:** Use this stream when you only need IDs; fetch details lazily on demand.
---

## 🧪 Testing Suggestions

* **Basic:** Verify non-empty tickets when known orders are open.
* **Add/Remove:** Open/close orders and confirm the diffing logic detects changes.
* **Timeout/Cancel:** Ensure context cancel cleanly stops both channels.
