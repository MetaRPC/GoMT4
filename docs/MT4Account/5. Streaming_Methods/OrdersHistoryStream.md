# OrdersHistoryStream

> Walk the full order history page by page, delivered through a channel — no manual paging loop.

**API Information:**

* **Low-level API:** `MT4Account.OrdersHistoryStream(ctx, sortType, from, to, pageSize)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **Built on:** [`OrdersHistory`](../2.%20Orders_History/OrdersHistory.md) — a client-side pager, **not** a live server stream. It calls `OrdersHistory` for each page and feeds the pages down a channel until the history is exhausted.
* **Payload:** `OrdersHistoryData` (one per page)

## 💬 Just the essentials

* **What it is.** A convenience wrapper that pages `OrdersHistory` for you and sends every page down a channel.
* **Why you need it.** Process large histories without tracking page numbers.
* **When to use.** Exports, backfills, analytics over closed orders.

---

## 🎯 Purpose

* Stream the whole closed-order history page by page
* Avoid writing a manual pagination loop
* Stop early by cancelling the context
* Reuse the same sort/window options as `OrdersHistory`

---

## Method Signature

```go
func (a *MT4Account) OrdersHistoryStream(
    ctx context.Context,
    sortType pb.EnumOrderHistorySortType,
    from, to *time.Time,   // nil = open-ended
    pageSize int32,
) (<-chan *pb.OrdersHistoryData, <-chan error)
```

It internally increments the page and calls [`OrdersHistory`](../2.%20Orders_History/OrdersHistory.md) with your `sortType` / `from` / `to` / `pageSize`, stopping when a page comes back empty.

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to stop paging early and close the channels. |
| `sortType` | `pb.EnumOrderHistorySortType` | yes | Sort order (see [`OrdersHistory`](../2.%20Orders_History/OrdersHistory.md#related-enums-from-proto)). |
| `from` | `*time.Time` | no | Range start; `nil` = open-ended. |
| `to` | `*time.Time` | no | Range end; `nil` = open-ended. |
| `pageSize` | `int32` | yes | Orders per page the wrapper requests. |

---

## ⬆️ Frame — `*pb.OrdersHistoryData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `orders_info` | repeated HistoryOrderInfo | `GetOrdersInfo()` | This page's orders |
| `total_count` | int32 | `GetTotalCount() int32` | Total across all pages |
| `next_page_token` | string | `GetNextPageToken() string` | Server paging cursor |

> The data channel is **closed** (you receive `nil`) when the history is exhausted. A non-nil value on the error channel ends it early.

---

## 🧩 Notes & Tips

* **Ends on an empty page** — the wrapper stops and closes the channel.
* **Not a live stream** — it finishes; it won't push future trades (use [`OnTrade`](./OnTrade.md)).
* **`cancel()` stops early** — safe to abandon mid-history.

---

## 🔗 Usage Examples

### 1) Count all closed orders, page by page

```go
import pb "github.com/MetaRPC/GoMT4/package"

func countHistory(ctx context.Context, acc *mt4.MT4Account) (int, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    pages, errs := acc.OrdersHistoryStream(ctx,
        pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC, nil, nil, 100)
    total := 0
    for {
        select {
        case p := <-pages:
            if p == nil { return total, nil } // done
            total += len(p.GetOrdersInfo())
        case err := <-errs:
            if err != nil { return 0, err }
        case <-ctx.Done():
            return total, nil
        }
    }
}
```

### 2) Sum realized profit over the last 30 days

```go
func realized30d(ctx context.Context, acc *mt4.MT4Account) (float64, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    from, to := time.Now().AddDate(0, 0, -30), time.Now()
    pages, errs := acc.OrdersHistoryStream(ctx,
        pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC, &from, &to, 200)
    var pnl float64
    for {
        select {
        case p := <-pages:
            if p == nil { return pnl, nil }
            for _, o := range p.GetOrdersInfo() { pnl += o.GetProfit() }
        case err := <-errs:
            if err != nil { return 0, err }
        case <-ctx.Done():
            return pnl, nil
        }
    }
}
```

### 3) Collect all into one slice

```go
func allHistory(ctx context.Context, acc *mt4.MT4Account) ([]*pb.HistoryOrderInfo, error) {
    pages, errs := acc.OrdersHistoryStream(ctx,
        pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC, nil, nil, 500)
    var all []*pb.HistoryOrderInfo
    for {
        select {
        case p := <-pages:
            if p == nil { return all, nil }
            all = append(all, p.GetOrdersInfo()...)
        case err := <-errs:
            if err != nil { return nil, err }
        }
    }
}
```

---

## ⚠️ Pitfalls

* **Channel close = `nil`** — check for it to know you're done.
* **Big histories.** A large `pageSize` is fewer round-trips but heavier pages — tune it.
* **Not live** — it won't include trades closed after it starts.

---

## 🧪 Testing Suggestions

* **Happy path (demo with history).** The summed page counts equal `GetTotalCount()`.
* **Window.** A 1-day range yields fewer pages than open-ended.
* **Cancel.** Cancelling mid-stream stops promptly.

---

## 📚 See also

* [OrdersHistory](../2.%20Orders_History/OrdersHistory.md) · [OnTrade](./OnTrade.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
