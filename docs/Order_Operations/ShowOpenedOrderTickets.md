# Getting Opened Order Tickets

> **Request:** fetch only the tickets (IDs) of currently opened orders
> Retrieve a lightweight ticket list without full order details.

---

### Code Example

```go
// Using service wrapper
service.ShowOpenedOrderTickets(context.Background())

// Or directly from MT4Account
ticketsData, err := mt4.OpenedOrdersTickets(context.Background())
if err != nil {
    log.Fatalf("Error retrieving opened order tickets: %v", err)
}

for _, ticket := range ticketsData.GetTickets() {
    fmt.Printf("Open Order Ticket: %d\n", ticket)
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOpenedOrderTickets(ctx context.Context)
```

---

## 🔽 Input

| Field | Type              | Description                          |
| ----- | ----------------- | ------------------------------------ |
| `ctx` | `context.Context` | For timeout and cancellation control |

---

## ⬆️ Output

This method prints the ticket IDs of opened orders to stdout.
Underlying response: `*pb.OpenedOrdersTicketsData`

| Field     | Type       | Description                          |
| --------- | ---------- | ------------------------------------ |
| `Tickets` | `[]uint64` | List of ticket IDs for opened orders |

---

## 🎯 Purpose

Retrieve only **open order IDs** without full details. Useful for:

* Rapid synchronization or order tracking
* Lightweight matching against local state
* Selecting targets for modification/cancellation

---

## 🧩 Notes & Tips

* **Tickets only:** To inspect volumes, prices, or symbols, follow up with `OpenedOrders` or `HistoryOrderByTicket`.
* **Uniqueness:** Ticket IDs are unique per account; always treat them as `uint64`.
* **Performance:** Ideal for high-frequency polling or lightweight checks.

---

## ⚠️ Pitfalls

* **No orders open:** The API returns an empty slice, not `nil`. Always handle gracefully.
* **Stale state:** If orders are rapidly opened/closed, snapshot may be outdated in milliseconds. For real-time, use streams if available.
