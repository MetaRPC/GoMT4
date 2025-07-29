# Closing an Order

> **Request:** close or delete an open order by its ticket
> Sends a request to the server to close or delete the specified order.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderCloseExample(context.Background(), 123456)

// Or directly from MT4Account
result, err := mt4.OrderClose(context.Background(), 123456, nil, nil, nil)
if err != nil {
    log.Fatalf("Error closing order: %v", err)
}

fmt.Printf("Closed: %s, Comment: %s\n", result.GetMode(), result.GetHistoryOrderComment())
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderCloseExample(ctx context.Context, ticket int32)
```

---

## 🔽 Input

Required:

* `ctx` (`context.Context`) — context for managing timeout or cancellation.
* `ticket` (`int32`) — The ticket number of the order to be closed.

Optional (if directly using MT4Account):

* `price` (`*float64`) — specific closing price.
* `slippage` (`*int32`) — acceptable slippage.
* `magic` (`*int32`) — magic number for identification.

The provided ticket must be a valid active order ID; otherwise, the server will return an error such as `Invalid ticket` or `Ticket not found`.

---

## ⬆️ Output

Prints result information to stdout and returns:

| Field                 | Type     | Description                                        |
| --------------------- | -------- | -------------------------------------------------- |
| `Mode`                | `string` | Operation mode result (e.g., "Closed", "Deleted"). |
| `HistoryOrderComment` | `string` | Server comment describing the result.              |

---

## 🎯 Purpose

This method enables manual closing or deletion of orders by ticket, useful for:

* Manual trade intervention or debugging tools
* Post-trade processing or cleanup
* Workflow simulations and order closure testing

---

### ❓ Why it’s commented out in code:

This method requires a **valid, active ticket number**. It’s often commented out by default to:

* ❌ Prevent runtime errors from invalid tickets
* ✅ Ensure intentional use with real ticket IDs

Test using valid tickets obtained from methods like `OpenedOrders` or `OpenedOrdersTickets`.
