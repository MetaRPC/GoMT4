# 🚀 Streaming Orders History (Example)

> **Request:** Stream historical orders with sensible defaults.
> This demo wrapper prints pages of order history and exits after a short timeout.

---

### Code Example

```go
// Quick use: prints orders from the last 30 days in pages (size=200)
// sorted by close time (DESC). Demo loop times out after ~30s.
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
* Demo timeout: \~30 seconds (select-case exit)

Under the hood, it calls the service API:

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

> Need more detail (open/close prices, times, lots, commission, etc.)? Use the low-level `StreamOrdersHistory` and print the extra fields you need. The proto model exposes them on each order.

---

## ENUM: `pb.EnumOrderHistorySortType`

Common values (check your generated `pb` for exact names):

* `HISTORY_SORT_BY_CLOSE_TIME_DESC`
* `HISTORY_SORT_BY_CLOSE_TIME_ASC`
* (others may exist depending on your schema)

**Order type enum:** depends on your `pb` (e.g., `ORDER_TYPE_BUY`, `ORDER_TYPE_SELL`, `ORDER_TYPE_BUY_LIMIT`, `ORDER_TYPE_SELL_LIMIT`, `ORDER_TYPE_BUY_STOP`, `ORDER_TYPE_SELL_STOP`).

---

## 🧩 Notes & Tips

* **Connection required:** Ensure you have called `ConnectByServerName`/`ConnectByHostPort` successfully; otherwise the underlying stream will fail fast.
* **Auto-reconnect:** The stream layer retries on transient `gRPC`/API errors with backoff + jitter.
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
