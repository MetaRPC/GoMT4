# 🚀 Streaming Orders History (Example)

> **Request:** Stream historical orders with sensible defaults.
> This demo wrapper prints pages of order history and exits after a short timeout.

---

### Code Example

```go
// Quick use: prints orders opened in the last 30 days in pages (size=200)
// sorted by close time (DESC). The demo loop stops after 3s without a page.
svc.StreamOrdersHistoryExample(ctx)
```

---

### Method Signature

```go
// Demo wrapper: safe defaults + console printing
func (s *MT4Service) StreamOrdersHistoryExample(ctx context.Context)
```

**Defaults used inside the example wrapper:**

* Date range: `from = now - 30 days`, `to = now`
* Sort: `HISTORY_SORT_BY_CLOSE_TIME_DESC`
* Page size: `200`
* Demo timeout: 3 seconds without a new page (select-case exit)

Under the hood, it calls the service API, which proxies `MT4Account.OrdersHistoryStream`:

```go
func (s *MT4Service) StreamOrdersHistory(
    ctx context.Context,
    sortType pb.EnumOrderHistorySortType,
    from, to *time.Time,
    pageSize int32,
) (<-chan *pb.OrdersHistoryData, <-chan error)
```

---

## 🔽 Input

| Field | Type              | Required | Description                                  |
| ----: | ----------------- | :------: | -------------------------------------------- |
| `ctx` | `context.Context` |     ✅    | Controls lifetime (cancel/timeout) of stream |

> No other parameters are required for the **Example** wrapper; it sets safe defaults.

---

## ⬆️ Output

**Console output** of selected fields from each page of `*pb.OrdersHistoryData`:

For each order (`OrdersInfo[]`), the example prints:

* `OrderType`
* `Ticket`
* `Symbol`
* `Profit`

Sample line format:

```
[HIST] <OrderType> | Ticket: <id> | <Symbol> | PnL: <profit>
```

> Need more detail (open/close prices, times, lots, commission, etc.)? Use `StreamOrdersHistory` (or `account.OrdersHistoryStream`) and print the extra fields you need from each `*pb.HistoryOrderInfo`.

How the paging works: `OrdersHistoryStream` calls `OrdersHistory` with page 1, 2, 3, … (`itemsPerPage = pageSize`) and sends each page to the channel; it stops and closes both channels at the first empty page or on an error. `from`/`to` filter by order **open** time and are read as broker server time.

---

## ENUM: `pb.EnumOrderHistorySortType`

Values (Go constants `pb.EnumOrderHistorySortType_<value>`):

* `HISTORY_SORT_BY_OPEN_TIME_ASC`, `HISTORY_SORT_BY_OPEN_TIME_DESC`
* `HISTORY_SORT_BY_CLOSE_TIME_ASC`, `HISTORY_SORT_BY_CLOSE_TIME_DESC`
* `HISTORY_SORT_BY_ORDER_TICKET_ID_ASC`, `HISTORY_SORT_BY_ORDER_TICKET_ID_DESC`

**Order type enum:** `pb.OpenedOrderType` (`OpenedOrderType_OO_OP_BUY`, `_OO_OP_SELL`, `_OO_OP_BUYLIMIT`, `_OO_OP_SELLLIMIT`, `_OO_OP_BUYSTOP`, `_OO_OP_SELLSTOP`).

---

## 🧩 Notes & Tips

* **Connection required:** Ensure you have called `ConnectByServerName`/`ConnectByHostPort` successfully; otherwise the first page request fails and the error arrives on `errCh`.
* **Auto-reconnect:** Each page request is a unary `OrdersHistory` call that retries on transient gRPC errors and terminal-not-found errors with backoff + jitter.
* **Customization:** For production, prefer `StreamOrdersHistory` (service API) and pass your own `from/to`, sorting, and `pageSize`.
* **Throughput:** If you process heavy pages, do it in a worker to avoid blocking the reader loop.

---

## ⚠️ Pitfalls

* **Ignoring `errCh`:** Always `select` on the error channel; otherwise you’ll miss terminal failures.
* **Infinite loops:** The example has a demo timeout; if you remove it, make sure you still have a cancellation path.
* **Large ranges:** Very wide date ranges + small pages can be slow; increase `pageSize` as appropriate.

---

## 🧪 Testing

* **Happy path:** Expect multiple pages if you have many orders in the last 30 days.
* **Network hiccups:** Brief drops should be retried transparently; you may notice short pauses.
* **Cancel path:** Cancel the passed `ctx`; both channels should close and the goroutine should exit.
