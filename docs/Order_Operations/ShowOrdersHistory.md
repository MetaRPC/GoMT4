# Getting Order History

> **Request:** retrieve historical orders for a specified time range
> Fetch all closed orders from the trading account history within a defined window.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints order history for the last 7 days.
svc.ShowOrdersHistory(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

from := time.Now().AddDate(0, 0, -7)
to   := time.Now()

ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
defer cancel()

history, err := account.OrdersHistory(
    ctx,
    pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
    &from, &to,
    nil, nil, // page & itemsPerPage (optional)
)
if err != nil {
    log.Fatalf("❌ OrdersHistory error: %v", err)
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

Method internally uses (when calling `MT4Account.OrdersHistory`):

| Field      | Type                          | Description                       |
| ---------- | ----------------------------- | --------------------------------- |
| `sortType` | `pb.EnumOrderHistorySortType` | Sorting logic.                    |
| `from`     | `*time.Time`                  | Start time of the history window. |
| `to`       | `*time.Time`                  | End time of the history window.   |
| `page`     | `*int32`                      | Page number (optional).           |
| `items`    | `*int32`                      | Items per page (optional).        |

Possible `EnumOrderHistorySortType` values:

* `HISTORY_SORT_BY_OPEN_TIME_ASC`
* `HISTORY_SORT_BY_OPEN_TIME_DESC`
* `HISTORY_SORT_BY_CLOSE_TIME_ASC`
* `HISTORY_SORT_BY_CLOSE_TIME_DESC`

---

## ⬆️ Output

Prints closed order information to stdout.
Underlying response: `*pb.OrdersHistoryData`

| Field        | Type    | Description                         |
| ------------ | ------- | ----------------------------------- |
| `OrdersInfo` | (slice) | List of historical (closed) orders. |

Printed fields typically include:

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

## 🧩 Notes & Tips

* **Pagination:** Use `page` and `itemsPerPage` for large ranges to avoid oversized payloads.
* **Sorting:** Choose `sortType` based on how you plan to display/export results (by open vs close time).
* **Time window:** If `from`/`to` are nil, the server may apply defaults; pass both for deterministic results.

---

## ⚠️ Pitfalls

* **Wide ranges:** Very large windows can be slow/heavy; prefer paged requests.
* **Broker discrepancies:** History can differ slightly across servers/brokers for the same symbol.
* **Status vs fills:** This is **order** history; if you need individual fills/deals, use the corresponding deals endpoint.

---

## 🎯 Purpose

Retrieve completed trades in a specified time frame for:

* Historical trade analysis and reporting
* Auditing and compliance exports
* Reconciliation with external systems
