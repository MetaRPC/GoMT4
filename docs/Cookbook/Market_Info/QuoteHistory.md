# 🕰️ QuoteHistory (GoMT4)

**Goal:** load **historical OHLC bars** for a symbol in a given time range — with safe batching for long periods.

> Real code refs in this repo:
>
> * Market info demo: `examples/mt4/MT4_service.go` (`ShowQuoteHistory`, `StreamQuoteHistoryExample`)
> * Account layer: `examples/mt4/MT4Account.go` (`QuoteHistory`, `QuoteHistoryStream`)

---

## ✅ 1) Preconditions

* `account` is connected; the symbol exists on the broker.
* `config.json` contains a valid `DefaultSymbol` or you pass a symbol explicitly.

---

## 🗓️ 2) Basic request (last N days)

`from` / `to` are sent as `FromTime` / `ToTime`, which the server reads as **broker server time**.

```go
sym := "EURUSD"
to := time.Now()
from := to.AddDate(0, 0, -7)

hist, err := account.QuoteHistory(ctx, sym, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, from, to)
if err != nil { return err }

for _, b := range hist.GetHistoricalQuotes() {
    t := b.GetTime().AsTime().Format("2006-01-02 15:04")
    fmt.Printf("%s %s O=%.5f H=%.5f L=%.5f C=%.5f\n", sym, t, b.GetOpen(), b.GetHigh(), b.GetLow(), b.GetClose())
}
```

> Mirrors the `ShowQuoteHistory` example in `MT4_service.go` (H1, last 5 days).

---

## 📦 3) Batch long ranges (chunked)

Large ranges can be heavy. `QuoteHistoryStream` splits the range into chunks and calls `QuoteHistory` per chunk:

```go
func LoadQuoteHistoryBatched(ctx context.Context, a *mt4.MT4Account, sym string, days int) ([]*pb.HistoryQuote, error) {
    to := time.Now()
    from := to.AddDate(0, 0, -days)

    barsCh, errCh := a.QuoteHistoryStream(ctx, sym, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M15, from, to, 24*time.Hour)

    var out []*pb.HistoryQuote
    for batch := range barsCh {
        out = append(out, batch.GetHistoricalQuotes()...)
    }
    // both channels are closed when the helper finishes; errCh holds the error, if any
    if err, ok := <-errCh; ok && err != nil {
        return out, err
    }
    return out, nil
}
```

---

## 🧪 4) Validate & spread (points)

Each bar carries its `Spread` (in points):

```go
for _, b := range bars {
    if b.GetOpen() <= 0 || b.GetClose() <= 0 { continue }
    if b.GetSpread() > 50 { // sanity check
        log.Printf("wide spread @ %s: %d points", b.GetTime().AsTime(), b.GetSpread())
    }
}
```

---

## 📤 5) Export to CSV

```go
func SaveBarsCSV(path, symbol string, bars []*pb.HistoryQuote) error {
    f, err := os.Create(path)
    if err != nil { return err }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    _ = w.Write([]string{"time","symbol","open","high","low","close","tick_volume","spread"})
    for _, b := range bars {
        _ = w.Write([]string{
            b.GetTime().AsTime().Format(time.RFC3339),
            symbol,
            fmt.Sprintf("%.5f", b.GetOpen()),
            fmt.Sprintf("%.5f", b.GetHigh()),
            fmt.Sprintf("%.5f", b.GetLow()),
            fmt.Sprintf("%.5f", b.GetClose()),
            fmt.Sprintf("%d", b.GetTickVolume()),
            fmt.Sprintf("%d", b.GetSpread()),
        })
    }
    w.Flush()
    return w.Error()
}
```

---

## ⚠️ Pitfalls

* **Empty result** → the terminal may not have downloaded that range yet; repeat the request after a short pause.
* **Timezones** → `Time` is broker server time (`TimeUtc` is UTC); `from`/`to` are read as broker server time too. `time.RFC3339` output of `Time` shows the broker clock with a `Z` suffix — use `TimeUtc` when you need true UTC.
* **Too big queries** → prefer chunking (`QuoteHistoryStream`).
* **Context & retries** → the call has an 8s default timeout and retries transient errors itself (see *Reliability*).

---

## 🔗 See also

* [GetQuote](GetQuote.md) — one‑shot live quote.
* [StreamQuotes](StreamQuotes.md) — continuous updates.
* [HistoryOrders](../Orders/HistoryOrders.md) — account history records.
