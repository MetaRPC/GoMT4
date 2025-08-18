# Modifying an Order

> **Request:** modify parameters (SL/TP) of an existing open order by its ticket
> Sends a request to update the stop loss and take profit values of a specified order.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderModifyExample(context.Background(), 123456)

// Or directly from MT4Account
newSL := 1.0500
newTP := 1.0900

modified, err := mt4.OrderModify(context.Background(), 123456, nil, &newSL, &newTP, nil)
if err != nil {
    log.Fatalf("Error modifying order: %v", err)
}

if modified {
    fmt.Println("Order successfully modified.")
} else {
    fmt.Println("Order was NOT modified.")
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderModifyExample(ctx context.Context, ticket int32)
```

---

## 🔽 Input

Required:

* `ctx` (`context.Context`) — context for managing timeout or cancellation.
* `ticket` (`int32`) — The ticket number of the order to modify.

Optional parameters (when directly using MT4Account):

| Field        | Type                     | Description                          |
| ------------ | ------------------------ | ------------------------------------ |
| `price`      | `*float64`               | New price (used mainly for pending). |
| `stopLoss`   | `*float64`               | New stop loss level.                 |
| `takeProfit` | `*float64`               | New take profit level.               |
| `expiration` | `*timestamppb.Timestamp` | New expiration (pending orders).     |

At least one of `price`, `stopLoss`, `takeProfit`, `expiration` must be provided.

---

## ⬆️ Output

The method prints whether modification succeeded.
Underlying response: `*pb.OrderModifyReply` → `Data.OrderWasModified` (`bool`).

---

## 🎯 Purpose

Adjust SL/TP (and, for pendings, price/expiration) on active orders—useful for dynamic risk management and strategy updates.

---

## 🧩 Notes & Tips

* **Digits & rounding:** Use `Digits` from `SymbolParams` to format/round SL/TP; avoid hardcoded decimals.
* **Volume/price rules:** Brokers enforce minimal distances/steps for SL/TP and pending prices; validate before sending to reduce rejects.
* **No-op guard:** If all optional params are `nil`, nothing will change; treat as a no-op.

---

## ⚠️ Pitfalls

* **Closed/non-existent ticket:** Will fail with a server-side error (e.g., ticket not found).
* **Invalid levels:** SL above price on sells (or below on buys), TP in the wrong direction, or levels inside the minimal distance will be rejected.
