# Getting Opened Order Tickets

> **Request:** fetch only the tickets (IDs) of currently opened orders
> Retrieve a lightweight ticket list without full order details.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints all currently opened order tickets.
svc.ShowOpenedOrderTickets(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

ticketsData, err := account.OpenedOrdersTickets(ctx)
if err != nil {
    log.Fatalf("❌ OpenedOrdersTickets error: %v", err)
}

for _, ticket := range ticketsData.GetTickets() {
    fmt.Printf("Open Order Ticket: %d\n", ticket)
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOpenedOrderTickets(ctx context.Context)
func (a *MT4Account) OpenedOrdersTickets(ctx context.Context) (*pb.OpenedOrdersTicketsData, error)
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
| `Tickets` | `[]int32`  | List of ticket IDs for opened orders |

---

## 🎯 Purpose

Retrieve only **open order IDs** without full details. Useful for:

* Rapid synchronization or order tracking
* Lightweight matching against local state
* Selecting targets for modification/cancellation

---

## 🧩 Notes & Tips

* **Tickets only:** To inspect volumes, prices, or symbols, follow up with `OpenedOrders` or `OrderSelect(ctx, ticket)`.
* **Uniqueness:** Ticket IDs are unique per account; in this API they are `int32`.
* **Performance:** Ideal for high-frequency polling or lightweight checks.

---

## ⚠️ Pitfalls

* **No orders open:** The API returns an empty slice, not `nil`. Always handle gracefully.
* **Stale state:** If orders are rapidly opened/closed, snapshot may be outdated in milliseconds. For real-time, use the `OnOpenedOrdersTickets` stream ([StreamOpenedOrderTickets](../Streaming/StreamOpenedOrderTickets.md)).
