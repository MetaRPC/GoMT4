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
func (s *MT4Service) ShowQuoteHistory(ctx context.Context, symbol string)
```

---

## 🔽 Input

| Field       | Type                           | Description                                |
| ----------- | ------------------------------ | ------------------------------------------ |
| `symbol`    | `string`                       | Trading symbol (e.g., "EURUSD").           |
| `timeframe` | `ENUM_QUOTE_HISTORY_TIMEFRAME` | Bar timeframe (e.g., hourly, daily, etc.). |
| `from`      | `time.Time`                    | Start of the historical range (UTC).       |
| `to`        | `time.Time`                    | End of the historical range (UTC).         |
| `ctx`       | `context.Context`              | For timeout or cancellation.               |

### ENUM: `ENUM_QUOTE_HISTORY_TIMEFRAME`

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
| `HistoricalQuotes` | `[]*pb.HistoricalQuoteBar` | List of historical OHLC candles. |

Each `*pb.HistoricalQuoteBar` includes:

| Field   | Type        | Description            |
| ------- | ----------- | ---------------------- |
| `Time`  | `timestamp` | Time of the bar (UTC). |
| `Open`  | `float64`   | Opening price.         |
| `High`  | `float64`   | Highest price.         |
| `Low`   | `float64`   | Lowest price.          |
| `Close` | `float64`   | Closing price.         |

---

## 🎯 Purpose

Load candlestick-style **historical price data** for a symbol. Ideal for:

* Charting historical candles
* Backtesting trading strategies
* Detecting technical analysis patterns

---

## 🧩 Notes & Tips

* **Range limits:** Very large ranges may be truncated. Query in smaller chunks.
* **Timezone:** All times are UTC. Convert to local TZ for charts.
* **Gaps:** Weekend and holiday gaps are normal; don’t misinterpret them as missing data.

---

## ⚠️ Pitfalls

* **Inconsistent data:** Different brokers can have slightly different history for the same symbol.
* **No future bars:** `to` cannot exceed the server’s latest bar — it will return empty beyond.
* **Alignment:** Bars align strictly to their timeframe (e.g., H1 always at :00 minutes).

---

## 🧪 Testing Suggestions

* **Happy path:** Request last 100 H1 bars for EURUSD → verify consistent OHLC values.
* **Edge case:** Request with `from > to` → expect error/empty response.
* **Stress test:** Fetch several months of M1 data → ensure chunking or iteration works.
