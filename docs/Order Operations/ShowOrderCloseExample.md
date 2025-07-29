# Sending a Market or Pending Order

> **Request:** send a trade order (market or pending)
> Sends a new order using the specified parameters and receives back execution details.

---

### Code Example

```go
// Using service wrapper
service.ShowOrderSendExample(context.Background(), "EURUSD")

// Or directly using MT4Account
result, err := mt4.OrderSend(
    context.Background(),
    "EURUSD",
    pb.OrderSendOperationType_OC_OP_BUY,
    0.1,
    nil,
    ptrInt32(5),
    ptrFloat64(1.0500),
    ptrFloat64(1.0900),
    ptrString("Go order test"),
    ptrInt32(123456),
    nil,
)

if err != nil {
    log.Fatalf("Error sending order: %v", err)
}

fmt.Printf("Order opened! Ticket: %d, Price: %.5f, Time: %s\n",
    result.GetTicket(),
    result.GetPrice(),
    result.GetOpenTime().AsTime().Format("2006-01-02 15:04:05"),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderSendExample(ctx context.Context, symbol string)
```

---

## 🔽 Input

Required:

| Field       | Type                        | Description                         |
| ----------- | --------------------------- | ----------------------------------- |
| `ctx`       | `context.Context`           | Timeout or cancellation management. |
| `symbol`    | `string`                    | Trading symbol (e.g., "EURUSD").    |
| `orderType` | `pb.OrderSendOperationType` | Type of order (market/pending).     |
| `volume`    | `float64`                   | Order volume in lots (e.g., 0.1).   |

Optional parameters (when directly using MT4Account):

| Field        | Type         | Description                         |
| ------------ | ------------ | ----------------------------------- |
| `price`      | `*float64`   | Order price for pending orders.     |
| `slippage`   | `*int32`     | Max slippage allowed (points).      |
| `stopLoss`   | `*float64`   | Stop Loss price.                    |
| `takeProfit` | `*float64`   | Take Profit price.                  |
| `comment`    | `*string`    | Optional order comment.             |
| `magic`      | `*int32`     | Magic number to tag the order.      |
| `expiration` | `*time.Time` | Expiration time for pending orders. |

---

## ⬆️ Output

Returns execution details:

| Field      | Type        | Description                         |
| ---------- | ----------- | ----------------------------------- |
| `Ticket`   | `int32`     | Unique order ID assigned by MT4.    |
| `Price`    | `float64`   | Actual execution price.             |
| `OpenTime` | `timestamp` | Time when order was executed (UTC). |

---

## 🎯 Purpose

Use this method to place new trade orders (market or pending), controlling volume, price, and risk parameters. The result confirms the assigned ticket number, price, and open time for tracking or logging.

---

### ❓ Notes

This method requires:

* Valid trading symbol
* Stable connection to MT4 terminal
* Terminal in trading-enabled state

Ensure proper configuration before using in production to avoid unexpected errors.
