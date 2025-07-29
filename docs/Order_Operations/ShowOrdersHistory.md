# Getting Order History

> **Request:** retrieve historical orders for a specified time range
> Fetch all closed orders from the trading account history within a defined window.

---

### Code Example

```go
// Using service wrapper
service.ShowOrdersHistory(context.Background())

// Or directly from MT4Account
from := time.Now().AddDate(0, 0, -7)
to := time.Now()

history, err := mt4.OrdersHistory(
    context.Background(),
    pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
    &from, &to, nil, nil,
)

if err != nil {
    log.Fatalf("Error retrieving order history: %v", err)
}

for _, order := range history.GetOrdersInfo() {
    fmt.Printf("[%s] Ticket: %d | Symbol: %s | Lots: %.2f | Open: %.5f | Close: %.5f | Profit: %.2f | Closed: %s\n",
        order.GetOrderType(),
        order.GetTicket(),
        order.GetSymbol(),
        order.GetLots(),
        order.GetOpenPrice(),
        order.GetClosePrice(),
        order.GetProfit(),
        order.GetCloseTime().AsTime().Format("2006-01-02 15:04:05"),
    )
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrdersHistory(ctx context.Context)
```

---

## 🔽 Input

Required:

* `ctx` (`context.Context`) — context for managing timeout or cancellation.

Method internally uses:

| Field      | Type                       | Description                       |
| ---------- | -------------------------- | --------------------------------- |
| `sortType` | `EnumOrderHistorySortType` | Sorting logic.                    |
| `from`     | `time.Time`                | Start time of the history window. |
| `to`       | `time.Time`                | End time of the history window.   |

Possible `EnumOrderHistorySortType` values:

* `HISTORY_SORT_BY_OPEN_TIME_ASC`
* `HISTORY_SORT_BY_OPEN_TIME_DESC`
* `HISTORY_SORT_BY_CLOSE_TIME_ASC`
* `HISTORY_SORT_BY_CLOSE_TIME_DESC`

Optional:

* `deadline` (`time.Time`) — optional timeout.
* `cancellationToken` (via context) — optional cancellation control.

---

## ⬆️ Output

Prints closed order information to stdout.

Underlying returned structure:

| Field        | Type          | Description                         |
| ------------ | ------------- | ----------------------------------- |
| `OrdersInfo` | `[]OrderInfo` | List of historical (closed) orders. |

Each `OrderInfo` includes:

| Field        | Type                 | Description                           |
| ------------ | -------------------- | ------------------------------------- |
| `Ticket`     | `int32`              | Unique ID of the order.               |
| `Symbol`     | `string`             | Trading symbol (e.g., EURUSD).        |
| `Lots`       | `float64`            | Volume of the order in lots.          |
| `OpenPrice`  | `float64`            | Entry price of the order.             |
| `ClosePrice` | `float64`            | Exit price of the order.              |
| `Profit`     | `float64`            | Final realized profit/loss.           |
| `OrderType`  | `ENUM_ORDER_TYPE_TF` | Type of order (Buy, Sell, etc.).      |
| `OpenTime`   | `timestamp`          | Time when the order was opened.       |
| `CloseTime`  | `timestamp`          | Time when the order was closed.       |
| `Sl`         | `float64`            | Stop Loss price (if set).             |
| `Tp`         | `float64`            | Take Profit price (if set).           |
| `Magic`      | `int32`              | Magic number for programmatic orders. |
| `Comment`    | `string`             | Custom comment attached to the order. |
| `Expiration` | `timestamp`          | Expiration time for pending orders.   |

---

### ENUM: `ENUM_ORDER_TYPE_TF`

| Value                  | Description        |
| ---------------------- | ------------------ |
| `OrderTypeTfBuy`       | Buy order          |
| `OrderTypeTfSell`      | Sell order         |
| `OrderTypeTfBuyLimit`  | Pending Buy Limit  |
| `OrderTypeTfSellLimit` | Pending Sell Limit |
| `OrderTypeTfBuyStop`   | Pending Buy Stop   |
| `OrderTypeTfSellStop`  | Pending Sell Stop  |

---

## 🎯 Purpose

This method retrieves completed trades within a specified time frame, useful for:

* Historical trade analysis
* Auditing or reporting
* Exporting trade logs for compliance or analytics

Essential for accessing historical trade data from MT4.
