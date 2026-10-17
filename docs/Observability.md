# 📈 Observability (Logs & Metrics) — *optional enhancement*

> ⚠️ Note: this section is **not yet in the codebase**, but is recommended for developers who want better visibility.

You can add lightweight observability to GoMT4:

- 🌐 **obs.go helper** — wrapper around `log.Printf` so logs have consistent format (INFO/WARN/ERROR + timestamp).
- ⏱ **Latency timers** — measure RPC duration (e.g., ShowQuote, ShowOrderSend).
- 📊 **Counters** — simple `atomic` metrics for errors, ticks, RPC calls (print every 10s).
- 🧩 **Prometheus exporter** — optional, if you want Grafana dashboards.

Make GoMT4 **debuggable in minutes**, not hours. Below are pragmatic logging patterns and lightweight metrics that fit your current codebase (`examples/mt4/*.go`, `main.go`).

---

## 🎯 Goals

* **See**: what RPCs happen, how often, with what latency.
* **Spot**: reconnect loops, stream stalls, broker rejections.
* **Prove**: system is healthy via simple counters/ratios.

---

## 🧭 Logging levels (simple & effective)

Use three coarse levels via the standard `log` pkg (no heavy deps):

* `INFO` — high‑level events (connect, subscribe, order actions).
* `WARN` — transient issues (retry, backoff, deadline exceeded).
* `ERROR` — final failures (stream aborted, order rejected).

**Helper:**

```go
package obs
import (
  "fmt"; "log"; "time"
)
func ts() string { return time.Now().Format(time.RFC3339) }
func Info(msg string, a ...any)  { log.Printf("INFO  %s | "+msg, append([]any{ts()}, a...)...) }
func Warn(msg string, a ...any)  { log.Printf("WARN  %s | "+msg, append([]any{ts()}, a...)...) }
func Error(msg string, a ...any) { log.Printf("ERROR %s | "+msg, append([]any{ts()}, a...)...) }
```

Place in `examples/mt4/obs/obs.go` (or inline now). Replace current `fmt.Println(...)` in hot paths.

**Redaction** (from Security & Secrets): never print passwords; mask tokens.

---

## 🔌 Streams — what to log

Target minimal but useful lines.

**`StreamQuotes`** (`examples/mt4/MT4_service.go`):

```go
obs.Info("stream.quotes.start symbols=%v", symbols)
// on first N ticks only (sampling)
if ticks%100 == 0 { obs.Info("stream.quotes.rate ticks=%d", ticks) }
// on error
obs.Error("stream.quotes.error err=%v", err)
// on normal end
obs.Info("stream.quotes.end ticks=%d", ticks)
```

**`StreamOpenedOrderProfits`**:

```go
obs.Info("stream.pnl.start interval_ms=%d", 1000)
if n%50 == 0 { obs.Info("stream.pnl.rate updates=%d", n) }
if err != nil { obs.Error("stream.pnl.error err=%v", err) }
```

**Reconnect/backoff** (Cookbook → `HandleReconnect.md`). The built‑in retry loops (`ExecuteWithReconnect`, `ExecuteStreamWithReconnect` in `examples/mt4/MT4Account.go`) do not log; add a line before their `waitWithCtx(...)` calls in your clone:

```go
obs.Warn("reconnect reason=%v delay=%s attempt=%d", reason, delay, attempt)
```

---

## ⏱️ Unary RPCs — latency & outcomes

Wrap calls with a tiny stopwatch. Time the `MT4Account` methods (they return data and errors; the `MT4Service.Show*` helpers only print):

```go
t0 := time.Now()
q, err := account.Quote(ctx, symbol)
dur := time.Since(t0)
if err != nil {
  obs.Error("rpc.quote symbol=%s dur=%s err=%v", symbol, dur, err)
} else {
  obs.Info("rpc.quote symbol=%s dur=%s bid=%.5f ask=%.5f", symbol, dur, q.GetBid(), q.GetAsk())
}
```

Apply to: `QuoteMany`, `SymbolParams`, `TickValueWithSize`, history getters (`OrdersHistory`, `QuoteHistory`), and order actions (`OrderSend`, `OrderModify`, `OrderClose`, `OrderCloseBy`, `OrderDelete`).

---

## 🧮 Minimal metrics (standard lib)

No Prometheus? Start with **process‑local counters** and periodic prints.

```go
var (
  ticks uint64; pnlUpdates uint64; rpcOK uint64; rpcErr uint64
)
func incr(p *uint64) { atomic.AddUint64(p, 1) }

// In streams
incr(&ticks)
// In rpc success/fail
incr(&rpcOK); incr(&rpcErr)

// Reporter goroutine
func StartMetricsLogger() {
  go func(){
    t := time.NewTicker(10 * time.Second)
    for range t.C {
      ok := atomic.LoadUint64(&rpcOK)
      er := atomic.LoadUint64(&rpcErr)
      tk := atomic.LoadUint64(&ticks)
      obs.Info("metrics rpc_ok=%d rpc_err=%d ticks=%d", ok, er, tk)
    }
  }()
}
```

Call `StartMetricsLogger()` once in `main.go`.

**Why it helps**: you instantly see rates and error spikes without extra tooling.

---

## 📊 (Optional) Prometheus exporter

If you later add infra, expose counters via `promhttp`.

```go
var (
  mRpcOK  = prometheus.NewCounter(prometheus.CounterOpts{Name: "gomt4_rpc_ok"})
  mRpcErr = prometheus.NewCounter(prometheus.CounterOpts{Name: "gomt4_rpc_err"})
  mTicks  = prometheus.NewCounter(prometheus.CounterOpts{Name: "gomt4_ticks_total"})
)
func init(){ prometheus.MustRegister(mRpcOK, mRpcErr, mTicks) }
// http.ListenAndServe(":2112", promhttp.Handler())
```

Keep sampling in logs even with Prometheus for quick local debugging.

---

## 🧵 Sampling & log volume

* **Sample** hot events (e.g., every 50th tick; or log once per symbol per second).
* **Group** repetitive warnings (e.g., backoff bursts) with attempt counters.
* **Bound** log file size if redirecting to disk (PowerShell: `Start-Transcript` or use a rotating writer).

---

## 📍 Code map (repo anchors)

* `examples/mt4/MT4_service.go` → replace `fmt.Println` in:

  * `StreamQuotes`, `StreamOpenedOrderProfits`, `StreamOpenedOrderTickets`, `StreamTradeUpdates`
  * `ShowQuote`, `ShowQuotesMany`, `ShowQuoteHistory`
  * `ShowOrderSendExample`, `ShowOrderModifyExample`, `ShowOrderCloseExample`
* `examples/main.go` → start metrics reporter early and print final summary on exit.

---

## ✅ Checklist

* [ ] Log start/end of each stream; sample mid‑flow ticks.
* [ ] Log RPC latency and result (ok/err) with symbol/ticket context.
* [ ] Keep secrets redacted.
* [ ] Add 10‑second metrics reporter (counters) in `main.go`.
* [ ] Backoff logs include reason + delay + attempt.

---

### See also

* **Performance Notes** — hot paths & batching
* **Security & Secrets** — redaction helper
* **Cookbook / Reliability** — `HandleReconnect`, `UnaryRetries`
