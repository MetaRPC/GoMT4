# 🗂️ HistoryOrders (GoMT4)

**Goal:** load **orders history** for a time range (optionally filtered by symbol), with practical paging and export examples.

> Real code refs in this repo:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrdersHistory`, `OrdersHistoryStream`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrdersHistory`, `StreamOrdersHistoryExample`)

---

## ✅ 1) Preconditions

* `account` is connected.
* `config.json` is valid and `DefaultSymbol` exists.

---

## 🗓️ 2) Basic range query (last 7 days)

The range filters by order **open** time and is read as broker server time; `nil` means no bound.

```go
to := time.Now()
from := to.AddDate(0, 0, -7)

hist, err := account.OrdersHistory(ctx,
    pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
    &from, &to,
    nil, nil /* page & itemsPerPage: all at once */,
)
if err != nil { return err }

for _, o := range hist.GetOrdersInfo() {
    fmt.Printf("%s ticket=%d symbol=%s lots=%.2f profit=%.2f open=%s close=%s\n",
        o.GetOrderType().String(), o.GetTicket(), o.GetSymbol(), o.GetLots(), o.GetProfit(),
        o.GetOpenTime().AsTime().Format(time.RFC3339),
        o.GetCloseTime().AsTime().Format(time.RFC3339))
}
```

---

## 🎯 3) Filter by symbol

The request has no symbol filter; filter on the client:

```go
sym := "EURUSD"
to := time.Now()
from := to.AddDate(0, 0, -30)

hist, err := account.OrdersHistory(ctx, pb.EnumOrderHistorySortType_HISTORY_SORT_BY_OPEN_TIME_ASC, &from, &to, nil, nil)
if err != nil { return err }

var mine []*pb.HistoryOrderInfo
for _, o := range hist.GetOrdersInfo() {
    if o.GetSymbol() == sym {
        mine = append(mine, o)
    }
}
fmt.Printf("%d history orders for %s\n", len(mine), sym)
```

---

## 📦 4) Page through long ranges

Long queries can be heavy. `OrdersHistoryStream` requests page 1, 2, 3, … until a page comes back empty:

```go
func loadHistoryPaged(ctx context.Context, a *mt4.MT4Account, days int) ([]*pb.HistoryOrderInfo, error) {
    to := time.Now()
    from := to.AddDate(0, 0, -days)

    pagesCh, errCh := a.OrdersHistoryStream(ctx, pb.EnumOrderHistorySortType_HISTORY_SORT_BY_OPEN_TIME_ASC, &from, &to, 500)

    var out []*pb.HistoryOrderInfo
    for page := range pagesCh {
        out = append(out, page.GetOrdersInfo()...)
    }
    if err, ok := <-errCh; ok && err != nil {
        return out, err
    }
    return out, nil
}
```

---

## 📤 5) Export to CSV (quick ‘n’ dirty)

```go
func saveCSV(path string, orders []*pb.HistoryOrderInfo) error {
    f, err := os.Create(path)
    if err != nil { return err }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    _ = w.Write([]string{"ticket","symbol","type","lots","profit","open","close"})
    for _, o := range orders {
        _ = w.Write([]string{
            strconv.Itoa(int(o.GetTicket())),
            o.GetSymbol(),
            o.GetOrderType().String(),
            fmt.Sprintf("%.2f", o.GetLots()),
            fmt.Sprintf("%.2f", o.GetProfit()),
            o.GetOpenTime().AsTime().Format(time.RFC3339),
            o.GetCloseTime().AsTime().Format(time.RFC3339),
        })
    }
    w.Flush()
    return w.Error()
}
```

---

## ⚠️ Pitfalls

* **Open-time window** → an order opened before `from` and closed inside the window is **not** returned; widen `from` if you need everything closed in a period.
* **Timezones** → `OpenTime` / `CloseTime` are broker server time; `OpenTimeUtc` / `CloseTimeUtc` are UTC.
* **Performance** → prefer paging for 30+ days ranges.
* **Context** → 8s default timeout per call when `ctx` has no deadline; transport errors are retried inside the client (see *Reliability*).

---

## 🔗 See also

* [CloseOrder](CloseOrder.md) — verify closed positions appear in history.
* [ModifyOrder](ModifyOrder.md) — changes in SL/TP reflected in history records.
* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — timeouts, reconnects & backoff patterns.
