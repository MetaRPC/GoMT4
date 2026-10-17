# Deleting a Pending Order

> **Request:** delete a pending order by its ticket
> Sends a request to remove a pending order from the trading terminal.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Deletes a pending order by ticket.
// ⚠️ Works only for pending orders (BuyLimit, SellLimit, etc.).
svc.ShowOrderDeleteExample(ctx, 123456)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.
// ⚠️ OrderDelete sends the same request as OrderClose: on a market Buy/Sell it CLOSES the position.

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := account.OrderDelete(ctx, 123456)
if err != nil {
    log.Fatalf("❌ OrderDelete error: %v", err)
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
func (a *MT4Account) OrderDelete(ctx context.Context, ticket int32) (*pb.OrderCloseDeleteData, error)
```

`OrderDelete(ctx, ticket)` is `OrderClose(ctx, ticket, nil, nil, nil)`.

---

## 🔽 Input

Required:

| Field    | Type              | Description                         |
| -------- | ----------------- | ----------------------------------- |
| `ctx`    | `context.Context` | Timeout or cancellation management. |
| `ticket` | `int32`           | Ticket number of the pending order. |

The ticket should reference a **pending** order — check its type (`OrderSelect` / `OpenedOrders`) first, because a market order with that ticket would be closed.

---

## ⬆️ Output

Returns `*pb.OrderCloseDeleteData`:

| Field                 | Type     | Description                             |
| --------------------- | -------- | --------------------------------------- |
| `Mode`                | `pb.OrderCloseDeleteMode` | `OCD_PENDING_ORDER` for a deleted pending order (`OCD_MARKET_ORDER` if a market order was closed). |
| `HistoryOrderComment` | `*string` | Comment of the resulting history order. |

---

## 🎯 Purpose

Remove limit/stop orders that are no longer needed:

* Cancel pending orders before execution
* Manage pending order queues
* Clean up unused test/training orders

---

## 🧩 Notes & Tips

* **Pending only (by convention):** For open market positions use `OrderClose`; `OrderDelete` would close them too, since both send `OrderCloseDelete`.
* **Idempotency:** If the order is already filled/expired/deleted, the API returns an error.
* **No price/slippage:** Delete uses only the ticket; it sends no lots, price or slippage.
