# 📡 StreamQuotes (GoMT4)

**Goal:** subscribe to **live quotes** for one or many symbols with auto‑reconnect and clean cancellation.

> Real code refs:
>
> * Streaming helpers: `examples/mt4/MT4Account.go` (`OnSymbolTick`, `ExecuteStreamWithReconnect`)
> * Demo usage: `examples/mt4/MT4_service.go` (`StreamQuotes`)

---

## ✅ 1) Preconditions

* `account` is connected; the symbols exist on the broker.
* You have a parent `ctx` to control lifetime.

---

## ▶️ 2) Subscribe to symbols

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

symbols := []string{"EURUSD", "GBPUSD"}

// Helper starts a stream and returns two channels:
//   dataCh: ticks (*pb.OnSymbolTickData)
//   errCh : the final error (after the built-in retries)
dataCh, errCh := account.OnSymbolTick(ctx, symbols)

for {
    select {
    case <-ctx.Done():
        return // stop gracefully
    case err := <-errCh:
        log.Printf("stream stopped: %v", err)
        return
    case pkt, ok := <-dataCh:
        if !ok {
            return // stream closed
        }
        t := pkt.GetSymbolTick()
        fmt.Printf("%s Bid=%.5f Ask=%.5f @ %s\n",
            t.GetSymbol(), t.GetBid(), t.GetAsk(), t.GetTime().AsTime().Format(time.RFC3339))
    }
}
```

> Under the hood the helper auto‑reconnects on `io.EOF` / `codes.Unavailable` / terminal-not-found using backoff+jitter from `MT4Account.go`.

---

## 🧭 3) Backoff & health checks

* Reconnects use exponential backoff with jitter (see constants in `MT4Account.go`).
* `ConnectByServerName` / `ConnectByHostPort` already run a short health‑check (`AccountSummary` with 3s timeout) after connecting.

---

## 🧰 4) Processing tips

* Keep the loop **non‑blocking**: if you do heavy work (DB, strategy), hand off ticks into a buffered channel/worker.
* Log the spread in points: `(Ask-Bid)/Point` (`Point` from `SymbolParams`) to catch broker anomalies.
* If you need OHLC bars, aggregate incoming ticks by timeframe on your side (or use `QuoteHistory`).

---

## ⚠️ Pitfalls

* **No reads** → if you stop reading `dataCh`, back‑pressure stalls the stream.
* **Symbol names** → suffixes like `EURUSD.m` are different instruments; use exact names.
* **Context canceled** → stream ends; always watch `<-ctx.Done()`.
* **Network flaps** → tune `backoffBase`/`backoffMax` (Reliability chapter).

---

## 🔄 Variations

* **Single symbol**: pass `[]string{"EURUSD"}`.
* **Dynamic subscribe**: keep a registry and restart the stream (cancel + new `OnSymbolTick`) with a new symbol set when needed.
* **Parallel consumers**: fan‑out ticks to multiple goroutines via a broadcast channel.

---

## 📎 See also

* [GetQuote](GetQuote.md) — one‑shot quote.
* [GetMultipleQuotes](GetMultipleQuotes.md) — batch fetch.
* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — timeouts, reconnects, backoff.
