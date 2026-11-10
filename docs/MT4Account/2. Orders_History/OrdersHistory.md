# OrdersHistory

> Closed-order history — sorted, time-filtered and paged.

**API Information:**

* **Low-level API:** `MT4Account.OrdersHistory(ctx, sortType, from, to, page, itemsPerPage)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **Payload:** `OrdersHistoryData` → `GetOrdersInfo()` → `HistoryOrderInfo`

## 💬 Just the essentials

* **What it is.** The account's closed trades, with sorting, a time window and pagination.
* **Why you need it.** P/L reports, trade journals, analytics over past trades.
* **When to use.** Exports and reports. For a hands-off page loop, use [`OrdersHistoryStream`](../5.%20Streaming_Methods/OrdersHistoryStream.md).

---

## 🎯 Purpose

* Read closed trades with full detail (open/close price, profit, swap)
* Sort by open/close time or ticket, ascending or descending
* Restrict to a `[from, to]` window
* Page through large histories

---

## Method Signature

```go
func (a *MT4Account) OrdersHistory(
    ctx context.Context,
    sortType pb.EnumOrderHistorySortType,
    from, to *time.Time,       // nil = open-ended
    page, itemsPerPage *int32, // nil = server defaults
) (*pb.OrdersHistoryData, error)
```

**Request message:**

```protobuf
message OrdersHistoryRequest {
  EnumOrderHistorySortType input_sort_mode = 1;
  optional Timestamp input_from = 2;      // (+ input_from_utc)
  optional Timestamp input_to = 3;        // (+ input_to_utc)
  optional int32 page_number = 4;
  optional int32 items_per_page = 5;
}
```

**Payload message:**

```protobuf
message OrdersHistoryData {
  repeated HistoryOrderInfo orders_info = 1;
  int32  total_count = 2;
  string next_page_token = 3;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout (history calls get a longer default). |
| `sortType` | `pb.EnumOrderHistorySortType` | yes | Sort field + direction (see enum). |
| `from` | `*time.Time` | no | Window start; `nil` = open-ended. |
| `to` | `*time.Time` | no | Window end; `nil` = open-ended. |
| `page` | `*int32` | no | 1-based page number; `nil` = server default. |
| `itemsPerPage` | `*int32` | no | Page size; `nil` = server default. |

---

## ⬆️ Output — `*pb.OrdersHistoryData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `orders_info` | repeated HistoryOrderInfo | `GetOrdersInfo()` | This page's closed orders |
| `total_count` | int32 | `GetTotalCount() int32` | Total across all pages |
| `next_page_token` | string | `GetNextPageToken() string` | Server paging cursor |

Each `HistoryOrderInfo`: `GetTicket()` `int32`, `GetSymbol()` `string`, `GetOrderType()` `OpenedOrderType`, `GetLots()`, `GetOpenPrice()`, `GetClosePrice()`, `GetProfit()`, `GetSwap()`, `GetCommision()` `float64`, `GetOpenTime()` / `GetCloseTime()` `*Timestamp`, `GetMagicNumber()` `int32`.

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧱 Related enums (from proto)

### `EnumOrderHistorySortType` — input `sortType`

| Constant | Value | Meaning |
|----------|-------|---------|
| `HISTORY_SORT_BY_OPEN_TIME_ASC` | 0 | By open time, oldest first |
| `HISTORY_SORT_BY_OPEN_TIME_DESC` | 1 | By open time, newest first |
| `HISTORY_SORT_BY_CLOSE_TIME_ASC` | 2 | By close time, oldest first |
| `HISTORY_SORT_BY_CLOSE_TIME_DESC` | 3 | By close time, newest first |
| `HISTORY_SORT_BY_ORDER_TICKET_ID_ASC` | 4 | By ticket, ascending |
| `HISTORY_SORT_BY_ORDER_TICKET_ID_DESC` | 5 | By ticket, descending |

> Pass directly: `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Window with `from`/`to`** to bound big histories; `nil` means open-ended on that side.
* **`GetTotalCount()`** lets you compute the number of pages for a chosen `itemsPerPage`.
* **Prefer the stream** ([`OrdersHistoryStream`](../5.%20Streaming_Methods/OrdersHistoryStream.md)) when you just want every page without managing numbers.

---

## 🔗 Usage Examples

### 1) Most recent closed trades

```go
import pb "github.com/MetaRPC/GoMT4/package"

func recentClosed(ctx context.Context, acc *mt4.MT4Account) error {
    page, size := int32(1), int32(20)
    d, err := acc.OrdersHistory(ctx,
        pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC, nil, nil, &page, &size)
    if err != nil { return err }
    for _, o := range d.GetOrdersInfo() {
        fmt.Printf("#%d %s %+.2f\n", o.GetTicket(), o.GetSymbol(), o.GetProfit())
    }
    return nil
}
```

### 2) Realized P/L over the last 30 days

```go
func realized30d(ctx context.Context, acc *mt4.MT4Account) (float64, error) {
    from := time.Now().AddDate(0, 0, -30)
    to := time.Now()
    page, size := int32(1), int32(500)
    var pnl float64
    for {
        d, err := acc.OrdersHistory(ctx,
            pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC, &from, &to, &page, &size)
        if err != nil { return 0, err }
        rows := d.GetOrdersInfo()
        if len(rows) == 0 { break }
        for _, o := range rows { pnl += o.GetProfit() + o.GetSwap() + o.GetCommision() }
        page++
    }
    return pnl, nil
}
```

### 3) How many pages?

```go
func pageCount(ctx context.Context, acc *mt4.MT4Account, size int32) (int32, error) {
    page := int32(1)
    d, err := acc.OrdersHistory(ctx,
        pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC, nil, nil, &page, &size)
    if err != nil { return 0, err }
    total := d.GetTotalCount()
    return (total + size - 1) / size, nil
}
```

---

## ⚠️ Pitfalls

* **Paging ends on an empty page** — stop when `GetOrdersInfo()` is empty.
* **Timestamps.** `from`/`to` are `time.Time`; mind the broker's time zone vs UTC.
* **Large pulls.** Use a sensible `itemsPerPage` and window, or the stream.

---

## 🧪 Testing Suggestions

* **Happy path (demo with history).** Page 1 returns rows; `GetTotalCount() > 0`.
* **Paging.** Walk pages until empty; the summed row count equals `GetTotalCount()`.
* **Window.** A 1-day `from`/`to` returns fewer rows than open-ended.

---

## 📚 See also

* [OrdersHistoryStream](../5.%20Streaming_Methods/OrdersHistoryStream.md) · [OpenedOrders](./OpenedOrders.md) · [OrderSelect](./OrderSelect.md)
* **Back to:** [2. Orders & History — Overview](./Orders_History.Overview.md) · [MT4Account — Overview](../Overview.md)
