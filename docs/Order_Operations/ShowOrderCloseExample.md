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
    res.GetMode(), // OCD_MARKET_ORDER (closed) or OCD_PENDING_ORDER (deleted)
    res.GetHistoryOrderComment(),
)

// Partial close: 0.05 lot at market with up to 10 points slippage
lots := 0.05
slip := int32(10)
_, err = account.OrderClose(ctx, 123456, &lots, nil, &slip)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderCloseExample(ctx context.Context, ticket int32)
func (a *MT4Account) OrderClose(ctx context.Context, ticket int32, lots, price *float64, slippage *int32) (*pb.OrderCloseDeleteData, error)
```

---

## 🔽 Input

| Field    | Type              | Description                     |
| -------- | ----------------- | ------------------------------- |
| `ctx`    | `context.Context` | Timeout / cancellation control. |
| `ticket` | `int32`           | Ticket of the order to close.   |
| `lots`   | `*float64`        | Optional volume for a partial close (`nil` = whole order). |
| `price`  | `*float64`        | Optional close price (`nil` = market). |
| `slippage` | `*int32`        | Optional slippage (points).     |

`ShowOrderCloseExample` passes `nil` for all three.

---

## ⬆️ Output

Returns `*pb.OrderCloseDeleteData`:

| Field                 | Type                      | Description                              |
| --------------------- | ------------------------- | ---------------------------------------- |
| `Mode`                | `pb.OrderCloseDeleteMode` | `OCD_MARKET_ORDER` (a market order was closed) or `OCD_PENDING_ORDER` (a pending order was deleted). |
| `HistoryOrderComment` | `*string`                 | Comment of the resulting history order (use `GetHistoryOrderComment()`). |

---

## 🎯 Purpose

Close or delete an order by ticket (one `OrderCloseDelete` request handles both: market orders are closed, pending orders are deleted).
Useful for manual interventions, post-trade cleanup, or testing order workflows. Default timeout: 5s when your `ctx` has no deadline.
