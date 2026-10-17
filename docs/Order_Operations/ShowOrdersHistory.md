# Getting Order History

> **Request:** retrieve historical orders for a specified time range
> Fetch orders from the trading account history (closed positions and deleted pending orders) whose **open time** falls within a defined window.

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
    nil, nil /* page & itemsPerPage (optional) */,
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
func (s *MT4Service) ShowOrdersHistory(ctx context.Context) // last 7 days, sorted by close time DESC
func (a *MT4Account) OrdersHistory(
    ctx context.Context,
    sortType pb.EnumOrderHistorySortType,
    from, to *time.Time,
    page, itemsPerPage *int32,
) (*pb.OrdersHistoryData, error)
```

---

## 🔽 Input

Required:

* `ctx` (`context.Context`) — context for managing timeout or cancellation.

Method internally uses (when calling `MT4Account.OrdersHistory`):

| Field      | Type                          | Description                       |
| ---------- | ----------------------------- | --------------------------------- |
| `sortType` | `pb.EnumOrderHistorySortType` | Sorting logic.                    |
| `from`     | `*time.Time`                  | Start of the window of order **open** times (broker server time; `nil` = no lower bound). |
| `to`       | `*time.Time`                  | End of the window of order **open** times (broker server time; `nil` = no upper bound). |
| `page`     | `*int32`                      | Page number, 1-based (optional; needs `itemsPerPage`). |
| `itemsPerPage` | `*int32`                  | Items per page (optional; needs `page`). |

Possible `pb.EnumOrderHistorySortType` values (Go constants `pb.EnumOrderHistorySortType_<value>`):

* `HISTORY_SORT_BY_OPEN_TIME_ASC`
* `HISTORY_SORT_BY_OPEN_TIME_DESC`
* `HISTORY_SORT_BY_CLOSE_TIME_ASC`
* `HISTORY_SORT_BY_CLOSE_TIME_DESC`
* `HISTORY_SORT_BY_ORDER_TICKET_ID_ASC`
* `HISTORY_SORT_BY_ORDER_TICKET_ID_DESC`

---

## ⬆️ Output

Prints closed order information to stdout.
Underlying response: `*pb.OrdersHistoryData`

| Field           | Type                     | Description                         |
| --------------- | ------------------------ | ----------------------------------- |
| `OrdersInfo`    | `[]*pb.HistoryOrderInfo` | List of historical orders.          |
| `TotalCount`    | `int32`                  | Total number of matching orders.    |
| `NextPageToken` | `string`                 | Cursor for the next page (cursor paging, not used by this client). |

Each `*pb.HistoryOrderInfo` includes:

| Field        | Type                 | Description                           |
| ------------ | -------------------- | ------------------------------------- |
| `Ticket`     | `int32`              | Unique ID of the order.               |
| `Symbol`     | `string`             | Trading symbol (e.g., EURUSD).        |
| `Lots`       | `float64`            | Volume of the order in lots.          |
| `OpenPrice`  | `float64`            | Entry price of the order.             |
| `ClosePrice` | `float64`            | Exit price of the order.              |
| `Profit`     | `float64`            | Final realized profit/loss.           |
| `Swap` / `Commision` | `float64`    | Swap and commission.                  |
| `OrderType`  | `pb.OpenedOrderType` | Type of order (`OO_OP_BUY`, `OO_OP_SELL`, …). |
| `OpenTime` / `OpenTimeUtc`   | `*timestamppb.Timestamp` | Open time (broker server time / UTC). |
| `CloseTime` / `CloseTimeUtc` | `*timestamppb.Timestamp` | Close time (broker server time / UTC). |
| `StopLoss`   | `float64`            | Stop Loss price (if set).             |
| `TakeProfit` | `float64`            | Take Profit price (if set).           |
| `MagicNumber`| `int32`              | Magic number for programmatic orders. |
| `Comment`    | `string`             | Custom comment attached to the order. |
| `ExpirationTime` | `*timestamppb.Timestamp` | Expiration time for pending orders. |

---

## 🧩 Notes & Tips

* **Pagination:** Use `page` and `itemsPerPage` for large ranges to avoid oversized payloads.
* **Sorting:** Choose `sortType` based on how you plan to display/export results (by open vs close time).
* **Time window:** The window filters by **open** time; `nil` means no bound on that side. The values are read as broker server time.
* **Paging helper:** `account.OrdersHistoryStream(ctx, sortType, from, to, pageSize)` walks the pages for you ([StreamOrdersHistoryExample](../Streaming/StreamOrdersHistoryExample.md)).

---

## ⚠️ Pitfalls

* **Wide ranges:** Very large windows can be slow/heavy; prefer paged requests.
* **Broker discrepancies:** History can differ slightly across servers/brokers for the same symbol.
* **Orders, not deals:** MT4 has no separate deals; each history order is one closed position or deleted pending order.

---

## 🎯 Purpose

Retrieve completed trades in a specified time frame for:

* Historical trade analysis and reporting
* Auditing and compliance exports
* Reconciliation with external systems
