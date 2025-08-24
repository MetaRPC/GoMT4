# 🗂️ HistoryOrders (GoMT4)

**Goal:** load **orders history** for a time range (optionally filtered by symbol), with practical batching and export examples.

> Real code refs in this repo:
>
> * Account methods: `examples/mt4/MT4Account.go` (e.g., `ShowOrdersHistory`, `OrdersHistoryRange`/similar)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrdersHistory` demo)

---

## ✅ 1) Preconditions

* MT4 terminal is connected and has the history for requested dates.
* `config.json` is valid and `DefaultSymbol` exists.

---

## 🗓️ 2) Basic range query (last 7 days)

```go
from := timestamppb.New(time.Now().AddDate(0, 0, -7))
to   := timestamppb.New(time.Now())

hist, err := account.ShowOrdersHistory(ctx, from, to, nil) // nil = all symbols
if err != nil { return err }

for _, o := range hist.GetOrders() {
    fmt.Printf("%s ticket=%d symbol=%s lots=%.2f profit=%.2f open=%s close=%s\n",
        o.GetType().String(), o.GetTicket(), o.GetSymbol(), o.GetVolume(), o.GetProfit(),
        o.GetOpenTime().AsTime().Format(time.RFC3339),
        o.GetCloseTime().AsTime().Format(time.RFC3339))
}
```

> Most brokers keep limited local history by default. If you see too few rows, open the symbol chart in MT4 to force a deeper history download.

---

## 🎯 3) Filter by symbol

```go
sym := ptr.String("EURUSD")
from := timestamppb.New(time.Now().AddDate(0, 0, -30))
to   := timestamppb.New(time.Now())

hist, err := account.ShowOrdersHistory(ctx, from, to, sym)
if err != nil { return err }

fmt.Printf("%d history orders for %s\n", len(hist.GetOrders()), *sym)
```

---

## 📦 4) Batch by days (for long ranges)

Long queries can be heavy and slow. Batch by day/week and append results.

```go
func loadHistoryBatched(ctx context.Context, a *MT4Account, sym *string, days int) ([]*pb.Order, error) {
    var out []*pb.Order
    end := time.Now()
    start := end.AddDate(0, 0, -days)

    for cur := start; cur.Before(end); cur = cur.AddDate(0, 0, 1) {
        from := timestamppb.New(cur)
        to   := timestamppb.New(cur.AddDate(0, 0, 1))
        h, err := a.ShowOrdersHistory(ctx, from, to, sym)
        if err != nil { return nil, err }
        out = append(out, h.GetOrders()...)
        // small pause to be gentle with terminal
        if err := waitWithCtx(ctx, 150*time.Millisecond); err != nil { return nil, err }
    }
    return out, nil
}
```

---

## 📤 5) Export to CSV (quick ‘n’ dirty)

```go
func saveCSV(path string, orders []*pb.Order) error {
    f, err := os.Create(path)
    if err != nil { return err }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    _ = w.Write([]string{"ticket","symbol","type","lots","profit","open","close"})
    for _, o := range orders {
        _ = w.Write([]string{
            strconv.FormatInt(o.GetTicket(), 10),
            o.GetSymbol(),
            o.GetType().String(),
            fmt.Sprintf("%.2f", o.GetVolume()),
            fmt.Sprintf("%.2f", o.GetProfit()),
            o.GetOpenTime().AsTime().Format(time.RFC3339),
            o.GetCloseTime().AsTime().Format(time.RFC3339),
        })
    }
    return w.Error()
}
```

---

## ⚠️ Pitfalls

* **Too few results** → MT4 hasn’t downloaded that range yet. Open the symbol chart or split the range into smaller chunks.
* **Timezones** → server time vs local time may differ; use `time.RFC3339` to avoid confusion.
* **Performance** → prefer batching for 30+ days ranges; add tiny sleeps between requests.
* **Context** → wrap calls with `context.WithTimeout` (3–6s) and retry only transport errors (see *Reliability*).

---

## 🔗 See also

* `CloseOrder.md` — verify closed positions appear in history.
* `ModifyOrder.md` — changes in SL/TP reflected in history records.
* `Reliability (en)` — timeouts, reconnects & backoff patterns.
