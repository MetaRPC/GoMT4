# QuoteHistory

> OHLC bars for a symbol over a time range, at a chosen timeframe.

**API Information:**

* **Low-level API:** `MT4Account.QuoteHistory(ctx, symbol, timeframe, from, to)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.MarketInfo`
* **RPC:** `QuoteHistory(QuoteHistoryRequest) → QuoteHistoryReply`
* **Payload:** `QuoteHistoryData` → `GetHistoricalQuotes()` → `HistoryQuote`

## 💬 Just the essentials

* **What it is.** Candlestick (OHLC) history for a symbol/timeframe/window.
* **Why you need it.** Charts, indicators, backtests.
* **When to use.** Any historical-bar need. For a chunked, hands-off pull use [`QuoteHistoryStream`](../5.%20Streaming_Methods/QuoteHistoryStream.md).

---

## 🎯 Purpose

* Fetch OHLC bars for one symbol
* Pick the timeframe (M1 … MN1)
* Bound the range with `from` / `to`
* Read tick/real volume and spread per bar

---

## Method Signature

```go
func (a *MT4Account) QuoteHistory(
    ctx context.Context,
    symbol string,
    timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
    from, to time.Time,
) (*pb.QuoteHistoryData, error)
```

**Request message:**

```protobuf
message QuoteHistoryRequest {
  string symbol = 1;
  ENUM_QUOTE_HISTORY_TIMEFRAME timeframe = 2;
  Timestamp from_time = 3;
  Timestamp to_time = 4;
}
```

**Payload message:**

```protobuf
message QuoteHistoryData { repeated HistoryQuote historical_quotes = 1; }

message HistoryQuote {
  string symbol = ...;
  double open = ...;  double high = ...;  double low = ...;  double close = ...;
  int64  tick_volume = ...;
  int64  real_volume = ...;
  int32  spread = ...;
  Timestamp time = ...;
  Timestamp time_utc = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbol` | `string` | yes | Instrument, e.g. `"EURUSD"`. |
| `timeframe` | `pb.ENUM_QUOTE_HISTORY_TIMEFRAME` | yes | Bar period (see enum). |
| `from` | `time.Time` | yes | Range start (a value, not a pointer). |
| `to` | `time.Time` | yes | Range end. |

---

## ⬆️ Output — `*pb.QuoteHistoryData`

`GetHistoricalQuotes()` returns `[]*pb.HistoryQuote`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `open` / `high` / `low` / `close` | double | `GetOpen()` / `GetHigh()` / `GetLow()` / `GetClose()` | OHLC |
| `tick_volume` | int64 | `GetTickVolume() int64` | Tick volume |
| `real_volume` | int64 | `GetRealVolume() int64` | Real volume (if provided) |
| `spread` | int32 | `GetSpread() int32` | Spread in points |
| `time` / `time_utc` | Timestamp | `GetTime()` / `GetTimeUtc()` | Bar open time (server / UTC) |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧱 Related enums (from proto)

### `ENUM_QUOTE_HISTORY_TIMEFRAME` — input `timeframe`

| Constant | Value | Period |
|----------|-------|--------|
| `QH_PERIOD_M1` | 0 | 1 minute |
| `QH_PERIOD_M5` | 1 | 5 minutes |
| `QH_PERIOD_M15` | 2 | 15 minutes |
| `QH_PERIOD_M30` | 3 | 30 minutes |
| `QH_PERIOD_H1` | 4 | 1 hour |
| `QH_PERIOD_H4` | 5 | 4 hours |
| `QH_PERIOD_D1` | 6 | 1 day |
| `QH_PERIOD_W1` | 7 | 1 week |
| `QH_PERIOD_MN1` | 8 | 1 month |

> Pass directly: `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **`from`/`to` are `time.Time` values** (not pointers) — pass concrete times.
* **Bar time is the open** of the candle (`GetTime()`).
* **Big ranges** can be heavy — chunk with [`QuoteHistoryStream`](../5.%20Streaming_Methods/QuoteHistoryStream.md).

---

## 🔗 Usage Examples

### 1) Last 100 H1 candles

```go
import pb "github.com/MetaRPC/GoMT4/package"

func lastH1(ctx context.Context, acc *mt4.MT4Account, symbol string) error {
    to := time.Now()
    from := to.Add(-100 * time.Hour)
    h, err := acc.QuoteHistory(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, from, to)
    if err != nil { return err }
    fmt.Printf("%d bars\n", len(h.GetHistoricalQuotes()))
    return nil
}
```

### 2) Highest high / lowest low in the window

```go
func rangeHiLo(ctx context.Context, acc *mt4.MT4Account, symbol string, from, to time.Time) (hi, lo float64, err error) {
    h, err := acc.QuoteHistory(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, from, to)
    if err != nil { return 0, 0, err }
    for _, b := range h.GetHistoricalQuotes() {
        if b.GetHigh() > hi { hi = b.GetHigh() }
        if lo == 0 || b.GetLow() < lo { lo = b.GetLow() }
    }
    return hi, lo, nil
}
```

### 3) Export bars as CSV

```go
func exportBars(ctx context.Context, acc *mt4.MT4Account, symbol string, from, to time.Time) error {
    h, err := acc.QuoteHistory(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_D1, from, to)
    if err != nil { return err }
    for _, b := range h.GetHistoricalQuotes() {
        fmt.Printf("%s,%.5f,%.5f,%.5f,%.5f,%d\n",
            b.GetTime().AsTime().Format("2006-01-02"),
            b.GetOpen(), b.GetHigh(), b.GetLow(), b.GetClose(), b.GetTickVolume())
    }
    return nil
}
```

---

## ⚠️ Pitfalls

* **Empty windows.** Weekends/holidays or a tiny range can return zero bars — guard `len(...)`.
* **Time zone.** `from`/`to` are interpreted against the broker/server time — be explicit.
* **Large pulls.** Prefer the chunked stream for months of M1 data.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** A day of H1 returns ~24 bars on a trading day.
* **OHLC sanity.** For each bar `low ≤ open,close ≤ high`.
* **Empty.** A future range returns zero bars without error.

---

## 📚 See also

* [QuoteHistoryStream](../5.%20Streaming_Methods/QuoteHistoryStream.md) · [Quote](./Quote.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
