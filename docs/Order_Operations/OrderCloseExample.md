# Closing an Order

> **Request:** close or delete an active order by its ticket
> Sends a request to terminate the specified trade.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderCloseExample(context.Background(), 123456)

// Or directly from MT4Account
res, err := mt4.OrderClose(context.Background(), 123456, nil, nil, nil)
if err != nil {
    log.Fatalf("Error closing order: %v", err)
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
