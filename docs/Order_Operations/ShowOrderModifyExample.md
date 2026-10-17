# Modifying an Order

> **Request:** modify parameters (SL/TP) of an existing open order by its ticket
> Sends a request to update the stop loss and take profit values of a specified order.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Modifies SL/TP for an order by ticket.
// ⚠️ This changes a live trade — use on demo or with caution.
svc.ShowOrderModifyExample(ctx, 123456)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

newSL := 1.0500
newTP := 1.0900

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

modified, err := account.OrderModify(ctx, 123456, nil, &newSL, &newTP, nil)
if err != nil {
    log.Fatalf("❌ OrderModify error: %v", err)
}

if modified {
    fmt.Println("✅ Order successfully modified.")
} else {
    fmt.Println("⚠️ Order was NOT modified.")
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderModifyExample(ctx context.Context, ticket int32) // sets SL=1.0500, TP=1.0900
func (a *MT4Account) OrderModify(
    ctx context.Context,
    ticket int32,
    price, stoploss, takeprofit *float64,
    expiration *timestamppb.Timestamp,
) (bool, error)
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
| `expiration` | `*timestamppb.Timestamp` | New expiration (pending orders; broker server time). |

`nil` fields are not sent (left unchanged); pass at least one of `price`, `stopLoss`, `takeProfit`, `expiration`.

---

## ⬆️ Output

`ShowOrderModifyExample` prints whether the modification succeeded. `account.OrderModify` returns that flag as `bool` (from `OrderModifyData.OrderWasModified`).

---

## 🎯 Purpose

Adjust SL/TP (and, for pendings, price/expiration) on active orders—useful for dynamic risk management and strategy updates.

---

## 🧩 Notes & Tips

* **Digits & rounding:** Use `Digits` from `SymbolParams` to format/round SL/TP; avoid hardcoded decimals.
* **Volume/price rules:** Brokers enforce minimal distances/steps for SL/TP and pending prices; validate before sending to reduce rejects.
* **No-op guard:** If all optional params are `nil`, the client still sends the request with only the ticket; skip the call instead.

---

## ⚠️ Pitfalls

* **Closed/non-existent ticket:** Will fail with a server-side error (e.g., ticket not found).
* **Invalid levels:** SL above price on sells (or below on buys), TP in the wrong direction, or levels inside the minimal distance will be rejected.
