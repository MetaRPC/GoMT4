# Closing an Order by Opposite Order

> **Request:** close one order using another opposite-position order
> Sends a request to close a position by matching it with an opposite order.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderCloseByExample(context.Background(), 123456, 654321)

// Or directly from MT4Account
result, err := mt4.OrderCloseBy(context.Background(), 123456, 654321)
if err != nil {
    log.Fatalf("Error closing by opposite: %v", err)
}

fmt.Printf("Closed by opposite: Profit=%.2f, Price=%.5f, Time: %s\n",
    result.GetProfit(),
    result.GetClosePrice(),
    result.GetCloseTime().AsTime().Format("2006-01-02 15:04:05"),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderCloseByExample(ctx context.Context, ticket int32, oppositeTicket int32)
```

---

## 🔽 Input

Required:

| Field            | Type              | Description                                |
| ---------------- | ----------------- | ------------------------------------------ |
| `ctx`            | `context.Context` | Context for timeout or cancellation.       |
| `ticket`         | `int32`           | The primary order ticket to be closed.     |
| `oppositeTicket` | `int32`           | The opposite-position order to close with. |

Both tickets must be valid and represent opposing open positions.

---

## ⬆️ Output

Returns the result of the closing operation:

| Field        | Type        | Description                         |
| ------------ | ----------- | ----------------------------------- |
| `Profit`     | `float64`   | Profit/loss realized from closing.  |
| `ClosePrice` | `float64`   | The closing price of the operation. |
| `CloseTime`  | `timestamp` | The time the orders were closed.    |

---

## 🎯 Purpose

Use this method to close one position with another opposite-position order.
This is useful for:

* Trade netting workflows
* Reducing exposure by pairing off positions
* Closing multiple positions efficiently

Ensure both tickets are valid and represent opposing trades for this operation to succeed.
