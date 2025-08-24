# 🕰️ QuoteHistory (GoMT4)

**Goal:** load **historical quotes** (Bid/Ask with timestamps) for a symbol in a given time range — with safe batching for long periods.

> Real code refs in this repo:
>
> * Market info demo: `examples/mt4/MT4_service.go` (see *ShowQuoteHistory* flow)
> * Account layer: `examples/mt4/MT4Account.go` (history wrapper & helpers)

---

## ✅ 1) Preconditions

* MT4 terminal is connected and the symbol is **visible** in *Market Watch*.
* `config.json` contains a valid `DefaultSymbol` or you pass a symbol explicitly.

---

## 🗓️ 2) Basic request (last N days)

```go
sym  := "EURUSD"
from := timestamppb.New(time.Now().AddDate(0, 0, -7))
to   := timestamppb.New(time.Now())

hist, err := account.ShowQuoteHistory(ctx, sym, from, to)
if err != nil { return err }

for _, q := range hist.GetQuotes() {
    t := q.GetTime().AsTime().Format(time.RFC3339)
    fmt.Printf("%s %s Bid=%.5f Ask=%.5f\n", sym, t, q.GetBid(), q.GetAsk())
}
```

> Mirrors the *ShowQuoteHistory* example flow in `MT4_service.go`.

---

## 📦 3) Batch long ranges (day-by-day)

Large ranges can be heavy. Split by days (or hours) and append.

```go
func LoadQuoteHistoryBatched(ctx context.Context, a *MT4Account, sym string, days int) ([]*pb.Quote, error) {
    var out []*pb.Quote
    end := time.Now()
    start := end.AddDate(0, 0, -days)

    for cur := start; cur.Before(end); cur = cur.AddDate(0, 0, 1) {
        from := timestamppb.New(cur)
        to   := timestamppb.New(cur.AddDate(0, 0, 1))
        h, err := a.ShowQuoteHistory(ctx, sym, from, to)
        if err != nil { return nil, err }
        out = append(out, h.GetQuotes()...)
        // Be gentle with terminal
        if err := waitWithCtx(ctx, 150*time.Millisecond); err != nil { return nil, err }
    }
    return out, nil
}
```

---

## 🧪 4) Validate & compute spread (pips)

```go
for _, q := range quotes {
    if q.GetAsk() <= 0 || q.GetBid() <= 0 { continue }
    spreadPips := (q.GetAsk() - q.GetBid()) / q.GetPoint()
    if spreadPips > 50 { // sanity check
        log.Printf("wide spread @ %s: %.1f pips", q.GetTime().AsTime(), spreadPips)
    }
}
```

---

## 📤 5) Export to CSV

```go
func SaveQuotesCSV(path, symbol string, quotes []*pb.Quote) error {
    f, err := os.Create(path)
    if err != nil { return err }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    _ = w.Write([]string{"time","symbol","bid","ask","spread_pips"})
    for _, q := range quotes {
        spread := (q.GetAsk() - q.GetBid()) / q.GetPoint()
        _ = w.Write([]string{
            q.GetTime().AsTime().Format(time.RFC3339),
            symbol,
            fmt.Sprintf("%.5f", q.GetBid()),
            fmt.Sprintf("%.5f", q.GetAsk()),
            fmt.Sprintf("%.1f", spread),
        })
    }
    return w.Error()
}
```

---

## ⚠️ Pitfalls

* **Empty result** → MT4 hasn’t downloaded that range yet; open the symbol chart to force data load.
* **Timezones** → MT4 server time may differ from local; always log with `time.RFC3339`.
* **Too big queries** → prefer batching; add tiny sleeps between requests.
* **Context & retries** → wrap with short timeout (3–6s) and retry only **transport** errors (see *Reliability*).

---

## 🔗 See also

* `GetQuote.md` — one‑shot live quote.
* `StreamQuotes.md` — continuous updates.
* `HistoryOrders.md` — account history records.
