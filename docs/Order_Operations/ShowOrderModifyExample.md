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

| Field        | Type         | Description                         |
| ------------ | ------------ | ----------------------------------- |
| `price`      | `*float64`   | Specific price for modification.    |
| `stopLoss`   | `*float64`   | New stop loss level.                |
| `takeProfit` | `*float64`   | New take profit level.              |
| `expiration` | `*time.Time` | Expiration time for pending orders. |

The provided ticket must be a valid active order ID; otherwise, the server will return an error such as `Invalid ticket` or `Ticket not found`.

---

## ⬆️ Output

The method prints modification result to stdout, indicating:

* Whether the modification was successful.
* Relevant server feedback or error messages in case of failure.

---

## 🎯 Purpose

This method allows manual updates to critical trade parameters (SL/TP) of active orders, useful for:

* Dynamic adjustment of risk management
* Trade strategy automation
* Manual intervention based on changing market conditions

Ensure you use valid parameters and active tickets for successful order modifications.
