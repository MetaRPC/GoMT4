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

Required:

ctx (context.Context) — context for managing timeout or cancellation. Can include a deadline or a cancellation token.

Optional context modifications:

Use context.WithDeadline to specify a deadline.

Use context.WithCancel to manage cancellation explicitly.
---

## ⬆️ Output

This method prints the ticket IDs of opened orders to stdout and does not directly return data.

Underlying returned structure:

| Field     | Type      | Description                           |
| --------- | --------- | ------------------------------------- |
| `Tickets` | `[]int32` | List of ticket IDs for opened orders. |

---

## 🎯 Purpose

Use this method when you only need ticket IDs of open orders, useful for:

* Rapid synchronization or matching processes
* Tracking order IDs without loading complete order details
* Quick selection for targeted operations such as modifications or cancellations

Provides an efficient, low-overhead alternative to retrieving full order details.
