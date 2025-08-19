# Closing an Order

> **Request:** close or delete an active order by its ticket
> Sends a request to terminate the specified trade.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Closes an order by ticket and prints result.
svc.ShowOrderCloseExample(ctx, 123456)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.
// ⚠️ This actually closes a trade — use on demo or with caution.

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

res, err := account.OrderClose(ctx, 123456, nil, nil, nil)
if err != nil {
    log.Fatalf("❌ OrderClose error: %v", err)
}

fmt.Printf("Closed: %s | Comment: %s\n",
    res.GetMode(),
    res.GetHistoryOrderComment(),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderCloseExample(ctx context.Context, ticket int32)
```

---

## 🔽 Input

| Field    | Type              | Description                     |
| -------- | ----------------- | ------------------------------- |
| `ctx`    | `context.Context` | Timeout / cancellation control. |
| `ticket` | `int32`           | Ticket of the order to close.   |
| `price`  | `*float64`        | Optional close price.           |
| `slip`   | `*int32`          | Optional slippage (points).     |
| `magic`  | `*int32`          | Optional magic ID.              |

---

## ⬆️ Output

Result object:

| Field                 | Type     | Description                              |
| --------------------- | -------- | ---------------------------------------- |
| `Mode`                | `string` | Operation result (e.g., "Closed").       |
| `HistoryOrderComment` | `string` | Server comment on the closure operation. |

---

## 🎯 Purpose

Close or delete an order by ticket.
Useful for manual interventions, post-trade cleanup, or testing order workflows.
