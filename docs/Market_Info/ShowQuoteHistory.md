# Getting Historical Quote Data

> **Request:** retrieve historical OHLC (candlestick) data for a given symbol
> Returns a list of time-based bars with open, high, low, close prices over a defined time range and timeframe.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints recent OHLC candles for the symbol.
svc.ShowQuoteHistory(ctx, "EURUSD")

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

from := time.Now().AddDate(0, 0, -5)
to   := time.Now()
timeframe := pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1

ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
defer cancel()

data, err := account.QuoteHistory(ctx, "EURUSD", timeframe, from, to)
if err != nil {
    log.Fatalf("❌ QuoteHistory error: %v", err)
}

for _, c := range data.GetHistoricalQuotes() {
    fmt.Printf("[%s] O: %.5f H: %.5f L: %.5f C: %.5f\n",
        c.GetTime().AsTime().Format("2006-01-02 15:04:05"),
        c.GetOpen(),
        c.GetHigh(),
        c.GetLow(),
        c.GetClose(),
    )
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowQuoteHistory(ctx context.Context, symbol string) // H1, last 5 days
func (a *MT4Account) QuoteHistory(
    ctx context.Context,
    symbol string,
    timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
    from, to time.Time,
) (*pb.QuoteHistoryData, error)
```

---

## 🔽 Input

| Field       | Type                           | Description                                |
| ----------- | ------------------------------ | ------------------------------------------ |
| `symbol`    | `string`                       | Trading symbol (e.g., "EURUSD").           |
| `timeframe` | `pb.ENUM_QUOTE_HISTORY_TIMEFRAME` | Bar timeframe (e.g., hourly, daily, etc.). |
| `from`      | `time.Time`                    | Start of the range (sent as `FromTime`, read as broker server time). |
| `to`        | `time.Time`                    | End of the range (sent as `ToTime`, read as broker server time). |
| `ctx`       | `context.Context`              | For timeout or cancellation.               |

### ENUM: `ENUM_QUOTE_HISTORY_TIMEFRAME`

Go constants are `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_<value>`, e.g. `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1`.

| Value           | Description    |
| --------------- | -------------- |
| `QH_PERIOD_M1`  | 1-minute bars  |
| `QH_PERIOD_M5`  | 5-minute bars  |
| `QH_PERIOD_M15` | 15-minute bars |
| `QH_PERIOD_M30` | 30-minute bars |
| `QH_PERIOD_H1`  | 1-hour bars    |
| `QH_PERIOD_H4`  | 4-hour bars    |
| `QH_PERIOD_D1`  | Daily bars     |
| `QH_PERIOD_W1`  | Weekly bars    |
| `QH_PERIOD_MN1` | Monthly bars   |

---

## ⬆️ Output

Returns `*pb.QuoteHistoryData`:

| Field              | Type                       | Description                      |
| ------------------ | -------------------------- | -------------------------------- |
| `HistoricalQuotes` | `[]*pb.HistoryQuote` | List of historical OHLC candles. |

Each `*pb.HistoryQuote` includes:

| Field        | Type                     | Description                                  |
| ------------ | ------------------------ | -------------------------------------------- |
| `Symbol`     | `string`                 | Symbol name.                                 |
| `Index`      | `int32`                  | Bar index.                                   |
| `Time`       | `*timestamppb.Timestamp` | Open time of the bar (broker server time).   |
| `TimeUtc`    | `*timestamppb.Timestamp` | The same moment in UTC.                      |
| `Open`       | `float64`                | Opening price.                               |
| `High`       | `float64`                | Highest price.                               |
| `Low`        | `float64`                | Lowest price.                                |
| `Close`      | `float64`                | Closing price.                               |
| `TickVolume` | `int64`                  | Tick volume.                                 |
| `RealVolume` | `int64`                  | Real volume (often 0 for forex).             |
| `Spread`     | `int32`                  | Spread in points.                            |

---

## 🎯 Purpose

Load candlestick-style **historical price data** for a symbol. Ideal for:

* Charting historical candles
* Backtesting trading strategies
* Detecting technical analysis patterns

---

## 🧩 Notes & Tips

* **Range limits:** Very large ranges may be truncated. Query in smaller chunks (`account.QuoteHistoryStream` does that).
* **Deprecated RPC:** `QuoteHistory` calls `MarketInfo.QuoteHistory`, which the server marks deprecated in favour of the paged `QuoteHistory.Bars` RPC (`pb.NewQuoteHistoryClient`, request `pb.QuoteHistoryBarsRequest`). The client does not wrap `Bars` yet; it still works through the generated client.
* **Timezone:** `Time` is broker server time, `TimeUtc` is UTC. The `from`/`to` you pass are read as broker server time, so `time.Now()` as `to` is shifted by the broker's UTC offset; widen the range if you need the latest bar.
* **Gaps:** Weekend and holiday gaps are normal; don’t misinterpret them as missing data.

---

## ⚠️ Pitfalls

* **Inconsistent data:** Different brokers can have slightly different history for the same symbol.
* **No future bars:** `to` cannot exceed the server’s latest bar — it will return empty beyond.
* **Alignment:** Bars align strictly to their timeframe (e.g., H1 always at :00 minutes).

---

## 🧪 Testing Suggestions

* **Happy path:** Request last 100 H1 bars for EURUSD → verify consistent OHLC values.
* **Edge case:** Request with `from > to` → the client returns `invalid time range: to < from` without calling the server.
* **Stress test:** Fetch several months of M1 data → ensure chunking or iteration works.
