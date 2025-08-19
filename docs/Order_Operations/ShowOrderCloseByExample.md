# Closing an Order by Opposite Order

> **Request:** close one order using another opposite-position order
> Sends a request to close a position by matching it with an opposite order. 

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Closes one order by its opposite order.
svc.ShowOrderCloseByExample(ctx, 123456, 654321)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.
// ⚠️ This action closes trades — use on demo or with caution.

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := account.OrderCloseBy(ctx, 123456, 654321)
if err != nil {
    log.Fatalf("❌ OrderCloseBy error: %v", err)
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

Returns `*pb.OrderCloseByData` with fields:

| Field        | Type        | Description                         |
| ------------ | ----------- | ----------------------------------- |
| `Profit`     | `float64`   | Profit/loss realized from closing.  |
| `ClosePrice` | `float64`   | The closing price of the operation. |
| `CloseTime`  | `timestamp` | The time the orders were closed.    |

---

## 🎯 Purpose

Close one position with another opposite-position order. Useful for:

* Trade netting workflows
* Reducing exposure by pairing positions
* Closing multiple positions efficiently

---

## 🧩 Notes & Tips

* **Same symbol required:** Both tickets must be for the **same symbol** and opposite directions.
* **Partial overlap:** If lot sizes differ, only the overlapping volume is closed; the larger position remains with the residual volume.
* **Ticket types:** Tickets are `int32` in your API; avoid mixing with `uint64` types used elsewhere.

