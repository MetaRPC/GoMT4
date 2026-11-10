# QuoteHistoryStream

> Pull a long history range as time-boxed chunks, each delivered through a channel.

**API Information:**

* **Low-level API:** `MT4Account.QuoteHistoryStream(ctx, symbol, timeframe, from, to, chunk)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **Built on:** [`QuoteHistory`](../3.%20Symbols_Quotes/QuoteHistory.md) — a client-side chunker, **not** a live server stream. It splits `[from..to]` into windows and calls `QuoteHistory` for each.
* **Payload:** `QuoteHistoryData` (one per chunk)

## 💬 Just the essentials

* **What it is.** A wrapper that fetches months of bars as manageable chunks and streams them.
* **Why you need it.** Backfill long ranges without one giant request or a hand-rolled window loop.
* **When to use.** Backtesting data, chart backfills.

---

## 🎯 Purpose

* Fetch a long OHLC range in bounded chunks
* Skip empty windows (weekends/holidays) automatically
* Stop early by cancelling the context
* Default the chunk size when you don't care

---

## Method Signature

```go
func (a *MT4Account) QuoteHistoryStream(
    ctx context.Context,
    symbol string,
    timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
    from, to time.Time,
    chunk time.Duration,   // 0 → defaults to 24h
) (<-chan *pb.QuoteHistoryData, <-chan error)
```

It walks `from → to` in steps of `chunk`, calling [`QuoteHistory`](../3.%20Symbols_Quotes/QuoteHistory.md) per window and sending non-empty results down the channel.

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to stop early and close the channels. |
| `symbol` | `string` | yes | Instrument, e.g. `"EURUSD"`. |
| `timeframe` | `pb.ENUM_QUOTE_HISTORY_TIMEFRAME` | yes | Bar period (see [`QuoteHistory`](../3.%20Symbols_Quotes/QuoteHistory.md#related-enums-from-proto)). |
| `from` | `time.Time` | yes | Range start. |
| `to` | `time.Time` | yes | Range end (must be ≥ `from`). |
| `chunk` | `time.Duration` | no | Window size per request; `0` → 24h. |

---

## ⬆️ Frame — `*pb.QuoteHistoryData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `historical_quotes` | repeated HistoryQuote | `GetHistoricalQuotes()` | Bars in this chunk (OHLC — see [`QuoteHistory`](../3.%20Symbols_Quotes/QuoteHistory.md)) |

> The data channel is **closed** (you receive `nil`) when the range is done. `to < from` is rejected with `invalid range: to < from` on the error channel.

---

## 🧩 Notes & Tips

* **Empty windows are skipped** — weekends/holidays just produce no frame, not an error.
* **`chunk = 0` → 24h** — a sensible default if you don't care.
* **Not live** — it completes over the fixed range; it won't push new bars.

---

## 🔗 Usage Examples

### 1) Count bars over the last month (weekly chunks)

```go
import pb "github.com/MetaRPC/GoMT4/package"

func countBars(ctx context.Context, acc *mt4.MT4Account, symbol string) (int, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    from, to := time.Now().AddDate(0, -1, 0), time.Now()
    chunks, errs := acc.QuoteHistoryStream(ctx, symbol,
        pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, from, to, 7*24*time.Hour)
    bars := 0
    for {
        select {
        case c := <-chunks:
            if c == nil { return bars, nil }
            bars += len(c.GetHistoricalQuotes())
        case err := <-errs:
            if err != nil { return 0, err }
        case <-ctx.Done():
            return bars, nil
        }
    }
}
```

### 2) Collect every bar for a backtest

```go
func collectBars(ctx context.Context, acc *mt4.MT4Account, symbol string) ([]*pb.HistoryQuote, error) {
    from := time.Now().AddDate(0, -3, 0) // last 3 months
    chunks, errs := acc.QuoteHistoryStream(ctx, symbol,
        pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, from, time.Now(), 0) // 0 → 24h
    var all []*pb.HistoryQuote
    for {
        select {
        case c := <-chunks:
            if c == nil { return all, nil }
            all = append(all, c.GetHistoricalQuotes()...)
        case err := <-errs:
            if err != nil { return nil, err }
        case <-ctx.Done():
            return all, nil
        }
    }
}
```

### 3) Daily highest-high scan

```go
func monthHigh(ctx context.Context, acc *mt4.MT4Account, symbol string) (float64, error) {
    from := time.Now().AddDate(0, -1, 0)
    chunks, errs := acc.QuoteHistoryStream(ctx, symbol,
        pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_D1, from, time.Now(), 0)
    var hi float64
    for {
        select {
        case c := <-chunks:
            if c == nil { return hi, nil }
            for _, b := range c.GetHistoricalQuotes() {
                if b.GetHigh() > hi { hi = b.GetHigh() }
            }
        case err := <-errs:
            if err != nil { return 0, err }
        case <-ctx.Done():
            return hi, nil
        }
    }
}
```

---

## ⚠️ Pitfalls

* **Channel close = `nil`** — check for it to know you're done.
* **`to < from`** is rejected — order your range correctly.
* **Not live** — it won't push bars that form after it starts.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** A month of H1 yields several chunks with bars; the summed count is plausible.
* **Empty windows.** A weekend-only range produces no frames, no error.
* **Bad range.** `to < from` sends the range error.

---

## 📚 See also

* [QuoteHistory](../3.%20Symbols_Quotes/QuoteHistory.md) · [OnSymbolTick](./OnSymbolTick.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
