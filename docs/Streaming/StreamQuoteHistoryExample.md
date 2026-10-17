# 📝 Streaming Quote History (Example)

> **Request:** Stream historical price bars (OHLC) using safe defaults.
> This demo wrapper prints batches of candles and exits after a short timeout.

---

### Code Example

```go
// Quick use: last 90 days of H1 candles for a symbol,
// delivered in weekly chunks. The demo loop stops after 5s without a batch.
svc.StreamQuoteHistoryExample(ctx, "EURUSD")
```

---

### Method Signature

```go
// Demo wrapper: safe defaults + console printing
func (s *MT4Service) StreamQuoteHistoryExample(ctx context.Context, symbol string)
```

**Defaults used inside the example wrapper:**

* Date range: `from = now - 90 days`, `to = now`
* Timeframe: `QH_PERIOD_H1`
* Chunk duration: `7 * 24h` (weekly)
* Demo timeout: 5 seconds without a new batch (select-case exit)

Under the hood, it calls the service API, which proxies `MT4Account.QuoteHistoryStream`:

```go
func (s *MT4Service) StreamQuoteHistory(
    ctx context.Context,
    symbol string,
    timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
    from, to time.Time,
    chunk time.Duration,
) (<-chan *pb.QuoteHistoryData, <-chan error)
```

---

## 🔽 Input

|    Field | Type              | Required | Description                                  |
| -------: | ----------------- | :------: | -------------------------------------------- |
|    `ctx` | `context.Context` |     ✅    | Controls lifetime (cancel/timeout) of stream |
| `symbol` | `string`          |     ✅    | Trading symbol, e.g. `"EURUSD"`              |

> The **Example** wrapper fixes timeframe/chunk/range to sensible defaults. Use the service API for full control.

---

## ⬆️ Output

**Console output** of OHLC bars from each `*pb.QuoteHistoryData` batch. For each candle, the example prints:

* `Time` (broker server time),
* `Open`, `Close` (and can be extended to `High`, `Low`).

Sample line format:

```
[YYYY-mm-dd HH:MM:SS] O: <open> C: <close>
```

> Each `*pb.HistoryQuote` also has `High`, `Low`, `TickVolume`, `RealVolume`, `Spread` and `TimeUtc`. Print whatever you need.

How the chunking works: `QuoteHistoryStream` splits `[from, to]` into windows of `chunk` (24h if `chunk <= 0`), calls `QuoteHistory` for each window, skips empty windows and sends the others to the channel; it closes both channels after the last window or on an error (`to` before `from` is an error).

---

## ENUM: `pb.ENUM_QUOTE_HISTORY_TIMEFRAME`

Values (Go constants `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_<value>`):

* `QH_PERIOD_M1`, `QH_PERIOD_M5`, `QH_PERIOD_M15`, `QH_PERIOD_M30`
* `QH_PERIOD_H1`, `QH_PERIOD_H4`
* `QH_PERIOD_D1`, `QH_PERIOD_W1`, `QH_PERIOD_MN1`

---

## 🧩 Notes & Tips

* **Connection required:** Ensure `ConnectByServerName` or `ConnectByHostPort` succeeded before streaming.
* **Chunking:** Larger `chunk` ⇒ fewer round-trips (better throughput). Smaller `chunk` ⇒ lower memory spikes.
* **Auto-reconnect:** Each chunk is a unary `QuoteHistory` call that retries on transient gRPC errors and terminal-not-found errors with exponential backoff + jitter.
* **Back-pressure:** If you do heavy work per batch, offload to a worker so the reader loop never blocks.

---

## ⚠️ Pitfalls

* **Ignoring `errCh`:** Always select on the error channel to catch terminal failures.
* **Huge ranges:** Very wide ranges on small timeframes can generate large data; consider paging by time (chunks) appropriately.
* **Timezone confusion:** `Time` is broker server time (`TimeUtc` is UTC), and the `from`/`to` you pass are read as broker server time.

---

## 🧪 Testing

* **Happy path:** Expect multiple batches for the last 90 days of H1 data.
* **Network hiccups:** Brief outages should be retried automatically; expect a short pause.
* **Cancel path:** Cancel the passed `ctx` and verify both channels close and the goroutine exits.
