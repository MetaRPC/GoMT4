# Deleting a Pending Order

> **Request:** delete a pending order by its ticket
> Sends a request to remove a pending order from the trading terminal.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderDeleteExample(context.Background(), 123456)

// Or directly from MT4Account
data, err := mt4.OrderDelete(context.Background(), 123456)
if err != nil {
    log.Fatalf("Error deleting order: %v", err)
}

fmt.Printf("Order deleted. Mode: %s, Comment: %s\n",
    data.GetMode(),
    data.GetHistoryOrderComment(),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderDeleteExample(ctx context.Context, ticket int32)
```

---

## 🔽 Input

Required:

| Field    | Type              | Description                         |
| -------- | ----------------- | ----------------------------------- |
| `ctx`    | `context.Context` | Timeout or cancellation management. |
| `ticket` | `int32`           | Ticket number of the pending order. |

The ticket must reference a **pending** order. Attempting to delete a non-pending order will result in an error.

---

## ⬆️ Output

Returns a result structure containing:

| Field                 | Type     | Description                             |
| --------------------- | -------- | --------------------------------------- |
| `Mode`                | `string` | Operation result (e.g., "Deleted").     |
| `HistoryOrderComment` | `string` | Server comment about the deleted order. |

---

## 🎯 Purpose

Use this method to manually delete pending orders that are no longer needed.
Helpful for:

* Canceling limit/stop orders before execution
* Managing pending order queues
* Cleaning up unused test/training orders

Ensure the ticket references a valid pending order to avoid errors.
