# Getting Historical Quote Data

> **Request:** retrieve historical OHLC (candlestick) data for a given symbol
> Returns a list of time-based bars with open, high, low, close prices over a defined time range and timeframe.

---

### Code Example

```go
// Using service wrapper
service.ShowQuoteHistory(context.Background(), "EURUSD")

// Or directly from MT4Account
from := time.Now().AddDate(0, 0, -5)
to := time.Now()
timeframe := pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1

data, err := mt4.QuoteHistory(context.Background(), "EURUSD", timeframe, from, to)
if err != nil {
    log.Fatalf("Error fetching quote history: %v", err)
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

Returns a `QuoteHistoryData` object with:

| Field              | Type                   | Description                      |
| ------------------ | ---------------------- | -------------------------------- |
| `HistoricalQuotes` | `[]HistoricalQuoteBar` | List of historical OHLC candles. |

Each `HistoricalQuoteBar` includes:

| Field   | Type        | Description            |
| ------- | ----------- | ---------------------- |
| `Time`  | `timestamp` | Time of the bar (UTC). |
| `Open`  | `float64`   | Opening price.         |
| `High`  | `float64`   | Highest price.         |
| `Low`   | `float64`   | Lowest price.          |
| `Close` | `float64`   | Closing price.         |

---

## 🎯 Purpose

Use this method to load candlestick-style **historical price data** for a symbol. Ideal for:

* Charting historical candles
* Backtesting trading strategies
* Detecting technical analysis patterns

Supports all major timeframes using `ENUM_QUOTE_HISTORY_TIMEFRAME`.
