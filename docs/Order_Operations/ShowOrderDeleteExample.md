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

The ticket must reference a **pending** order.

---

## ⬆️ Output

Returns `*pb.OrderCloseDeleteData`:

| Field                 | Type     | Description                             |
| --------------------- | -------- | --------------------------------------- |
| `Mode`                | `string` | Operation result (e.g., "Deleted").     |
| `HistoryOrderComment` | `string` | Server comment about the deleted order. |

---

## 🎯 Purpose

Remove limit/stop orders that are no longer needed:

* Cancel pending orders before execution
* Manage pending order queues
* Clean up unused test/training orders

---

## 🧩 Notes & Tips

* **Pending only:** For open market positions use `OrderClose`, not `OrderDelete`.
* **Idempotency:** If the order is already filled/expired/cancelled, the API may return an error with a broker comment.
* **No price/slippage:** Delete uses only the ticket; price and slippage parameters are not applicable.
